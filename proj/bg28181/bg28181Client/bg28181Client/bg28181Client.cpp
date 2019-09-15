#include "eXosip2/eXosip.h"
#include "bg28181Client.h"


bg28181Client::bg28181Client()
{

}

bg28181Client::~bg28181Client()
{

}

int bg28181Client::Initialize(const char *local_ip, unsigned short local_port, int net_type)
{
	int errCode = 0;

	eXosip_t *excontext = (eXosip_t *)context_;

	errCode = eXosip_init(excontext);
	if (errCode <= 0)
	{
		return errCode;
	}

	return errCode;
}