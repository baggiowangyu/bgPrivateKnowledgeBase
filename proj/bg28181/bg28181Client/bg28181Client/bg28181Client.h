#ifndef _bg28181Client_H_
#define _bg28181Client_H_


class bg28181Client
{
public:
	bg28181Client();
	~bg28181Client();

public:
	/**
	 * ��ʼ���ͻ��ˣ���Ҫ�м�������
	 * ������
	 *	@
	 *	@
	 *	@
	 * ����ֵ��
	 */
	int Initialize(const char *local_ip, unsigned short local_port, int net_type);

private:
	void *context_;
};

#endif//_bg28181Client_H_
