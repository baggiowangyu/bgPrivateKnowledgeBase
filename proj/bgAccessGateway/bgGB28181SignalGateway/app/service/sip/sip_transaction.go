package sip

import (
	"errors"
	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/transaction"
	"strconv"
)


/*
解析SIP请求
依次返回以下信息：
- SIP请求method,
- 客户端IP
- 客户端端口号
- 客户端采用协议
- 客户端用户ID（可以理解为国标ID）
*/
func ParseRequest(request *base.Request) (base.Method, string, uint16, string, string) {

	var method base.Method = base.BYE
	var host string = ""
	var port uint16 = 0
	var protocol string = ""
	var user_code string = ""

	method = request.Method

	// 从Header中拿到源IP地址，端口号，传输协议
	sip_headers := request.Headers("Via")
	for _, sip_header := range sip_headers {
		via_header := sip_header.(*base.ViaHeader)
		for _, via_hop := range *via_header {
			host = via_hop.Host
			port = *via_hop.Port
			protocol = via_hop.ProtocolName
			break
		}
	}

	// 从From.Address.User中拿到国标编码（或者叫用户ID）
	sip_headers = request.Headers("From")
	for _, sip_header := range sip_headers {
		from_header := sip_header.(*base.FromHeader)
		address_sipuri := from_header.Address.(*base.SipUri)
		user_code = address_sipuri.User.(base.String).S

	}

	return method, host, port, protocol, user_code
}

func ReadRequestHeaderExpiredValue(request *base.Request) (int64, error) {
	sip_headers := request.Headers("Expires")
	if len(sip_headers) <= 0 {
		return 0, errors.New("No such key")
	}

	var expired_time int64
	var err error
	for _, sip_header := range sip_headers {
		println(sip_header)
		generic_header := sip_header.(*base.GenericHeader)
		expired_time, err = strconv.ParseInt(generic_header.Contents, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	return expired_time, nil
}

/*
响应200OK
*/
func Response200OK(server_transaction *transaction.ServerTransaction, usercode string, host string, port uint16) error {
	var err error

	request := server_transaction.Origin()

	// 返回200 OK，并且撞见
	response := base.NewResponse(
		"SIP/2.0",
		200,
		"OK",
		[]base.SipHeader{},
		"",
	)

	base.CopyHeaders("Via", request, response)
	base.CopyHeaders("From", request, response)
	base.CopyHeaders("To", request, response)
	base.CopyHeaders("Call-Id", request, response)
	base.CopyHeaders("CSeq", request, response)
	response.AddHeader(
		&base.ContactHeader{
			DisplayName: nil,
			Address:     &base.SipUri{
				IsEncrypted: false,
				User:        base.String{S:usercode},
				Password:    nil,
				Host:        host,
				Port:        &port,
				UriParams:   nil,
				Headers:     nil,
			},
			Params:      nil,
		})

	server_transaction.Respond(response)
	server_transaction.Ack()

	return err
}

func Response401Unauthorized(server_transaction *transaction.ServerTransaction, usercode string, host string, port uint16) error {
	var err error

	request := server_transaction.Origin()

	// 返回200 OK，并且撞见
	response := base.NewResponse(
		"SIP/2.0",
		401,
		"Unauthorized",
		[]base.SipHeader{},
		"",
	)

	base.CopyHeaders("Via", request, response)
	base.CopyHeaders("From", request, response)
	base.CopyHeaders("To", request, response)
	base.CopyHeaders("Call-Id", request, response)
	base.CopyHeaders("CSeq", request, response)
	response.AddHeader(
		&base.ContactHeader{
			DisplayName: nil,
			Address:     &base.SipUri{
				IsEncrypted: false,
				User:        base.String{S:usercode},
				Password:    nil,
				Host:        host,
				Port:        &port,
				UriParams:   nil,
				Headers:     nil,
			},
			Params:      nil,
		})

	server_transaction.Respond(response)
	server_transaction.Ack()

	return err
}
