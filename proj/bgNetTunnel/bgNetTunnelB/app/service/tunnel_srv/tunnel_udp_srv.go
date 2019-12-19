package tunnel_srv

import (
	"bgNetTunnelB/app/service/cli_mgr"
	"bgNetTunnelB/tunnel_protocol"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/xtaci/kcp-go"
)

type TunnelUdpServer struct {
	kcp_listener *kcp.Listener
	//udp_server *gudp.Server
	tunnel_A_conn *kcp.UDPSession
	is_running int
	tunnel_proto *tunnel_protocol.TunnelProto
	enable_tunnel_crypto bool
	client_mgr *cli_mgr.ClientMgr
}

func (t *TunnelUdpServer) Initialize(tp *tunnel_protocol.TunnelProto, climgr *cli_mgr.ClientMgr, key []byte) {
	var err error

	t.enable_tunnel_crypto = g.Config().GetBool("tunnel.enable_crypto")
	t.client_mgr = climgr
	_ = t.client_mgr.Initialize(tp, t, t.enable_tunnel_crypto)

	// 从配置文件读取出监听端口
	srv_address := g.Config().GetString("tunnel.udp_address")
	//t.udp_server = g.UDPServer(srv_address)
	//t.udp_server.SetAddress(srv_address)
	//t.udp_server.SetHandler(t.ConnHandler)
	block, _ := kcp.NewAESBlockCrypt(key)
	t.kcp_listener, err = kcp.ListenWithOptions(srv_address, block, 10, 3)
	if err != nil {
		glog.Error(err)
		return
	}

	t.tunnel_proto = tp
	t.tunnel_A_conn = nil

	// 由于服务端在主线程启动会阻塞主线程，我们这里用一个协程来负责启动隧道服务端
	go t.Start()
}

func (t *TunnelUdpServer) Start() {
	//var err error
	for {
		kcp_cli_session, err := t.kcp_listener.AcceptKCP()
		if err != nil {
			glog.Error(err)
			continue
		}

		if t.tunnel_A_conn != nil {
			kcp_cli_session.Close()
			glog.Error(errors.New("Tunnel A conn has already exist."))
			continue
		}

		// 启动接收线程
		go t.ConnHandler(kcp_cli_session)
	}
}

func (t *TunnelUdpServer) ConnHandler(client_conn *kcp.UDPSession) {
	// 隧道客户端接入了
	glog.Debugf("TunnelTcpServer::ConnHandler() %s connect to tunnel server.", client_conn.RemoteAddr().String())
	if t.tunnel_A_conn == nil {
		t.tunnel_A_conn = client_conn
	} else {
		glog.Error("Already has a tunnel a connection...")
		_ = client_conn.Close()
		return
	}

	for {
		data, err := t.RecvPkg()
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

func (t *TunnelUdpServer) PostDataToTunnel(data []byte) error {
	var err error

	// 直接找到隧道连接对象，发送出去
	err = t.SendPkg(data)
	if err != nil {
		//glog.Debug("Send tunnel data failed :\n", data)
		glog.Error(err)
	} else {
		//glog.Debug("Send tunnel data :\n", data)
	}

	return err
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//isSymbol表示有无符号
func BytesToInt(b []byte, isSymbol bool)  (int, error){
	if isSymbol {
		return bytesToIntS(b)
	}
	return bytesToIntU(b)
}

//字节数(大端)组转成int(无符号的)
func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0},b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0,fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

//字节数(大端)组转成int(有符号)
func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0},b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0,fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}

//整形转换成字节
func IntToBytes(n int,b byte) ([]byte,error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 3,4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	}
	return nil,fmt.Errorf("IntToBytes b param is invaild")
}
///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (t *TunnelUdpServer) SendPkg(data []byte) error {
	// 这里我们采用策略来传输，首先传输一个魔术头，代表一次数据包传输开始
	magic_data := []byte("BGMagic")
	send_bytes, err := t.tunnel_A_conn.Write(magic_data)
	if err != nil {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send magic data failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send magic data succeed. Send length : %d", send_bytes)
	}

	// 第二次传输数据长度
	data_len := len(data)
	data_len_bytes, err := IntToBytes(data_len, 4)
	send_bytes, err = t.tunnel_A_conn.Write(data_len_bytes)
	if err != nil {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send data length failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send data length succeed. Send length : %d, Data length : %d",
			send_bytes, data_len)
	}

	// 第三次传输完整数据
	send_bytes, err = t.tunnel_A_conn.Write(data)
	if err != nil {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send data failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpServer::SendPkg] KCP send data succeed.")
	}

	return err
}

func (t *TunnelUdpServer) RecvPkg() ([]byte, error) {
	// 先构建缓冲区，4K字节，目前看来KCP包单包大小也不会超过2048字节
	recv_buffer := make([]byte, 4096)

	// 首先接收魔术头
	recv_bytes, err := t.tunnel_A_conn.Read(recv_buffer)
	if err != nil {
		glog.Debugf("[TunnelUdpServer::RecvPkg] KCP recv magic data failed.")
		glog.Error(err)
		return nil, err
	} else {
		glog.Debugf("[TunnelUdpServer::RecvPkg] KCP recv magic data succeed. Recv length : %d", recv_bytes)
	}

	// 检查魔术头
	magic_data := []byte("BGMagic")
	if !bytes.Equal(recv_buffer[:recv_bytes], magic_data) {
		glog.Debugf("[TunnelUdpServer::RecvPkg]Magic data not equal.")
		return nil, errors.New("Magic data not equal.")
	}

	// 接收传输数据长度
	var data_len int = 0
	recv_bytes, err = t.tunnel_A_conn.Read(recv_buffer)
	if err != nil {
		glog.Debugf("[TunnelUdpServer::RecvPkg] KCP recv data length failed.")
		glog.Error(err)
		return nil, err
	} else {
		data_len, err = bytesToIntS(recv_buffer[:recv_bytes])
		glog.Debugf("[TunnelUdpServer::RecvPkg] KCP recv data length succeed. Recv data length : %d", data_len)
	}

	// 循环接收完整数据
	var result []byte
	result = make([]byte, 0)
	var total_recv_data_len int = 0
	for {
		recv_bytes, err = t.tunnel_A_conn.Read(recv_buffer)
		if err != nil {
			glog.Debugf("[TunnelUdpServer::RecvPkg] KCP recv data length failed.")
			glog.Error(err)
			return nil, err
		} else {
			result = append(result, recv_buffer[:recv_bytes]...)
			// 计算已经接收的数据量，如果已经达到预期目标，则返回
			total_recv_data_len = total_recv_data_len + recv_bytes

			if total_recv_data_len == data_len {
				break
			}
		}
	}

	return result, err
}