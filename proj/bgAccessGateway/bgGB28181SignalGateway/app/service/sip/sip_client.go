package sip

import (
	"bgGB28181SignalGateway/app/service/data"
	"fmt"
	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/transaction"
	"strconv"
	"strings"
)

/*
SIP客户端，表示一个SIP设备
*/
type SipClient struct {

	client_usercode string	// 客户端设备编码
	client_ip string		// 客户端IP
	client_port uint16		// 客户端PORT
	client_protocol string	// 客户端接入协议
	client_sip_call_id string
	client_sip_cseq uint32
	client_domain string	// 客户端国标ID前10位

	server_host string
	server_port uint16
	server_gbcode string
	server_domain string	// 国标ID前10位

	// 2020-04-13 新增加的终端成员
	client_endpoint endpoint

	register_time int64		// 注册时间
	expired_time int64		// 有效期
	latest_heartbeat_time int64	// 最后一次心跳时间
	location_info data.GxxGmLocationInfo // 缓存最新的GPS数据
	exception_info data.GxxGmExceptionInfo
	base_info data.GxxGmDeviceBaseInfo

	task_queue *gqueue.Queue // 平台下发命令队列
}

func (client *SipClient) SendCatlogQueryRequest(server_transaction_mgr *transaction.Manager) error {
	client.client_sip_cseq += 1

	// 构建请求参数
	recipient := &base.SipUri{
		User : base.String{S: client.client_usercode},
		Host : client.client_domain,	// 这里应该是国标前10位
	}

	branch := "z9hG4bK" + strconv.Itoa(grand.Intn(9))
	tag := strconv.FormatInt(gtime.Now().Unix(), 10)
	call_id := fmt.Sprintf("%s@%s", grand.Str(10), client.client_ip)

	// 构建请求信息
	request_body := fmt.Sprintf("<?xml version=\"1.0\" encoding=\"gb2312\"?>\n<Query>\n\t" +
		"<CmdType>Catalog</CmdType>\n\t<SN>1</SN>\n\t<DeviceID>%s</DeviceID>\n</Query>\n", client.client_usercode)

	request := base.NewRequest("MESSAGE", recipient,"SIP/2.0",
		[]base.SipHeader{
			Via("UDP", client.server_host, client.server_port, branch), // 这里的IP端口应该是服务端的
			From(client.server_gbcode, client.server_domain, "UDP", tag),
			To(client.client_usercode, client.client_domain, tag),
			//Contact(client.client_usercode, client.client_ip),
			CSeq(client.client_sip_cseq, "MESSAGE"),
			&base.GenericHeader{
				HeaderName: "Content-Type",
				Contents: "Application/MANSCDP+xml",
			},
			&base.GenericHeader{
				HeaderName: "User-Agent",
				Contents: "GxxGm",
			},
			CallId(call_id),
			ContentLength(uint32(len(request_body))),
		}, request_body)

	// 这里发下去就可以了，数据会从服务器统一入口发回来
	// 这里发送是OK了，但是还不是很清楚
	client_transaction := server_transaction_mgr.Send(request, fmt.Sprintf("%v:%v", client.client_ip, client.client_port))
	response := <- client_transaction.Responses()
	glog.Debug(response.GetBody())

	return nil
}

/*
SIP请求处理
*/
func (client *SipClient)RequestHandler(server_transaction *transaction.ServerTransaction, request *base.Request) {
	switch {
	case strings.EqualFold(string(request.Method), "INVITE"):
		//	点流
		client.HandleINVITE(server_transaction, request)
	case strings.EqualFold(string(request.Method), "ACK"):
		// 应答
		client.HandleACK(server_transaction, request)
	case strings.EqualFold(string(request.Method), "CANCEL"):
		// 取消
		client.HandleCANCEL(server_transaction, request)
	case strings.EqualFold(string(request.Method), "BYE"):
		// 结束点流
		client.HandleBYE(server_transaction, request)
	case strings.EqualFold(string(request.Method), "REGISTER"):
		//	处理注册相关的内容
		client.HandleREGISTER(server_transaction, request)
	case strings.EqualFold(string(request.Method), "OPTIONS"):
		//
		client.HandleOPTIONS(server_transaction, request)
	case strings.EqualFold(string(request.Method), "SUBSCRIBE"):
		// 订阅
		client.HandleSUBSCRIBE(server_transaction, request)
	case strings.EqualFold(string(request.Method), "NOTIFY"):
		// 通知
		client.HandleNOTIFY(server_transaction, request)
	case strings.EqualFold(string(request.Method), "REFER"):
		// 引用
		client.HandleREFER(server_transaction, request)
	default:
		// 常规消息，基本上90%的通信都是此类通信
		client.HandleMESSAGE(server_transaction, request)
	}
}

/*
处理INVITE请求
参数：
@request 	sip请求
*/
func (client *SipClient) HandleINVITE(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleACK(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleCANCEL(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleBYE(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleREGISTER(server_transaction *transaction.ServerTransaction, request *base.Request) {
	// 这里是处理设备异常断线后重新接入后的注册流程处理。
	// 根据配置来决定是否使用鉴权。
	// 评估一下是否清空基本信息、定位信息、异常信息。
	// 需要重新计算注册是否过期
	/*client_method*/_, client_host, client_port, /*client_protocol*/_, client_usercode := ParseRequest(request)

	var err error
	// 1、检查当前服务器是否开启了鉴权
	if Sip_enable_Authenticate {
		glog.Debug(request)

		// 2、若开启了鉴权，则判断请求中是否携带鉴权信息
		sip_headers := request.Headers("Authorization")
		if len(sip_headers) == 0 {
			// 需要返回401
			err := Response401Unauthorized(server_transaction, client_usercode, client_host, client_port)
			if err != nil {
				glog.Error(err)
				return
			}
		}

		for _, sip_header := range sip_headers {
			glog.Debug(sip_header)
		}
		// 3、若未携带，则返回401，若携带了则校验鉴权信息
		// 4、鉴权信息校验失败，则返回401，否则返回200 OK
	} else {
		client.expired_time, err = ReadRequestHeaderExpiredValue(request)
		if err != nil {
			glog.Error(err)
			// 这里应该返回一个其他错误的
			return
		}
		client.register_time = gtime.Now().Unix()
		client.latest_heartbeat_time = gtime.Now().Unix()

		// 返回200 OK
		err := Response200OK(server_transaction, client_usercode, client_host, client_port)
		if err != nil {
			glog.Error(err)
			return
		}

		// 根据配置，决定是否发送设备目录查询命令

	}

	return
}

func (client *SipClient) HandleOPTIONS(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleSUBSCRIBE(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleNOTIFY(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleREFER(server_transaction *transaction.ServerTransaction, request *base.Request) {

}

func (client *SipClient) HandleMESSAGE(server_transaction *transaction.ServerTransaction, request *base.Request) {
	var xml_data string
	glog.Debug(request.String())

	request_data := request.Body

	// 这里想办法挖掉不可见字符，例如\n、\t
	xml_data = strings.ReplaceAll(request_data, "\n", "")

	glog.Debugf("[[[MESSAGE]]] INFO:\n%s", xml_data)

	request_data_json := gjson.New(xml_data)
	if request_data_json.Contains("Notify") {
		notify_json := request_data_json.GetJson("Notify")
		cmd_type := notify_json.GetString("CmdType")

		if strings.Compare(cmd_type, "Keepalive") == 0 {
			// 心跳数据
			client.HandleGB28181StandardKeepalive(server_transaction)
		} else if strings.Compare(cmd_type, "TransData") == 0 {
			// 透传数据
			info := notify_json.GetString("Info")
			client.HandleExtraInfo(server_transaction, info)
		}
	} else if request_data_json.Contains("Response") {
		// 这里是平台下发命令后的应答数据
		response_json := request_data_json.GetJson("Response")
		cmd_type := response_json.GetString("CmdType")

		if strings.Compare(cmd_type, "Catalog") == 0 {
			// 设备目录
			device_list_json := request_data_json.GetJson("DeviceList")
			client.HandleQueryCatalogResponseInfo(server_transaction, device_list_json)
		}
	}
}

func (client *SipClient) HandleGB28181StandardKeepalive(server_transaction *transaction.ServerTransaction) {
	glog.Debugf("接收到[%s]发送的心跳信息", client.client_usercode)

	// 更新客户端最后一次心跳时间
	client.latest_heartbeat_time = gtime.Now().Unix()

	// 响应200 OK，好像没什么用
	err := Response200OK(server_transaction, client.client_usercode, client.client_ip, client.client_port)
	if err != nil {
		glog.Error(err)
		return
	}
}

func (client *SipClient) HandleExtraInfo(server_transaction *transaction.ServerTransaction, info string) {
	var origin_data_handled string

	// 优先回复200 OK，已收到数据
	err := Response200OK(server_transaction, client.client_usercode, client.client_ip, client.client_port)
	if err != nil {
		glog.Error(err)
	}

	// base64解码
	origin_data, err := gbase64.DecodeToString(info)
	if err != nil {
		glog.Error(err)
		return
	}

	////////////////////////////////////////////////////////////////////////////////////////
	// 处理一下多余字符，并将其组成完整的XML
	for {
		pos := strings.Index(origin_data, "> ")
		if pos >= 0 {
			origin_data = strings.ReplaceAll(origin_data, "> ", ">")
			continue
		} else {
			break
		}
	}

	for {
		pos := strings.Index(origin_data, " <")
		if pos >= 0 {
			origin_data = strings.ReplaceAll(origin_data, " <", "<")
			continue
		} else {
			break
		}
	}

	origin_data_handled = strings.ReplaceAll(origin_data, "\n", "")
	origin_data_handled = strings.ReplaceAll(origin_data, "\t", "")
	origin_data_handled = fmt.Sprintf("<?xml version=\"1.0\" encoding=\"GB2312\"?><Info>%s</Info>", origin_data_handled)
	// xml组装完成
	////////////////////////////////////////////////////////////////////////////////////////

	// 这里都是我们的扩展协议
	trans_data_json := gjson.New(origin_data_handled)

	sub_cmd_type := trans_data_json.GetJson("Info").GetString("SubCmdType")
	if sub_cmd_type == "LocationInfo" {
		// 扩展协议：定位信息
		location_json := trans_data_json.GetJson("Info").GetMap("Location")
		client.HandleExtraLocationInfo(location_json)

	} else if sub_cmd_type == "DeviceException" {
		// 扩展协议：设备异常信息
		exception_map := trans_data_json.GetJson("Info").GetMap("Exceptions")
		client.HandleExtraExceptionInfo(exception_map)

	} else if sub_cmd_type == "DeviceInfo" {
		// 扩展协议：设备基础信息
		base_info_map := trans_data_json.GetJson("Info").GetMap("DeviceStates")
		client.HandleExtraBaseInfo(base_info_map)

	}
}

func (client *SipClient) HandleQueryCatalogResponseInfo(server_transaction *transaction.ServerTransaction, device_list_json *gjson.Json) {
	// 响应200 OK，好像没什么用
	err := Response200OK(server_transaction, client.client_usercode, client.client_ip, client.client_port)
	if err != nil {
		glog.Error(err)
		return
	}

	// 这里处理设备列表
	glog.Debug(device_list_json.Export())
}

func (client *SipClient) HandleExtraLocationInfo(info map[string]interface{}) {
	// 扩展定位信息，保存到当前客户端缓存中，这里未判断字段是否存在
	client.location_info.Usercode = client.client_usercode
	_, exist := info["LocationTime"]
	if exist {
		client.location_info.Timevalue = info["LocationTime"].(string)
	}

	_, exist = info["Longitude"]
	if exist {
		client.location_info.Timevalue = info["Longitude"].(string)
	}

	_, exist = info["Latitude"]
	if exist {
		client.location_info.Timevalue = info["Latitude"].(string)
	}

	_, exist = info["Speed"]
	if exist {
		client.location_info.Timevalue = info["Speed"].(string)
	}

	_, exist = info["Direction"]
	if exist {
		client.location_info.Timevalue = info["Direction"].(string)
	}

	_, exist = info["Height"]
	if exist {
		client.location_info.Timevalue = info["Height"].(string)
	}

	_, exist = info["Radius"]
	if exist {
		client.location_info.Timevalue = info["Radius"].(string)
	}
}

func (client *SipClient) HandleExtraExceptionInfo(info map[string]interface{}) {
	client.exception_info.Usercode = client.client_usercode

	_, exist := info["Storage"]
	if exist {
		client.exception_info.Storage = info["Storage"].(string)
	}

	_, exist = info["Battery"]
	if exist {
		client.exception_info.Storage = info["Battery"].(string)
	}

	_, exist = info["CCD"]
	if exist {
		client.exception_info.Storage = info["CCD"].(string)
	}

	_, exist = info["MIC"]
	if exist {
		client.exception_info.Storage = info["MIC"].(string)
	}

	_, exist = info["POSITION"]
	if exist {
		client.exception_info.Storage = info["POSITION"].(string)
	}
}

func (client *SipClient) HandleExtraBaseInfo(info map[string]interface{}) {
	client.base_info.Usercode = client.client_usercode

	_, exist := info["Carrieroperator"]
	if exist {
		client.base_info.Carrieroperator = info["Carrieroperator"].(string)
	}

	_, exist = info["Nettype"]
	if exist {
		client.base_info.Carrieroperator = info["Nettype"].(string)
	}

	_, exist = info["Signal"]
	if exist {
		client.base_info.Carrieroperator = info["Signal"].(string)
	}

	_, exist = info["Battery"]
	if exist {
		client.base_info.Carrieroperator = info["Battery"].(string)
	}

	_, exist = info["Storage"]
	if exist {
		client.base_info.Carrieroperator = info["Storage"].(string)
	}

	_, exist = info["CPU"]
	if exist {
		client.base_info.Carrieroperator = info["CPU"].(string)
	}

	_, exist = info["Version"]
	if exist {
		client.base_info.Carrieroperator = info["Version"].(string)
	}

	_, exist = info["LocalRecord"]
	if exist {
		client.base_info.Carrieroperator = info["LocalRecord"].(string)
	}

	_, exist = info["ChargeState"]
	if exist {
		client.base_info.Carrieroperator = info["ChargeState"].(string)
	}
}