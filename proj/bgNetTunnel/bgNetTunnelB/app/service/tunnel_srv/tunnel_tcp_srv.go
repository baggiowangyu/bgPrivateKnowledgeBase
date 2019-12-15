/*
TCP隧道服务端
本隧道服务端是一个TCP服务端，持有一个隧道服务端观察者接口，用于向其反馈接收到的数据
本隧道服务端实现了数据发送接口，可以抽象为隧道服务端接口，供隧道服务端调用
 */
package tunnel_srv

import (
	"bgNetTunnelB/app/service/cli_mgr"
	"bgNetTunnelB/tunnel_protocol"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
)

type TunnelTcpServer struct {
	tcp_server *gtcp.Server
	tunnel_A_conn *gtcp.Conn
	is_running int
	tunnel_proto *tunnel_protocol.TunnelProto
	enable_tunnel_crypto bool
	client_mgr *cli_mgr.ClientMgr
}

func (t *TunnelTcpServer) Initialize(tp *tunnel_protocol.TunnelProto, climgr *cli_mgr.ClientMgr) {

	t.enable_tunnel_crypto = g.Config().GetBool("tunnel.enable_crypto")
	t.client_mgr = climgr
	_ = t.client_mgr.Initialize(tp, t, t.enable_tunnel_crypto)

	// 从配置文件读取出监听端口
	srv_address := g.Config().GetString("tunnel.tcp_address")
	t.tcp_server = g.TCPServer(srv_address)
	t.tcp_server.SetAddress(srv_address)
	t.tcp_server.SetHandler(t.ConnHandler)
	t.tunnel_proto = tp
	t.tunnel_A_conn = nil

	// 由于服务端在主线程启动会阻塞主线程，我们这里用一个协程来负责启动隧道服务端
	go t.Start()
}

func (t *TunnelTcpServer) Start() {

	err := t.tcp_server.Run()
	if err != nil {
		glog.Error(err)
	}
}

func (t *TunnelTcpServer) ConnHandler(client_conn *gtcp.Conn) {
	// 隧道客户端接入了
	glog.Debugf("TunnelTcpServer::ConnHandler() %s connect to tunnel server.", client_conn.RemoteAddr().String())
	if t.tunnel_A_conn == nil {
		t.tunnel_A_conn = client_conn
	} else {
		glog.Error("Already has a tunnel a connection...")
		_ = client_conn.Close()
		return
	}

	opt := gtcp.PkgOption{
		HeaderSize:4,
		MaxDataSize : 16777215,
	}

	for {
		// 等待隧道客户端发送数据，这里使用RecvPkg，防止隧道内发生粘包，导致解包失败,opt内定义了，单包最大4GB字节
		data, err := client_conn.RecvPkg(opt)
		if err != nil {
			// 发生异常了，那么我们清空连接信息，等待下一个隧道A端连过来
			t.tunnel_A_conn = nil
			glog.Debug("Recv tunnel data failed.")
			glog.Error(err)
			break
		} else {
			glog.Debugf("Recv tunnel data :\n%s", data)
		}

		// 接收到数据，首先反序列化，得到整包，判断是业务数据还是异常数据，是业务数据则扔给客户端管理器
		tunnel_sec_protocol_data, err := t.tunnel_proto.Unmarshal(data)
		if err != nil {
			glog.Error(err)
			continue
		}

		if tunnel_sec_protocol_data.Main == tunnel_protocol.MainType_Business {
			// 向客户端管理器发送业务数据
			t.client_mgr.RecvDataFromTunnel(tunnel_sec_protocol_data.Data)
		} else if tunnel_sec_protocol_data.Main == tunnel_protocol.MainType_Exception {
			// 处理异常，可能是A端客户端主动断开了，我们在这里一并通知客户端管理器，并发执行
			go t.client_mgr.RecvExceptionFromTunnel(tunnel_sec_protocol_data.Data)
		}
	}
}

func (t *TunnelTcpServer) PostDataToTunnel(data []byte) error {
	var err error

	opt := gtcp.PkgOption{
		HeaderSize : 4,
		MaxDataSize : 16777215,
	}

	// 直接找到隧道连接对象，发送出去
	err = t.tunnel_A_conn.SendPkg(data, opt)
	if err != nil {
		//glog.Debug("Send tunnel data failed :\n", data)
		glog.Error(err)
	} else {
		//glog.Debug("Send tunnel data :\n", data)
	}

	return err
}