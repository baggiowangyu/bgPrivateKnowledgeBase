/*
映射对象
映射对象负责提供客户端接入、客户端数据转发与应答等工作
映射对象持有完整的映射信息，包括但不限于：映射端口、源IP、源端口、网络类型、映射状态
 */
package mapping_srv

import (
	"bgNetTunnelA/app/service/tunnel_cli"
	"bgNetTunnelA/tunnel_protocol"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/gqueue"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/glog"
	"strconv"
)

type MappingBaseInfo struct {
	Mapping_id	int
	Mapping_port int
	Source_ip string
	Source_port int
	Net_type string
	Is_running int
}

type MappingObject struct {
	Info *MappingBaseInfo

	Tcp_server *gtcp.Server
	Tcp_client_table map[string]*gtcp.Conn

	Udp_server *gudp.Server
	Udp_client_table map[string]*gudp.Conn

	Tunnel_client_interface tunnel_cli.TunnelClientInterface
	Enable_crypto bool

	Response_queue *gqueue.Queue
	TunnelProtocolObject *tunnel_protocol.TunnelProto
}

func (m *MappingObject) Initialize(base_info *MappingBaseInfo, tunnel_client_inter tunnel_cli.TunnelClientInterface, t *tunnel_protocol.TunnelProto, enable_crypto bool) error {
	m.Info = base_info
	m.Tunnel_client_interface = tunnel_client_inter
	m.TunnelProtocolObject = t
	m.Enable_crypto = enable_crypto

	// 根据网络类型，初始化服务端
	server_address := fmt.Sprintf("0.0.0.0:%d", m.Info.Mapping_port)
	if m.Info.Net_type == "TCP" {
		m.Tcp_server = g.TCPServer(server_address)
		m.Tcp_server.SetAddress(server_address)
		m.Tcp_server.SetHandler(m.TcpConnectHandler)
		m.Tcp_client_table = make(map[string]*gtcp.Conn, 64)

	} else if m.Info.Net_type == "UDP" {
		m.Udp_server.SetAddress(server_address)
		m.Udp_server.SetHandler(m.UdpConnectHandler)
		m.Udp_client_table = make(map[string]*gudp.Conn, 64)

	}

	return nil
}

func (m *MappingObject) Start() {
	go m.MainThread()
}

func (m *MappingObject) MainThread() {
	var err error
	if m.Info.Net_type == "TCP" {
		err = m.Tcp_server.Run()
	} else if m.Info.Net_type == "UDP" {
		err = m.Udp_server.Run()
	} else {
		err = errors.New("Unsupported net type")
	}

	if err != nil {
		glog.Error(err)
		return
	}
}

func (m *MappingObject) SendDataToClient(data *tunnel_protocol.TunnelProtocol) error {
	var err error

	// 这里收到的才是真实的数据，想办法加到一个队列里面去吧，另外再搞一个线程扫描队列来逐一回复，否则会因为某一条网络通信阻塞而导致所有请求阻塞
	if data.Net == tunnel_protocol.NetType_TCP {
		conn, exists := m.Tcp_client_table[data.SrcCliAddr]
		if exists {
			go func() {
				err = conn.Send(data.Data)
				if err != nil {
					glog.Debug("[MappingObject::SendDataToClient] Send data to client failed.")
					glog.Error(err)
				} else {
					glog.Debug("[MappingObject::SendDataToClient] Send data to client succeed.")
				}
			}()

		}
	} else if data.Net == tunnel_protocol.NetType_UDP {
		conn, exists := m.Udp_client_table[data.SrcCliAddr]
		if exists {
			go func() {
				err = conn.Send(data.Data)
				if err != nil {
					glog.Debug("[MappingObject::SendDataToClient] Send data to client failed.")
					glog.Error(err)
				} else {
					glog.Debug("[MappingObject::SendDataToClient] Send data to client succeed.")
				}
			}()
		}
	}

	return err
}

func (m *MappingObject) TcpConnectHandler(client_conn *gtcp.Conn) {
	client_address := client_conn.RemoteAddr().String()
	glog.Infof("[MappingObject::TcpConnectHandler] Client connected. Address : %s", client_address)

	// 通知B端，有客户端接入，然后加入客户端表
	m.Tcp_client_table[client_address] = client_conn
	mapping_id := m.Info.Mapping_id
	dst_server_address := m.Info.Source_ip + ":" + strconv.Itoa(m.Info.Source_port)

	for {
		// 由于在后面使用SendPkg和RecvPkg有包大小限制，似乎是65535，我们这里限制一下每次接受的缓冲区大小为40K，后面可以适当调整
		data, err := client_conn.Recv(-1)

		if err != nil {
			glog.Debugf("[MappingObject::TcpConnectHandler] Port %d recv data from %s failed.", m.Info.Mapping_port, client_address)

			// 出现错误，通知B端，客户端已断开
			mashal_data, err := m.TunnelProtocolObject.BuildDataV1_0(nil, int32(mapping_id), client_address,
				dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Exception,
				tunnel_protocol.SubType_Request, m.Enable_crypto)
			if err != nil {
				glog.Debugf("[MappingObject::TcpConnectHandler] Build exception Marshaled TunnelProtocolObject failed.")
				glog.Error(err)
				break
			}

			// 向B端发送协议数据
			err = m.Tunnel_client_interface.PostDataToTunnel(mashal_data)
			if err != nil {
				glog.Debugf("[MappingObject::TcpConnectHandler] Post exception data to tunnel-B failed")
				glog.Error(err)
			}

			break
		} else {
			data_len := len(data)
			group_len := 1024 * 1024 * 4
			send_count := data_len / group_len
			remain_data_len := data_len % group_len
			glog.Debugf("[MappingObject::TcpConnectHandler] Port %d recv data from %s. length : %d", m.Info.Mapping_port, client_address, data_len)

			for index := 0; index < send_count; index++ {
				start_pos := index * group_len
				end_pos := (index + 1) * group_len
				send_data := data[start_pos:end_pos]

				// 组装数据
				tunnel_sec_data, err := m.TunnelProtocolObject.BuildDataV1_0(send_data, int32(mapping_id), client_address,
					dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
					tunnel_protocol.SubType_Response, m.Enable_crypto)
				if err != nil {
					glog.Debugf("[MappingObject::TcpConnectHandler] Build business Marshaled TunnelProtocolObject failed.")
					glog.Error(err)
					continue
				}

				// 通过接口扔给隧道
				err = m.Tunnel_client_interface.PostDataToTunnel(tunnel_sec_data)
				if err != nil {
					glog.Debugf("[MappingObject::TcpConnectHandler] Post business data to tunnel-B failed")
					glog.Error(err)
					continue
				}
			}

			if remain_data_len > 0 {
				start_pos := data_len - remain_data_len
				end_pos := data_len
				send_data := data[start_pos:end_pos]

				// 组装数据
				tunnel_sec_data, err := m.TunnelProtocolObject.BuildDataV1_0(send_data, int32(mapping_id), client_address,
					dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
					tunnel_protocol.SubType_Response, m.Enable_crypto)
				if err != nil {
					glog.Debugf("[MappingObject::TcpConnectHandler] Build last business Marshaled TunnelProtocolObject failed.")
					glog.Error(err)
					continue
				}

				// 通过接口扔给隧道
				err = m.Tunnel_client_interface.PostDataToTunnel(tunnel_sec_data)
				if err != nil {
					glog.Debugf("[MappingObject::TcpConnectHandler] Post last business data to tunnel-B failed")
					glog.Error(err)
					continue
				}
			}
		}
	}

	// 将客户端从客户端集合中移除
	delete(m.Tcp_client_table, client_address)
}

func (m *MappingObject) UdpConnectHandler(client_conn *gudp.Conn) {
	client_address := client_conn.RemoteAddr().String()

	// 通知B端，有客户端接入，然后加入客户端表
	m.Udp_client_table[client_address] = client_conn
	mapping_id := m.Info.Mapping_id
	dst_server_address := m.Info.Source_ip + ":" + strconv.Itoa(m.Info.Source_port)

	for {
		// 由于在后面使用SendPkg和RecvPkg有包大小限制，似乎是65535，我们这里限制一下每次接受的缓冲区大小为40K，后面可以适当调整
		data, err := client_conn.Recv(-1)
		if err != nil {
			glog.Debugf("[MappingObject::UdpConnectHandler] Port %d recv data from %s failed.", m.Info.Mapping_port, client_address)

			// 出现错误，通知B端，客户端已断开
			mashal_data, err := m.TunnelProtocolObject.BuildDataV1_0(nil, int32(mapping_id), client_address,
				dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Exception,
				tunnel_protocol.SubType_Request, m.Enable_crypto)
			if err != nil {
				glog.Debugf("[MappingObject::UdpConnectHandler] Build exception Marshaled TunnelProtocolObject failed.")
				glog.Error(err)
				break
			}

			// 向B端发送协议数据
			err = m.Tunnel_client_interface.PostDataToTunnel(mashal_data)
			if err != nil {
				glog.Debugf("[MappingObject::UdpConnectHandler] Post exception data to tunnel-B failed")
				glog.Error(err)
			}

			break
		} else {
			data_len := len(data)
			group_len := 1024 * 1024 * 4
			send_count := data_len / group_len
			remain_data_len := data_len % group_len
			glog.Debugf("[MappingObject::UdpConnectHandler] Port %d recv data from %s. length : %d", m.Info.Mapping_port, client_address, data_len)

			for index := 0; index < send_count; index++ {
				start_pos := index * group_len
				end_pos := (index + 1) * group_len
				send_data := data[start_pos:end_pos]

				tunnel_sec_data, err := m.TunnelProtocolObject.BuildDataV1_0(send_data, int32(mapping_id), client_address,
					dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
					tunnel_protocol.SubType_Response, m.Enable_crypto)
				if err != nil {
					glog.Debugf("[MappingObject::UdpConnectHandler] Build business Marshaled TunnelProtocolObject failed.")
					glog.Error(err)
					continue
				}

				// 向B端发送协议数据
				err = m.Tunnel_client_interface.PostDataToTunnel(tunnel_sec_data)
				if err != nil {
					glog.Debugf("[MappingObject::UdpConnectHandler] Post business data to tunnel-B failed")
					glog.Error(err)
				}
			}

			if remain_data_len > 0 {
				start_pos := data_len - remain_data_len
				end_pos := data_len
				send_data := data[start_pos:end_pos]

				// 组装数据
				tunnel_sec_data, err := m.TunnelProtocolObject.BuildDataV1_0(send_data, int32(mapping_id), client_address,
					dst_server_address, tunnel_protocol.NetType_TCP, tunnel_protocol.MainType_Business,
					tunnel_protocol.SubType_Response, m.Enable_crypto)
				if err != nil {
					glog.Debugf("[MappingObject::UdpConnectHandler] Build last business Marshaled TunnelProtocolObject failed.")
					glog.Error(err)
					continue
				}

				// 通过接口扔给隧道
				err = m.Tunnel_client_interface.PostDataToTunnel(tunnel_sec_data)
				if err != nil {
					glog.Debugf("[MappingObject::UdpConnectHandler] Post last business data to tunnel-B failed")
					glog.Error(err)
					continue
				}
			}
		}
	}

	// 将客户端从客户端集合中移除
	delete(m.Udp_client_table, client_address)
}
