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

	// 初始化eXosip环境
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

	// 本地监听
	int transport = net_type == NET_TCP ? IPPROTO_TCP : IPPROTO_UDP;
	errCode = eXosip_listen_addr(excontext, transport, local_ip_.c_str(), local_port_, AF_INET, 0);
	if (errCode != 0)
	{
		// 无法监听端口
		eXosip_quit(excontext);
		return errCode;
	}

	// 这里应该启动一个线程处理接收到的消息
	if (working_thread_.isRunning())
	{
		// 设法先停止线程运行
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

	// 首先，清理鉴权信息
	eXosip_clear_authentication_info(excontext);

	// 接下来，重新增加鉴权信息
	errCode = eXosip_add_authentication_info(excontext, username, username, password, "MD5", NULL);
	if (errCode > 0)
	{
		// 添加鉴权信息失败
		return errCode;
	}

	// 接下来组织报文
	osip_message_t *reg = NULL;
	char proxy[1024] = { 0 };
	char from[1024] = { 0 };
	char contact[1024] = { 0 };
	sprintf(from, "sip:%s@%s:%d", username, server_ip, server_port);

	// 这里用来生成REGISTER后面跟着的一段内容，这里按照28181的文档说明，格式应当为： sip:SIP服务器编码@目的域名或IP地址端口 
	sprintf(proxy, "sip:%s@%s:%d", server_gbcode, server_ip, server_port);
	sprintf(contact, "sip:%s@%s:%d", local_gbcode_.c_str(), local_ip_.c_str(), local_port_);

	// 锁定
	eXosip_lock(excontext);
	int reg_id = eXosip_register_build_initial_register(excontext, from, proxy, contact, expired, &reg);
	if (reg_id < 0)
	{
		// 初始化鉴权信息失败
		eXosip_unlock(excontext);
		eXosip_quit(excontext);
		return -1;
	}

	errCode = eXosip_register_send_register(excontext, reg_id, reg);
	eXosip_unlock(excontext);
	if (errCode != 0)
	{
		// 发送注册失败！
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
	//	// 说明鉴权信息有问题，使用sip_event的rid重新产生一个注册包
	//	osip_message_t *reg = NULL;
	//	errCode = eXosip_register_build_register(excontext, sip_event->rid, reg_expired_, &reg);
	//	
	//	// 取回认证的字符串authorization
	//	osip_authorization_t *auth = NULL;
	//	char *strAuth = NULL;
	//	osip_message_get_authorization(reg, 0, &auth);
	//	osip_authorization_to_str(auth, &strAuth);

	//	//保存认证字符串
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
		// 没找到这个函数的定义，这里暂时注释掉看看
		//eXosip_automatic_refresh();
		eXosip_automatic_action(excontext);
		eXosip_unlock(excontext);

		if (sip_event == NULL)
			continue;

		switch (sip_event->type)
		{
		// 注册有关的事件
		case EXOSIP_REGISTRATION_SUCCESS:
			// user is registred
			// 注册成功
			OutputDebugStringA("WorkingThread >>> EXOSIP_REGISTRATION_SUCCESS\n");
			
			break;
		case EXOSIP_REGISTRATION_FAILURE:
			// user is not registred
			// 注册失败，判断是否是401错误
			OutputDebugStringA("WorkingThread >>> EXOSIP_REGISTRATION_FAILURE\n");
			client->ReregisterBy401(sip_event);
			break;

		/* INVITE related events within calls */
		// 在调用中与INVITE有关的事件
		case EXOSIP_CALL_INVITE:
			// announce a new call
			// 一个新的CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_INVITE\n");
			break;
		case EXOSIP_CALL_REINVITE:
			// announce a new INVITE within call
			// 在一个CALL里面又收到一个新的INVITE
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REINVITE\n");
			break;
		case EXOSIP_CALL_NOANSWER:
			// announce no answer within the timeout
			// 超时无响应
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_NOANSWER\n");
			break;
		case EXOSIP_CALL_PROCEEDING:
			// announce processing by a remote app
			// 远程应用正在处理这个CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_PROCEEDING\n");
			break;
		case EXOSIP_CALL_RINGING:
			// announce ringback
			// 正在响铃
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_RINGING\n");
			break;
		case EXOSIP_CALL_ANSWERED:
			// announce start of call
			// 开始一个CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_ANSWERED\n");
			break;
		case EXOSIP_CALL_REDIRECTED:
			// announce a redirection
			// 重定向一个CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REDIRECTED\n");
			break;
		case EXOSIP_CALL_REQUESTFAILURE:
			// announce a request failure
			// 呼叫请求失败
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_REQUESTFAILURE\n");
			break;
		case EXOSIP_CALL_SERVERFAILURE:
			// announce a server failure
			// 一个服务端错误
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_SERVERFAILURE\n");
			break;
		case EXOSIP_CALL_GLOBALFAILURE:
			// announce a global failure
			// 一个全局错误
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_GLOBALFAILURE\n");
			break;
		case EXOSIP_CALL_ACK:
			// ACK received for 200ok to INVITE
			// INVITE的200 OK已经收到ACK了
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_ACK\n");
			break;

		case EXOSIP_CALL_CANCELLED:
			// announce that call has been cancelled
			// 取消一个CALL
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_CANCELLED\n");
			break;

		/* request related events within calls (except INVITE) */
		// 调用中与请求有关的事件(INVITE除外)
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
			// 在这个CALL中收到一个BYE
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_CLOSED\n");
			break;

		/* for both UAS & UAC events */
		// UAS和UAC都会收到的事件
		case EXOSIP_CALL_RELEASED:
			// call context is cleared
			OutputDebugStringA("WorkingThread >>> EXOSIP_CALL_RELEASED\n");
			break;

		/* events received for request outside calls */
		// 接收请求外部调用的事件
		case EXOSIP_MESSAGE_NEW:
			// announce new incoming request.
			// 接收到一个消息
			// 基本上，上级平台所有的命令都会在这里收到（未测试点流的情况）
			// 为了保障程序效率，在这里收到请求后，直接扔到请求队列中，后面安排线程池进行处理
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_NEW\n");

			// 首先判断是否是一个MESSAGE
			if (MSG_IS_MESSAGE(sip_event->request))
			{
				osip_body_t *sip_body = NULL;
				osip_message_get_body(sip_event->request, 0, &sip_body);
				OutputDebugStringA(sip_body->body);
				OutputDebugStringA("\n");

				// 将收到的请求扔进队列，或缓存中
				// 在放进缓存队列前，确认好回复时需要哪些参数，组织一个结构体进行缓存
				// 剩下的事情就是起线程池，多线程处理信息了
				//client->request_queue_.push();

				//按照规则，回复200 OK信息
				osip_message_t *answer = NULL;
				eXosip_message_build_answer(excontext, sip_event->tid, 200, &answer);
				eXosip_message_send_answer(excontext, sip_event->tid, 200, answer);
			}
			break;
		case EXOSIP_MESSAGE_PROCEEDING:
			// announce a 1xx for request.
			// 接收到一个1XX的请求
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_PROCEEDING\n");
			break;
		case EXOSIP_MESSAGE_ANSWERED:
			// announce a 200ok
			// 接收到一个200 OK
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_ANSWERED\n");
			break;
		case EXOSIP_MESSAGE_REDIRECTED:
			// announce a failure.
			// 接收到一个失败信息
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_REDIRECTED\n");
			break;
		case EXOSIP_MESSAGE_REQUESTFAILURE:
			// announce a failure.
			// 接收到一个失败信息
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_REQUESTFAILURE\n");
			break;
		case EXOSIP_MESSAGE_SERVERFAILURE:
			// announce a failure.
			// 接收到一个失败信息
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_SERVERFAILURE\n");
			break;
		case EXOSIP_MESSAGE_GLOBALFAILURE:
			// announce a failure.
			// 接收到一个失败信息
			OutputDebugStringA("WorkingThread >>> EXOSIP_MESSAGE_GLOBALFAILURE\n");
			break;
					
		/* Presence and Instant Messaging */
		// 状态呈现与即时通信
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
			// 一个服务端错误
			OutputDebugStringA("WorkingThread >>> EXOSIP_NOTIFICATION_SERVERFAILURE\n");
			break;
		case EXOSIP_NOTIFICATION_GLOBALFAILURE:
			// announce a global failure
			// 一个全局错误
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