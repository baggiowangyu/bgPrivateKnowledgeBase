package sip

import "github.com/stefankopieczek/gossip/base"

// Utility methods for creating headers.

func Via(transport string, host string, port uint16, branch string) *base.ViaHeader {
	params := base.NewParams()
	params.Add("branch", base.String{S: branch})
	params.Add("rport", nil)

	return &base.ViaHeader{
		&base.ViaHop{
			ProtocolName:    "SIP",
			ProtocolVersion: "2.0",
			Transport:       transport,
			Host:            host,
			Port:            &port,
			Params:			 params,
		},
	}
}

func To(client_gbcode string, host string, tag string) *base.ToHeader {
	params := base.NewParams()
	if tag != "" {
		params.Add("tag", base.String{S: tag})
	}

	header := &base.ToHeader{
		//DisplayName: base.String{S: client_gbcode},
		Address: &base.SipUri{
			User: base.String{S: client_gbcode},
			Host: host,
			UriParams: base.NewParams(),
		},
		Params: params,
	}

	return header
}

func From(client_gbcode string, host string, transport string, tag string) *base.FromHeader {
	params := base.NewParams()
	if tag != "" {
		params.Add("tag", base.String{S: tag})
	}

	uri_params := base.NewParams()
	uri_params.Add("transport", base.String{S: transport})

	header := &base.FromHeader{
		//DisplayName: base.String{S: client_gbcode},
		Address: &base.SipUri{
			User: base.String{S: client_gbcode},
			Host: host,
			//UriParams: uri_params,
		},
		Params: params,
	}

	return header
}

func Contact(client_gbcode string, host string) *base.ContactHeader {
	return &base.ContactHeader{
		//DisplayName: base.String{S: client_gbcode},
		Address: &base.SipUri{
			User: base.String{S: client_gbcode},
			Host: host,
		},
	}
}

func CSeq(seqno uint32, method base.Method) *base.CSeq {
	return &base.CSeq{
		SeqNo:      seqno,
		MethodName: method,
	}
}

func CallId(callid string) *base.CallId {
	header := base.CallId(callid)
	return &header
}

func ContentLength(l uint32) base.ContentLength {
	return base.ContentLength(l)
}