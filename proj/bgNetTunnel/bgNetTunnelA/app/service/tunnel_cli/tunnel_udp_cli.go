package tunnel_cli

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/xtaci/kcp-go"
	"time"
)

type TunnelUdpClient struct {
	tunnel_B_address string
	block kcp.BlockCrypt
	conn *kcp.UDPSession
	is_running int
	tunnel_client_obsever_interface TunnelClientObserverInterface

}

func (t *TunnelUdpClient) Initialize(tunnel_b_address string, inter TunnelClientObserverInterface, key []byte) error {
	t.tunnel_B_address = tunnel_b_address
	t.block, _ = kcp.NewAESBlockCrypt(key)
	t.tunnel_client_obsever_interface = inter

	return t.ConnectToTunnelB()
}

func (t *TunnelUdpClient) ConnectToTunnelB() error {
	var err error

	reconnect_interval := g.Config().GetInt("tunnel.tunnel_reconnect_interval")
	// 尝试连接到目标服务器，这里设置为永久重试，只不过有一个重试间隔时间
	for {
		t.conn, err = kcp.DialWithOptions(t.tunnel_B_address, t.block, 10, 3)
		//t.conn, err = gudp.NewConn(t.tunnel_B_address)
		if err != nil {
			glog.Warningf(err.Error())
			time.Sleep(time.Duration(reconnect_interval) * time.Millisecond)
			continue
		}

		break
	}

	// 连接成功，启动接收线程开始接收
	go t.RecvThread()

	return err
}

func (t *TunnelUdpClient) PostDataToTunnel(data []byte) error {

	// 这里是按照UDP包发送的，我们可能要引入
	err := t.SendPkg(data)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::PostDataToTunnel] SendData failed.")
		glog.Error(err)
	} else {
		glog.Debugf("[TunnelUdpClient::PostDataToTunnel] SendData succeed:\n%s", data)
	}

	return err
}

func (t *TunnelUdpClient) RecvThread() {
	// 接收线程，接收隧道B端发来的数据
	t.is_running = 1

	for {
		// 这里采用RecvPkg

		data, err := t.RecvPkg()
		if err != nil {
			// 当这里出现异常的时候，我们应当进行无限次重连，确保通道能正常恢复
			glog.Debugf("[TunnelUdpClient::RecvThread] Recv tunnel data failed")
			glog.Error(err)
			break
		} else {
			glog.Debugf("[TunnelUdpClient::RecvThread] Recv tunnel data :\n%s", data)
		}

		// 这里应该往上层观察者扔了
		err = t.tunnel_client_obsever_interface.PeekDataFromTunnel(data)
	}

	t.is_running = 0

	// 启动重连线程
	go t.ConnectToTunnelB()
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

func (t *TunnelUdpClient) SendPkg(data []byte) error {
	// 这里我们采用策略来传输，首先传输一个魔术头，代表一次数据包传输开始
	magic_data := []byte("BGMagic")
	send_bytes, err := t.conn.Write(magic_data)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send magic data failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send magic data succeed. Send length : %d", send_bytes)
	}

	// 第二次传输数据长度
	data_len := len(data)
	data_len_bytes, err := IntToBytes(data_len, 4)
	send_bytes, err = t.conn.Write(data_len_bytes)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send data length failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send data length succeed. Send length : %d, Data length : %d",
			send_bytes, data_len)
	}

	// 第三次传输完整数据
	send_bytes, err = t.conn.Write(data)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send data failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debugf("[TunnelUdpClient::SendPkg] KCP send data succeed.")
	}

	return err
}

func (t *TunnelUdpClient) RecvPkg() ([]byte, error) {
	// 先构建缓冲区，4K字节，目前看来KCP包单包大小也不会超过2048字节
	recv_buffer := make([]byte, 4096)

	// 首先接收魔术头
	recv_bytes, err := t.conn.Read(recv_buffer)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::RecvPkg] KCP recv magic data failed.")
		glog.Error(err)
		return nil, err
	} else {
		glog.Debugf("[TunnelUdpClient::RecvPkg] KCP recv magic data succeed. Recv length : %d", recv_bytes)
	}

	// 检查魔术头
	magic_data := []byte("BGMagic")
	if !bytes.Equal(recv_buffer[:recv_bytes], magic_data) {
		glog.Debugf("[TunnelUdpClient::RecvPkg] Magic data not equal.")
		return nil, errors.New("Magic data not equal.")
	}

	// 接收传输数据长度
	var data_len int = 0
	recv_bytes, err = t.conn.Read(recv_buffer)
	if err != nil {
		glog.Debugf("[TunnelUdpClient::RecvPkg] KCP recv data length failed.")
		glog.Error(err)
		return nil, err
	} else {
		data_len, err = bytesToIntS(recv_buffer[:recv_bytes])
		glog.Debugf("[TunnelUdpClient::RecvPkg] KCP recv data length succeed. Recv data length : %d", data_len)
	}

	// 循环接收完整数据
	var result []byte
	result = make([]byte, 0)
	var total_recv_data_len int = 0
	for {
		recv_bytes, err = t.conn.Read(recv_buffer)
		if err != nil {
			glog.Debugf("[TunnelUdpClient::RecvPkg] KCP recv data length failed.")
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
