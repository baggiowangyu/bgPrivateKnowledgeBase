#include <winsock.h>
#include "eXosip2/eXosip.h"
#include "osip2/osip.h"
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

	return errCode;
}

int bg28181Client::Register(const char *server_ip, unsigned short server_port, const char *server_gbcode, const char *username, const char *password, int expired)
{
	int errCode = 0;

	eXosip_t *excontext = (eXosip_t *)context_;

	// ���ȣ������Ȩ��Ϣ
	eXosip_clear_authentication_info(excontext);

	// ���������������Ӽ�Ȩ��Ϣ
	errCode = eXosip_add_authentication_info(excontext, username, username, password, NULL, NULL);
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

	// ����������һ���̣߳����ڽ��շ��������ص���Ϣ

	return errCode;
}