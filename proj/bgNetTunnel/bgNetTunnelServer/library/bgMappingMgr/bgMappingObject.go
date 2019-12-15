package bgMappingMgr

import (
	"bgNetTunnelServer/library/bgClientMgr"
	"bgNetTunnelServer/library/bgNetProtocol/bgNetMessage"
	"bgNetTunnelServer/library/bgNetTunnel"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/glog"
	"strconv"
)

/*

映射对象

1、这里是真实客户端数据入口

 */
type MappingObject struct {
	mapping_id			int
	tcp_server 			*gtcp.Server
	udp_server 			*gudp.Server
	net_type			string

	client_mgr			bgClientMgr.ClientMgr
	Net_tunnel_client	*bgNetTunnel.NetTunnelClient
}

func (m *MappingObject) TCPHandler(conn *gtcp.Conn) {
	client_addr := conn.Conn.RemoteAddr().String()
	mapping_addr := conn.Conn.LocalAddr().String()

	glog.Infof("Client connect >>> mapping address : %s; client address : %s", mapping_addr, client_addr)

	// 有消息进来，首先检查客户端表内是否有此客户端，如果没有则创建客户端对象
	_, err := m.client_mgr.FindClientObject(client_addr)
	if err != nil {
		m.client_mgr.AddClientObject(client_addr, mapping_addr, "TCP", conn, nil)
	}

	// 判断连接的映射地址，取得映射ID
	msg_v1 := bgNetMessage.NetMessageV1{}
	msg_v1.MsgType = 1
	msg_v1.MappingId = int32(m.mapping_id)
	msg_v1.ClientId = client_addr

	for {
		data, err := conn.Recv(-1)
		if err != nil {
			// 这里可能意味着设备已经断开了
			glog.Info(err)
			break
		}

		data_b64 := gbase64.EncodeToString(data)
		msg_v1.MessageBody = data_b64

		// 调用隧道接口发送数据，不能直接饮用service包，用chan吗？
		err = m.Net_tunnel_client.SendMessage(msg_v1)
	}
}

func (m *MappingObject) UDPHandler(conn *gudp.Conn) {
	client_addr := conn.RemoteAddr().String()
	mapping_addr := conn.LocalAddr().String()
	println("UDPHandler: mapping address : " + mapping_addr + "; client address : " + client_addr)

	// 有消息进来，首先检查客户端表内是否有此客户端，如果没有则创建客户端对象
	_, err := m.client_mgr.FindClientObject(client_addr)
	if err == nil {
		m.client_mgr.AddClientObject(client_addr, mapping_addr, "UDP", nil, conn)
	}
}

func (m *MappingObject) SendMsgToClient(client_addr string, data []byte) error {
	// 首先根据client_id找到对应的conn连接对象
	client_object, err := m.client_mgr.FindClientObject(client_addr)
	if err != nil {
		glog.Error(err)
		return err
	}

	// 直接发送
	if client_object.Conn_type == "TCP" {
		err = client_object.Tcp_conn.Send(data)
	} else if client_object.Conn_type == "UDP" {
		err = client_object.Udp_conn.Send(data)
	}

	if err != nil {
		glog.Error(err)
	}

	return err
}

func (m *MappingObject) Startup(mapping_id int, mapping_ip string, mapping_port int, net_type string, tunnel_client *bgNetTunnel.NetTunnelClient) error {
	var err error

	// 初始化客户端管理对象
	m.client_mgr.Initialize()
	m.Net_tunnel_client = tunnel_client

	address_string := mapping_ip + ":" + strconv.Itoa(mapping_port)
	if net_type == "TCP" {

		// 初始化TCP服务器
		m.tcp_server = g.TCPServer("tcp://" + address_string)
		m.tcp_server.SetAddress(address_string)
		m.tcp_server.SetHandler(m.TCPHandler)
		// 这里启动后会挂起，应该采用goruntime协程处理
		go m.RunTcpServer()
		//err = m.tcp_server.Run()

	} else {

		// 初始化UDP服务器
		m.udp_server = g.UDPServer("udp://" + address_string)
		m.udp_server.SetAddress(address_string)
		m.udp_server.SetHandler(m.UDPHandler)
		//err = m.udp_server.Run()
		go m.RunUdpServer()
	}

	m.mapping_id = mapping_id
	m.net_type = net_type
	return err
}

func (m *MappingObject) RunTcpServer()  {
	_ = m.tcp_server.Run()
}

func (m *MappingObject) RunUdpServer() {
	_ = m.udp_server.Run()
}

func (m *MappingObject) Shutdown() error {
	var err error

	if m.net_type == "TCP" {
		err = m.tcp_server.Close()
	} else {
		err = m.udp_server.Close()
	}

	return err
}
