package bgNetTunnel

import (
	"bgNetTunnelServer/library/bgNetProtocol"
	"bgNetTunnelServer/library/bgNetProtocol/bgNetMessage"
	"errors"
	"github.com/gogf/gf/os/glog"
	"github.com/golang/protobuf/proto"
)

type NetTunnelClient struct {
	// 网络通道接口（抽象出来的）
	net_turnnel_interface 	NetTunnelServerInterface

	net_protocol_object		bgNetProtocol.NetProtocol
}

func (n *NetTunnelClient) Initialize(tunnel_srv_ip string, tunnel_srv_port int, user string, pass string, send_dir string, recv_dir string, tunnel_type string, callback NetTunnelServerRecvCallback) error {
	var err error

	err = n.net_protocol_object.Initialize("AES", []byte("1234567890123456"), []byte("1234567890123456"))
	if err != nil {
		glog.Debug("NetTunnelClient::Initialize Initialize net protocol object failed.")
		glog.Error(err)
		panic(err)
	} else {
		glog.Debug("NetTunnelClient::Initialize Initialize net protocol object succeed.")
	}

	if tunnel_type == "FTP" {

		ftp_tunnel_client := new(FTPTunnelClient)
		n.net_turnnel_interface = ftp_tunnel_client
		ftp_tunnel_client.Initialize(send_dir, recv_dir)
		ftp_tunnel_client.RecvData(callback)
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

func (n *NetTunnelClient) SendMessage(msg bgNetMessage.NetMessageV1) error {
	result, err := proto.Marshal(&msg)
	if err != nil {
		glog.Debug("NetTunnelClient::SendMessage marshal message failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("NetTunnelClient::SendMessage marshal message succeed.")
	}

	// 执行加密
	cipher_data, err := n.net_protocol_object.EncryptData(result)
	if err != nil {
		glog.Debug("NetTunnelClient::SendMessage encrypt marshaled data failed.")
		return err
	} else {
		glog.Debug("NetTunnelClient::SendMessage encrypt marshaled data succeed.")
	}

	// 发送数据
	err = n.net_turnnel_interface.SendData(cipher_data)
	return err
}
