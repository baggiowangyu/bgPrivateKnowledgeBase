/*
本对象为抽象的隧道客户端对象
本对象初始化时，根据配置文件config.toml中的隧道类型来创建对应协议的隧道客户端对象
本对象初始化时，需要传入一个观察者接口，用于向对应的观察者反馈接收到的信息
 */
package tunnel_cli

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

type TunnelClientInterface interface {
	// 向隧道服务端发送数据
	PostDataToTunnel(data []byte) error
}

/*
TunnelClientObserverInterface
隧道客户端观察者接口
一般建议有隧道客户端持有者实现此接口
 */
type TunnelClientObserverInterface interface {
	PeekDataFromTunnel(data []byte) error
}

type TunnelClient struct {
	client_interface TunnelClientInterface
}

func (t *TunnelClient) Initialize(inter TunnelClientObserverInterface) error {
	var err error

	// 根据配置文件生成对应的通道对象
	tunnel_type := g.Config().GetString("tunnel.type")
	if tunnel_type == "TCP" {
		//
		tunnel_b_address := g.Config().GetString("tunnel.b_address")
		tunnel_tcp_client := new(TunnelTcpClient)
		err = tunnel_tcp_client.Initialize(tunnel_b_address, inter)
		if err != nil {
			glog.Error(err)
			return err
		}

		t.client_interface = tunnel_tcp_client
		glog.Debug("TunnelClient::Initialize() Create a TCP tunnel client...")

	} else if tunnel_type == "UDP" {
		//
	} else if tunnel_type == "LOCAL" {
		sdir := g.Config().GetString("tunnel.local_send_dir")
		rdir := g.Config().GetString("tunnel.local_recv_dir")
		tunnel_local_client := new(TunnelLocalClient)
		tunnel_local_client.Initialize(rdir, sdir, inter)
		if err != nil {
			glog.Error(err)
			return err
		}

		t.client_interface = tunnel_local_client
		glog.Debug("TunnelClient::Initialize() Create a LOCAL tunnel client...")

	} else if tunnel_type == "FTP" {
		//
	} else if tunnel_type == "SIP" {
		//
	} else {
		//
	}

	return err
}

func (t *TunnelClient) PostDataToTunnel(data []byte) error {
	err := t.client_interface.PostDataToTunnel(data)
	return err
}