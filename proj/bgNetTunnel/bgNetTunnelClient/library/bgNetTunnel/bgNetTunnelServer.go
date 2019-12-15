package bgNetTunnel

import (
	"bgNetTunnelClient/library/bgNetProtocol"
	"bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"
	"errors"
	"github.com/gogf/gf/os/glog"
	"github.com/golang/protobuf/proto"
)

type NetTunnelServer struct {
	// 网络通道接口
	net_turnnel_interface 	NetTunnelServerInterface

	// 隧道协议对象
	net_protocol_object		bgNetProtocol.NetProtocol

	// 映射管理器的回调对象
	net_tunnel_callback		NetTunnelServerRecvCallback
}

func (n *NetTunnelServer) Initialize(tunnel_srv_ip string, tunnel_srv_port int, user string, pass string, send_dir string, recv_dir string, tunnel_type string, callback NetTunnelServerRecvCallback) error {

	var err error
	err = n.net_protocol_object.Initialize("AES", []byte("1234567890123456"), []byte("1234567890123456"))
	if err != nil {
		glog.Debug("NetTunnelClient::net_protocol_object::Initialize failed.")
		glog.Error(err)

		// 考虑一下，这里用panic是否合适
		panic(err)
	} else {
		glog.Debug("NetTunnelClient::net_protocol_object::Initialize succeed.")
	}

	if tunnel_type == "FTP" {

		ftp_tunnel_client := new(FTPTunnelClient)
		n.net_turnnel_interface = ftp_tunnel_client
		n.net_tunnel_callback = callback
		ftp_tunnel_client.Initialize(send_dir, recv_dir)
		ftp_tunnel_client.RecvData(n.TunnelRecvCallback)
		err = nil

	} else if tunnel_type == "TCP" {

		glog.Fatal("NetTunnelClient::Initialize TCP tunnel is not implemented.")
		err = errors.New("TCP tunnel is not implemented.")

	} else if tunnel_type == "UDP" {

		glog.Fatal("NetTunnelClient::Initialize UDP tunnel is not implemented.")
		err = errors.New("UDP tunnel is not implemented.")

	}else if tunnel_type == "SIP" {

		glog.Fatal("NetTunnelClient::Initialize SIP tunnel is not implemented.")
		err = errors.New("SIP tunnel is not implemented.")

	}

	return err
}

func (n *NetTunnelServer) SendMessage(msg bgNetMessage.NetMessageV1) error {
	//var b []byte
	//result, err := msg.XXX_Marshal(b, false)
	result, err := proto.Marshal(&msg)
	if err != nil {
		glog.Debug("NetTunnelClient::SendMessage Marshal NetMessageV1 failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("NetTunnelClient::SendMessage Marshal NetMessageV1 succeed.")
	}

	// 执行加密
	cipher_data, err := n.net_protocol_object.EncryptData(result)
	if err != nil {
		glog.Debug("NetTunnelClient::SendMessage service.NetPeorocolObject.EncryptData failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("NetTunnelClient::SendMessage service.NetPeorocolObject.EncryptData succeed.")
	}

	// 发送数据
	err = n.net_turnnel_interface.SendData(cipher_data)
	if err != nil {
		glog.Debug("NetTunnelClient::SendMessage n.net_turnnel_interface.SendData failed.")
		glog.Error(err)
	} else {
		glog.Debug("NetTunnelClient::SendMessage n.net_turnnel_interface.SendData succeed.")
	}

	return err
}

func (n *NetTunnelServer) TunnelRecvCallback(data []byte) error {

	glog.Debug("NetTunnelServer::TunnelRecvCallback recv data.")

	// 首先协议栈解密，反序列化
	plain_data, err := n.net_protocol_object.DecryptData(data)
	if err != nil {
		glog.Debug("NetTunnelServer::TunnelRecvCallback NetProtocol decrypt data failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("NetTunnelServer::TunnelRecvCallback NetProtocol decrypt data succeed.")
	}

	err = n.net_tunnel_callback(plain_data)
	if err != nil {
		glog.Error(err)
		return err
	}

	return nil
}