#ifndef _bg28181Client_H_
#define _bg28181Client_H_

#include <string>
#include "bg28181Def.h"

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
	int Register(const char *server_ip, unsigned short server_port, unsigned char *server_gbcode, const char *username, const char *password, int expired);

private:
	void *context_;

	std::string local_ip_;
	unsigned short local_port_;
	std::string local_gbcode_;
};

#endif//_bg28181Client_H_
