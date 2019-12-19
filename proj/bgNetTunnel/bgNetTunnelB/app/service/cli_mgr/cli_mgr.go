/*
客户端管理器

用于管理所有虚拟客户端的生命周期
 */
package cli_mgr

import (
	"bgNetTunnelB/tunnel_protocol"
	"errors"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/glog"
)

type ClientObserverInterface interface {
	// 从通道拾取消息
	PostDataToTunnel(data []byte) error
}

type ClientInfo struct {
	address string
	nettype string
}

type ClientMgr struct {
	tcp_clients map[string]*gtcp.Conn
	udp_clients map[string]*gudp.Conn

	tunnel_proto *tunnel_protocol.TunnelProto
	tunnel_enable_crypto bool

	// 这里实际上保存的是具体的隧道服务端对象，例如TunnelTcpServer等
	client_observer_interface ClientObserverInterface
}

func (c *ClientMgr) Initialize(tp *tunnel_protocol.TunnelProto, inter ClientObserverInterface, enable_crypto bool) error {
	c.tcp_clients = make(map[string]*gtcp.Conn, 0)
	c.udp_clients = make(map[string]*gudp.Conn, 0)
	c.tunnel_proto = tp
	c.client_observer_interface = inter
	c.tunnel_enable_crypto = enable_crypto
	return nil
}

func (c *ClientMgr) FindTcpClient(client_tag string) (*gtcp.Conn, error) {
	client, exist := c.tcp_clients[client_tag]
	if !exist {
		return client, errors.New("Not Found.")
	} else {
		return client, nil
	}
}

func (c *ClientMgr) FindUdpClient(client_tag string) (*gudp.Conn, error) {
	client, exist := c.udp_clients[client_tag]
	if !exist {
		return client, errors.New("Not Found.")
	} else {
		return client, nil
	}
}

/**
获取所有客户端信息
 */
func (c *ClientMgr) GetAllClients() []ClientInfo {
	// 先搞tcp的客户端
	infos := make([]ClientInfo, 0)
	for _, tcp_cli := range c.tcp_clients {
		info := ClientInfo{
			address: tcp_cli.RemoteAddr().String(),
			nettype: "TCP",
		}

		infos = append(infos, info)
	}

	for _, udp_cli := range c.udp_clients {
		info := ClientInfo{
			address: udp_cli.RemoteAddr().String(),
			nettype: "UDP",
		}

		infos = append(infos, info)
	}

	return infos
}

func (c *ClientMgr) CreateClient(mapping_id int, src_client_address string, dst_server_address string, net_type string) error {
	// 首先创建到服务器的连接
	if net_type == "TCP" {
		dst_conn, err := gtcp.NewConn(dst_server_address)
		if err != nil {
			glog.Errorf("[ClientMgr::CreateClient] Create client and connect to %s error. %s", dst_server_address, err.Error())
			return err
		}

		// 启动接收线程
		go c.TcpClientRecvThread(mapping_id, src_client_address, dst_server_address, net_type, dst_conn)

		// 加入客户端表
		c.tcp_clients[src_client_address] = dst_conn

	} else if net_type == "UDP" {
		dst_conn, err := gudp.NewConn(dst_server_address)
		if err != nil {
			glog.Errorf("[ClientMgr::CreateClient] Create client and connect to %s error. %s", dst_server_address, err.Error())
			return err
		}

		// 启动接收线程
		go c.UdpClientRecvThread(mapping_id, src_client_address, dst_server_address, net_type, dst_conn)

		// 加入客户端表
		c.udp_clients[src_client_address] = dst_conn
	}

	return nil
}

func (c *ClientMgr) DestroyClient(src_client_address string, net_type string) {
	if net_type == "TCP" {
		dst_conn, exist := c.tcp_clients[src_client_address]
		if !exist {
			return
		}

		// 关闭连接
		_ = dst_conn.Close()
		// 从map中移除
		delete(c.tcp_clients, src_client_address)
		// 通知隧道A端，此模拟客户端已断开连接，那边也可以断开连接了

	} else if net_type == "UDP" {
		dst_conn, exist := c.udp_clients[src_client_address]
		if !exist {
			return
		}

		// 关闭连接
		_ = dst_conn.Close()
		// 从map中移除
		delete(c.udp_clients, src_client_address)
		// 通知隧道A端，此模拟客户端已断开连接，那边也可以断开连接了

	}
}

func (c *ClientMgr) AsynchronousRecvDataFromTunnel(data *tunnel_protocol.TunnelProtocol)  {
	var err error
	// 这里接收到业务数据，查找对应客户端(不存在的话就新建)，异步发送数据
	if data.Net == tunnel_protocol.NetType_TCP {
	Find_Tcp_Client:
		client_conn, exists := c.tcp_clients[data.SrcCliAddr]
		if !exists {
			err = c.CreateClient(int(data.MappingID), data.SrcCliAddr, data.DstSrvAddr, data.Net.String())
			if err != nil {
				glog.Debugf("[ClientMgr::AsynchronousRecvDataFromTunnel] Create client failed.")
				glog.Error(err)
				return
			}

			goto Find_Tcp_Client
		}

		err = client_conn.Send(data.Data)
	} else if data.Net == tunnel_protocol.NetType_UDP {
	Find_Udp_Client:
		client_conn, exists := c.udp_clients[data.SrcCliAddr]
		if !exists {
			err = c.CreateClient(int(data.MappingID), data.SrcCliAddr, data.DstSrvAddr, data.Net.String())
			if err != nil {
				glog.Debugf("[ClientMgr::AsynchronousRecvDataFromTunnel] Create client failed.")
				glog.Error(err)
				return
			}

			goto Find_Udp_Client
		}

		err = client_conn.Send(data.Data)
	}
}

func (c *ClientMgr) RecvDataFromTunnel(data *tunnel_protocol.TunnelProtocol) {
	// 为了提高效率，我们在这里调用异步发送接口
	go c.AsynchronousRecvDataFromTunnel(data)
}

func (c *ClientMgr) RecvExceptionFromTunnel(data *tunnel_protocol.TunnelProtocol) {
	// 这里主要是接收A端客户端主动断开的情况，我们在这里只需要找到对应的模拟客户端key，将其强行断开，并从列表中移除即可
	c.DestroyClient(data.SrcCliAddr, data.Net.String())
}

func (c *ClientMgr) TcpClientRecvThread(mapping_id int, src_client_address string, dst_server_address string, net_type string, dst_conn *gtcp.Conn) {
	// 在这里接收服务器返回数据的时候要留一个心眼，由于访问较大文件的时候，服务器会直接返回大于我们pkg的限额（15MB左右），我们在这里可以将数据拆分
	// 当服务器返回的数据大于4MB的时候，我们将数据以4MB的大小进行拆分，循环通过隧道发回给A端

	for {
		data_from_server, err := dst_conn.Recv(-1)

		// 连接出现问题，要么意外断开，要么某一方主动断开
		// 此时我们需要调用管理器的结束客户端方法，结束掉当前客户端
		if err != nil {
			c.DestroyClient(src_client_address, net_type)

			// 组装数据
			tunnel_sec_data, err := c.tunnel_proto.BuildDataV1_0(nil, int32(mapping_id), src_client_address,
				dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Exception,
				tunnel_protocol.SubType_Request, c.tunnel_enable_crypto)
			if err == nil {
				// 通知A端，断开客户端连接
				err = c.client_observer_interface.PostDataToTunnel(tunnel_sec_data)
				if err == nil {
					return
				}
			}

			glog.Error(err)
			return
		}

		data_from_server_len := len(data_from_server)
		group_len := 1024 * 1024 * 4
		send_count := data_from_server_len / group_len
		remain_data_len := data_from_server_len % group_len

		for index := 0; index < send_count; index++ {
			start_pos := index * group_len
			end_pos := (index + 1) * group_len
			send_data := data_from_server[start_pos:end_pos]

			// 组装数据
			tunnel_sec_data, err := c.tunnel_proto.BuildDataV1_0(send_data, int32(mapping_id), src_client_address,
				dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
				tunnel_protocol.SubType_Response, c.tunnel_enable_crypto)
			if err != nil {
				glog.Error(err)
				continue
			}

			// 通过接口扔给隧道
			err = c.client_observer_interface.PostDataToTunnel(tunnel_sec_data)
			if err != nil {
				glog.Error(err)
				continue
			}
		}

		if remain_data_len > 0 {
			start_pos := data_from_server_len - remain_data_len
			end_pos := data_from_server_len
			send_data := data_from_server[start_pos:end_pos]

			// 组装数据
			tunnel_sec_data, err := c.tunnel_proto.BuildDataV1_0(send_data, int32(mapping_id), src_client_address,
				dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
				tunnel_protocol.SubType_Response, c.tunnel_enable_crypto)
			if err != nil {
				glog.Error(err)
				continue
			}

			// 通过接口扔给隧道
			err = c.client_observer_interface.PostDataToTunnel(tunnel_sec_data)
			if err != nil {
				glog.Error(err)
				continue
			}
		}
	}
}

func (c *ClientMgr) UdpClientRecvThread(mapping_id int, src_client_address string, dst_server_address string, net_type string, dst_conn *gudp.Conn) {
	// 在这里接收服务器返回数据的时候要留一个心眼，由于访问较大文件的时候，服务器会直接返回大于我们pkg的限额（15MB左右），我们在这里可以将数据拆分
	// 当服务器返回的数据大于4MB的时候，我们将数据以4MB的大小进行拆分，循环通过隧道发回给A端
	for {
		data_from_server, err := dst_conn.Recv(-1)

		//tunnel_data := tunnel_protocol.TunnelProtocol{
		//	MappingID : int32(mapping_id),
		//	SrcCliAddr : src_client_address,
		//	DstSrvAddr : dst_server_address,
		//	Net : tunnel_protocol.NetType_UDP,
		//	Data : data_from_server,
		//}

		if err != nil {
			// 连接出现问题，要么意外断开，要么某一方主动断开
			// 此时我们需要调用管理器的结束客户端方法，结束掉当前客户端
			c.DestroyClient(src_client_address, net_type)

			tunnel_sec_data, err := c.tunnel_proto.BuildDataV1_0(data_from_server, int32(mapping_id), src_client_address,
				dst_server_address, tunnel_protocol.NetType_UDP, tunnel_protocol.MainType_Exception,
				tunnel_protocol.SubType_Request, c.tunnel_enable_crypto)
			//tunnel_sec_data, err := c.tunnel_proto.Marshal("1.0", tunnel_protocol.MainType_Exception, tunnel_protocol.SubType_Request, c.tunnel_enable_crypto, &tunnel_data)
			if err == nil {
				// 通知A端，断开客户端连接
				err = c.client_observer_interface.PostDataToTunnel(tunnel_sec_data)
				if err == nil {
					return
				}
			}

			glog.Error(err)
			return
		}

		tunnel_sec_data, err := c.tunnel_proto.BuildDataV1_0(data_from_server, int32(mapping_id), src_client_address,
			dst_server_address, tunnel_protocol.NetType_UDP, tunnel_protocol.MainType_Business,
			tunnel_protocol.SubType_Response, c.tunnel_enable_crypto)
		//tunnel_sec_data, err := c.tunnel_proto.Marshal("1.0", tunnel_protocol.MainType_Business, tunnel_protocol.SubType_Response, c.tunnel_enable_crypto, &tunnel_data)
		if err != nil {
			glog.Error(err)
			continue
		}

		// 通过接口扔给隧道
		err = c.client_observer_interface.PostDataToTunnel(tunnel_sec_data)
		if err != nil {
			glog.Error(err)
			continue
		}
	}
}