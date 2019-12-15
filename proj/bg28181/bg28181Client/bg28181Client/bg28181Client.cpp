#include <winsock.h>
#include "eXosip2/eXosip.h"
#include "eXosip2/eX_setup.h"
#include "eXosip2/eX_register.h"
#include "eXosip2/eX_options.h"
#include "eXosip2/eX_message.h"
#include "osip2/osip.h"
#include "osipparser2/osip_message.h"
#include "osipparser2/osip_parser.h"
#include "osipparser2/osip_port.h"
#include "bg28181Client.h"



bg28181Client::bg28181Client()
{

}

bg28181Client::~bg28181Client()
{

}

int bg28181Client::Initialize(const char *local_ip, unsigned short local_port, const char *local_gb_code, NET_TYPE net_type)
{
	int errCode = 0;

	// ��ʼ��eXosip����
	context_ = (void *)eXosip_malloc();
	eXosip_t *excontext = (eXosip_t *)context_;

	errCode = eXosip_init(excontext);
	if (errCode != 0)
	{
		return errCode;
	}

	eXosip_set_user_agent(excontext, "UAC/1.0");

	local_ip_ = local_ip;
	local_port_ = local_port;
	local_gbcode_ = local_gb_code;

	// ���ؼ���
	int transport = net_type == NET_TCP ? IPPROTO_TCP : IPPROTO_UDP;
	errCode = eXosip_listen_addr(excontext, transport, local_ip_.c_str(), local_port_, AF_INET, 0);
	if (errCode != 0)
	{
		// �޷������˿�
		eXosip_quit(excontext);
		return errCode;
	}

	// ����Ӧ������һ���̴߳�����յ�����Ϣ
	if (working_thread_.isRunning())
	{
		// �跨��ֹͣ�߳�����
		need_working_thread_exit_ = true;
		working_thread_.join();
	}
	else
	{
		need_working_thread_exit_ = false;
		working_thread_.start(WorkingThread, this);
	}

	return errCode;
}

int bg28181Client::Register(const char *server_ip, unsigned short server_port, const char *server_gbcode, const char *username, const char *password, int expired)
{
	int errCode = 0;
	eXosip_t *excontext = (eXosip_t *)context_;

	server_ip_ = server_ip;
	server_port_ = server_port;
	server_gbcode_ = server_gbcode;
	server_username_ = username;
	server_password_ = password;
	reg_expired_ = expired;

	// ���ȣ������Ȩ��Ϣ
	eXosip_clear_authentication_info(excontext);

	// ���������������Ӽ�Ȩ��Ϣ
	errCode = eXosip_add_authentication_info(excontext, username, username, password, "MD5", NULL);
	if (errCode > 0)
	{
		// ��Ӽ�Ȩ��Ϣʧ��
		return errCode;
	}

	// ��������֯����
	osip_message_t *reg = NULL;
	char proxy[1024] = { 0 };
	char from[1024] = { 0 };
	char contact[1024] = { 0 };
	sprintf(from, "sip:%s@%s:%d", username, server_ip, server_port);

	// ������������REGISTER������ŵ�һ�����ݣ����ﰴ��28181���ĵ�˵������ʽӦ��Ϊ�� sip:SIP����������@Ŀ��������IP��ַ�˿� 
	sprintf(proxy, "sip:%s@%s:%d", server_gbcode, server_ip, server_port);
	sprintf(contact, "sip:%s@%s:%d", local_gbcode_.c_str(), local_ip_.c_str(), local_port_);

	// ����
	eXosip_lock(excontext);
	int reg_id = eXosip_register_build_initial_register(excontext, from, proxy, contact, expired, &reg);
	if (reg_id < 0)
	{
		// ��ʼ����Ȩ��Ϣʧ��
		eXosip_unlock(excontext);
		eXosip_quit(excontext);
		return -1;
	}

	errCode = eXosip_register_send_register(excontext, reg_id, reg);
	eXosip_unlock(excontext);
	if (errCode != 0)
	{
		// ����ע��ʧ�ܣ�
		return errCode;
	}

	return errCode;
}

int bg28181Client::Unregister()
{
	int errCode = 0;

	return errCode;
}

int bg28181Client::ReregisterBy401(void* event_info)
{
	int errCode = 0;
	//eXosip_t *excontext = (eXosip_t *)context_;
	//eXosip_event_t *sip_event = (eXosip_event_t *)event_info;

	//if (sip_event->response->status_code == 401)
	//{
	//	// ˵����Ȩ��Ϣ�����⣬ʹ��sip_event��rid���²���һ��ע���
	//	osip_message_t *reg = NULL;
	//	errCode = eXosip_register_build_register(excontext, sip_event->rid, reg_expired_, &reg);
	//	
	//	// ȡ����֤���ַ���authorization
	//	osip_authorization_t *auth = NULL;
	//	char *strAuth = NULL;
	//	osip_message_get_authorization(reg, 0, &auth);
	//	osip_authorization_to_str(auth, &strAuth);

	//	//������֤�ַ���
	//	auth_info_ = strAuth;
	//	delete [] strAuth;

	//	eXosip_register_send_register(excontext, sip_event->rid, reg);
	//}

	return errCode;
}


void bg28181Client::WorkingThread(void* lpParam)
{
	bg28181Client *client = (bg28181Client *)lpParam;
	eXosip_t *excontext = (eXosip_t *)client->context_;
	eXosip_event_t *sip_event = NULL;
	osip_message_t *ack = NULL;
	osip_message_t *answer = NULL;

	for (;;)
	{
		sip_event = eXosip_event_wait(excontext, 0, 50);
		eXosip_lock(excontext);

		eXosip_default_action(excontext, sip_event);
		// û�ҵ���������Ķ��壬������ʱע�͵�����
		//eXosip_automatic_refresh();
		eXosip_automatic_action(excontext);
		eXosip_unlock(excontext);

		if (sip_event == NULL)
			continue;

		switch (sip_event->type)
		{
		// ע���йص��¼�
		case EXOSIP_REGISTRATION_SUCCESS:
			// user is registred
			// ע��ɹ�
			OutputDebugStringA("WorkingThread >>> EXOSIP_REGISTRATION_SUCCESS\n");
			
			break;
		case EXOSIP_REGISTRATION_FAILURE:
			// user is not registred
			// ע��ʧ�ܣ��ж��Ƿ���401����
			OutputDebugStringA("WorkingThread >>> EXOSIP_REGISTRATION_FAILURE\n");
			client->ReregisterBy401(sip_event);
			break;

		/* INVITE related events within calls */
		// �ڵ�������INVITE�йص��¼�
		case EXOSIP_CALL_INVITE:
			// announce a new call
			// һ���µ�CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_INVITE\n");
			break;
		case EXOSIP_CALL_REINVITE:
			// announce a new INVITE within call
			// ��һ��CALL�������յ�һ���µ�INVITE
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REINVITE\n");
			break;
		case EXOSIP_CALL_NOANSWER:
			// announce no answer within the timeout
			// ��ʱ����Ӧ
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_NOANSWER\n");
			break;
		case EXOSIP_CALL_PROCEEDING:
			// announce processing by a remote app
			// Զ��Ӧ�����ڴ������CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_PROCEEDING\n");
			break;
		case EXOSIP_CALL_RINGING:
			// announce ringback
			// ��������
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_RINGING\n");
			break;
		case EXOSIP_CALL_ANSWERED:
			// announce start of call
			// ��ʼһ��CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_ANSWERED\n");
			break;
		case EXOSIP_CALL_REDIRECTED:
			// announce a redirection
			// �ض���һ��CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REDIRECTED\n");
			break;
		case EXOSIP_CALL_REQUESTFAILURE:
			// announce a request failure
			// ��������ʧ��
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REQUESTFAILURE\n");
			break;
		case EXOSIP_CALL_SERVERFAILURE:
			// announce a server failure
			// һ������˴���
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_SERVERFAILURE\n");
			break;
		case EXOSIP_CALL_GLOBALFAILURE:
			// announce a global failure
			// һ��ȫ�ִ���
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_GLOBALFAILURE\n");
			break;
		case EXOSIP_CALL_ACK:
			// ACK received for 200ok to INVITE
			// INVITE��200 OK�Ѿ��յ�ACK��
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_ACK\n");
			break;

		case EXOSIP_CALL_CANCELLED:
			// announce that call has been cancelled
			// ȡ��һ��CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_CANCELLED\n");
			break;

		/* request related events within calls (except INVITE) */
		// �������������йص��¼�(INVITE����)
		case EXOSIP_CALL_MESSAGE_NEW:
			// announce new incoming request
			// 
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_NEW\n");
			break;
		case EXOSIP_CALL_MESSAGE_PROCEEDING:
			// announce a 1xx for request
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_PROCEEDING\n");
			break;
		case EXOSIP_CALL_MESSAGE_ANSWERED:
			// announce a 200ok
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_ANSWERED\n");
			break;
		case EXOSIP_CALL_MESSAGE_REDIRECTED:
			// announce a failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_REDIRECTED\n");
			break;
		case EXOSIP_CALL_MESSAGE_REQUESTFAILURE:
			// announce a failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_REQUESTFAILURE\n");
			break;
		case EXOSIP_CALL_MESSAGE_SERVERFAILURE:
			// announce a failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_SERVERFAILURE\n");
			break;
		case EXOSIP_CALL_MESSAGE_GLOBALFAILURE:
			// announce a failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_MESSAGE_GLOBALFAILURE\n");
			break;

		case EXOSIP_CALL_CLOSED:
			// a BYE was received for this call
			// �����CALL���յ�һ��BYE
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_CLOSED\n");
			break;

		/* for both UAS & UAC events */
		// UAS��UAC�����յ����¼�
		case EXOSIP_CALL_RELEASED:
			// call context is cleared
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_RELEASED\n");
			break;

		/* events received for request outside calls */
		// ���������ⲿ���õ��¼�
		case EXOSIP_MESSAGE_NEW:
			// announce new incoming request.
			// ���յ�һ����Ϣ
			// �����ϣ��ϼ�ƽ̨���е�������������յ���δ���Ե����������
			// Ϊ�˱��ϳ���Ч�ʣ��������յ������ֱ���ӵ���������У����氲���̳߳ؽ��д���
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_NEW\n");

			// �����ж��Ƿ���һ��MESSAGE
			if (MSG_IS_MESSAGE(sip_event->request))
			{
				osip_body_t *sip_body = NULL;
				osip_message_get_body(sip_event->request, 0, &sip_body);
				OutputDebugStringA(sip_body->body);
				OutputDebugStringA("\n");

				// ���յ��������ӽ����У��򻺴���
				// �ڷŽ��������ǰ��ȷ�Ϻûظ�ʱ��Ҫ��Щ��������֯һ���ṹ����л���
				// ʣ�µ�����������̳߳أ����̴߳�����Ϣ��
				//client->request_queue_.push();

				//���չ��򣬻ظ�200 OK��Ϣ
				osip_message_t *answer = NULL;
				eXosip_message_build_answer(excontext, sip_event->tid, 200, &answer);
				eXosip_message_send_answer(excontext, sip_event->tid, 200, answer);
			}
			break;
		case EXOSIP_MESSAGE_PROCEEDING:
			// announce a 1xx for request.
			// ���յ�һ��1XX������
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_PROCEEDING\n");
			break;
		case EXOSIP_MESSAGE_ANSWERED:
			// announce a 200ok
			// ���յ�һ��200 OK
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_ANSWERED\n");
			break;
		case EXOSIP_MESSAGE_REDIRECTED:
			// announce a failure.
			// ���յ�һ��ʧ����Ϣ
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_REDIRECTED\n");
			break;
		case EXOSIP_MESSAGE_REQUESTFAILURE:
			// announce a failure.
			// ���յ�һ��ʧ����Ϣ
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_REQUESTFAILURE\n");
			break;
		case EXOSIP_MESSAGE_SERVERFAILURE:
			// announce a failure.
			// ���յ�һ��ʧ����Ϣ
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_SERVERFAILURE\n");
			break;
		case EXOSIP_MESSAGE_GLOBALFAILURE:
			// announce a failure.
			// ���յ�һ��ʧ����Ϣ
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_GLOBALFAILURE\n");
			break;
					
		/* Presence and Instant Messaging */
		// ״̬�����뼴ʱͨ��
		case EXOSIP_SUBSCRIPTION_NOANSWER:
			// announce no answer
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_NOANSWER\n");
			break;
		case EXOSIP_SUBSCRIPTION_PROCEEDING:
			// announce a 1xx
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_PROCEEDING\n");
			break;
		case EXOSIP_SUBSCRIPTION_ANSWERED:
			// announce a 200ok
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_ANSWERED\n");
			break;
		case EXOSIP_SUBSCRIPTION_REDIRECTED:
			// announce a redirection
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_REDIRECTED\n");
			break;
		case EXOSIP_SUBSCRIPTION_REQUESTFAILURE:
			// announce a request failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_REQUESTFAILURE\n");
			break;
		case EXOSIP_SUBSCRIPTION_SERVERFAILURE:
			// announce a server failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_SERVERFAILURE\n");
			break;
		case EXOSIP_SUBSCRIPTION_GLOBALFAILURE:
			// announce a global failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_GLOBALFAILURE\n");
			break;
		case EXOSIP_SUBSCRIPTION_NOTIFY:
			// announce new NOTIFY request
			OutputDebugStringA("WorkingThread >>> EXOSIP_SUBSCRIPTION_NOTIFY\n");
			break;

		case EXOSIP_IN_SUBSCRIPTION_NEW:
			// announce new incoming SUBSCRIBE/REFER.*/
			OutputDebugStringA("WorkingThread >>> EXOSIP_IN_SUBSCRIPTION_NEW\n");
			break;

		case EXOSIP_NOTIFICATION_NOANSWER:
			// announce no answer
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_NOANSWER\n");
			break;
		case EXOSIP_NOTIFICATION_PROCEEDING:
			// announce a 1xx
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_PROCEEDING\n");
			break;
		case EXOSIP_NOTIFICATION_ANSWERED:
			// announce a 200ok
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_ANSWERED\n");
			break;
		case EXOSIP_NOTIFICATION_REDIRECTED:
			// announce a redirection
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_REDIRECTED\n");
			break;
		case EXOSIP_NOTIFICATION_REQUESTFAILURE:
			// announce a request failure
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_REQUESTFAILURE\n");
			break;
		case EXOSIP_NOTIFICATION_SERVERFAILURE:
			// announce a server failure
			// һ������˴���
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_SERVERFAILURE\n");
			break;
		case EXOSIP_NOTIFICATION_GLOBALFAILURE:
			// announce a global failure
			// һ��ȫ�ִ���
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_GLOBALFAILURE\n");
			break;
		default:
			OutputDebugStringA("WorkingThread >>> UNKNOWN EVENT\n");
			break;
		}

		eXosip_event_free(sip_event);
	}

	return ;
}