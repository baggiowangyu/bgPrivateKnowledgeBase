package sip

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/os/gtime"
	"github.com/stefankopieczek/gossip/base"
	"github.com/stefankopieczek/gossip/transaction"
	"github.com/stefankopieczek/gossip/transport"
	"os"
)

/*
SIP服务端对象
持有一个sip客户端清单
持有一个配置参数：是否开启鉴权
持有一个SIP传输管理器
持有一个SIP事务处理器
*/
type SipServer struct {
	Sip_server_host string
	Sip_server_port uint16
	Sip_server_gbcode string
	Sip_server_username string
	Sip_server_password string
	Sip_server_domain string

	// SIP客户端清单
	Sip_clients map[string]*SipClient
	Clients_map_lock *gmlock.Locker

	// SIP传输管理器
	Sip_transport_manager transport.Manager
	// SIP事务管理器
	Sip_transaction_manager *transaction.Manager

	// 鉴权开关
	Sip_enable_Authenticate bool
}

var sip_server *SipServer = new(SipServer)

func init() {
	// 创建一个SIP服务，并启动
	err := sip_server.Startup()
	if err != nil {
		os.Exit(4)
	}
}

func (server *SipServer) LockMap() {
	//server.Clients_map_lock.Lock("Clients_map")
}

func (server *SipServer) UnlockMap() {
	//server.Clients_map_lock.Unlock("Clients_map")
}

func (server *SipServer)Startup() error {
	var err error

	// 初始化客户端清单空间
	server.Sip_clients = make(map[string]*SipClient, 0)
	server.Clients_map_lock = new(gmlock.Locker)

	server.Sip_server_host = g.Config().GetString("GB28181.host")
	server.Sip_server_port = g.Config().GetUint16("GB28181.port")
	server.Sip_server_gbcode = g.Config().GetString("GB28181.gbcode")
	server.Sip_server_username = g.Config().GetString("GB28181.username")
	server.Sip_server_password = g.Config().GetString("GB28181.password")
	server.Sip_server_domain = server.Sip_server_gbcode[:10]

	// 初始化传输管理器
	server.Sip_transport_manager, err = transport.NewManager("udp")
	if err != nil {
		glog.Error(err)
		return err
	}

	// 初始化事务管理器
	host_string := fmt.Sprintf("%s:%d", server.Sip_server_host, server.Sip_server_port)
	server.Sip_transaction_manager, err = transaction.NewManager(server.Sip_transport_manager, host_string)
	if err != nil {
		glog.Error(err)
		return err
	}

	// 配置服务器鉴权信息
	Sip_enable_Authenticate = false

	// 启动一个线程，用于获取请求信息
	go server.SipRequestRecv()
	return err
}

func (server *SipServer)SipRequestRecv() {
	for {
		// 获取一个SIP请求
		server_transaction := <- server.Sip_transaction_manager.Requests()

		// 解析请求，获取来源Host信息
		request := server_transaction.Origin()
		client_method, /*client_host*/_, /*client_port*/_, /*client_protocol*/_, client_usercode := ParseRequest(request)

		// 查找发来请求的客户端，这里检查也需要考虑是否加锁
		server.LockMap()
		sip_client, exists := server.Sip_clients[client_usercode]
		server.UnlockMap()
		if !exists {
			// 客户端不存在
			if client_method == base.REGISTER {
				// 请求为Register，尝试注册
				go server.HandleRegisterAsyn(server_transaction, request)
			} else {
				// 啥都不做，看看有没有什么办法可以断开这个连接
				server_transaction.Delete()
				server_transaction.Destination()
			}
		} else {
			// 客户端存在，则交给客户端处理，异步的
			go sip_client.RequestHandler(server_transaction, request)
		}
	}
}

func (server *SipServer)HandleRegisterAsyn(server_transaction *transaction.ServerTransaction, request *base.Request) {

	/*client_method*/_, client_host, client_port, client_protocol, client_usercode := ParseRequest(request)

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
		sip_client := &SipClient{
			client_ip : client_host,
			client_port : client_port,
			client_protocol : client_protocol,
			client_usercode : client_usercode,
			client_sip_cseq: 0,
			client_domain: client_usercode[:10],

			server_host: server.Sip_server_host,
			server_port: server.Sip_server_port,
			server_gbcode: server.Sip_server_gbcode,
			server_domain: server.Sip_server_domain,
		}

		server.LockMap()
		server.Sip_clients[client_usercode] = sip_client
		server.UnlockMap()

		var err error
		sip_client.expired_time, err = ReadRequestHeaderExpiredValue(request)
		sip_client.register_time = gtime.Now().Unix()
		sip_client.latest_heartbeat_time = gtime.Now().Unix()

		// 返回200 OK
		err = Response200OK(server_transaction, client_usercode, client_host, client_port)
		if err != nil {
			glog.Error(err)
			return
		}

		// 根据配置，决定是否发送设备目录查询命令
		sip_client.SendCatlogQueryRequest(server.Sip_transaction_manager)
		//client_catlog_response := RequestQueryCatlog(server.Sip_transaction_manager, client_usercode, client_host, client_port)
		//glog.Debug(client_catlog_response.Body)
	}

	return
}
