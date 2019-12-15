#ifndef _bg28181Client_H_
#define _bg28181Client_H_

#include <string>
#include <queue>
#include "bg28181Def.h"
#include "Poco/Thread.h"

class bg28181Client
{
public:
	bg28181Client();
	~bg28181Client();

public:
	/**
	 * ��ʼ���ͻ��ˣ���Ҫ�м�������
	 * ������
	 *	@local_ip		�ͻ���IP
	 *	@local_port		�ͻ��˶˿�
	 *	@local_gb_code	�ͻ��˹������
	 *	@net_type		�ͻ�����������
	 * ����ֵ��
	 *	@
	 */
	int Initialize(const char *local_ip, unsigned short local_port, const char *local_gb_code, NET_TYPE net_type);

	/**
	 * ע�ᵽ������
	 */
	int Register(const char *server_ip, unsigned short server_port, const char *server_gbcode, const char *username, const char *password, int expired);

	/**
	 * �ӷ�����ע��
	 */
	int Unregister();

private:
	int ReregisterBy401(void* event_info);

private:
	void *context_;

	std::string local_ip_;
	unsigned short local_port_;
	std::string local_gbcode_;

	std::string server_ip_;
	unsigned short server_port_;
	std::string server_gbcode_;
	std::string server_username_;
	std::string server_password_;
	int reg_expired_;

public:
	std::string auth_info_;

public:
	std::queue<std::string> request_queue_;

private:
	Poco::Thread working_thread_;
	bool need_working_thread_exit_;
	static void WorkingThread(void* lpParam);
};

#endif//_bg28181Client_H_
