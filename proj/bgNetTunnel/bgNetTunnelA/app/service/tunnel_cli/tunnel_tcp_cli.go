/*
TCP隧道客户端
本隧道客户端是一个TCP客户端，持有一个隧道客户端观察者接口，用于向其反馈接收到的数据
本隧道客户端实现了数据发送接口，可以抽象为隧道客户端接口，供隧道客户端调用
 */
package tunnel_cli

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/os/glog"
	"time"
)

type TunnelTcpClient struct {
	tunnel_B_address string
	conn *gtcp.Conn
	is_running int
	tunnel_client_obsever_interface TunnelClientObserverInterface

}

func (t *TunnelTcpClient) Initialize(tunnel_b_address string, inter TunnelClientObserverInterface) error {
	t.tunnel_B_address = tunnel_b_address
	t.tunnel_client_obsever_interface = inter

	return t.ConnectToTunnelB()
}

func (t *TunnelTcpClient) ConnectToTunnelB() error {
	var err error

	reconnect_interval := g.Config().GetInt("tunnel.tunnel_reconnect_interval")
	// 尝试连接到目标服务器，这里设置为永久重试，只不过有一个重试间隔时间
	for {
		t.conn, err = gtcp.NewConn(t.tunnel_B_address)
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

func (t *TunnelTcpClient) PostDataToTunnel(data []byte) error {

	opt := gtcp.PkgOption{
		HeaderSize:4,
		MaxDataSize : 16777215,
	}

	// 隧道内的数据尽量不要粘包，我们整包发
	err := t.conn.SendPkg(data, opt)
	if err != nil {
		glog.Debugf("SendData failed.")
		glog.Error(err)
	} else {
		glog.Debugf("SendData succeed:\n%s", data)
	}

	return err
}

func (t *TunnelTcpClient) RecvThread() {
	// 接收线程，接收隧道B端发来的数据
	t.is_running = 1
	opt := gtcp.PkgOption{
		HeaderSize:4,
		MaxDataSize : 16777215,
	}

	for {
		// 这里采用RecvPkg

		data, err := t.conn.RecvPkg(opt)
		if err != nil {
			// 当这里出现异常的时候，我们应当进行无限次重连，确保通道能正常恢复
			glog.Debugf("Recv tunnel data failed")
			glog.Error(err)
			break
		} else {
			glog.Debugf("Recv tunnel data :\n%s", data)
		}

		// 这里应该往上层观察者扔了
		err = t.tunnel_client_obsever_interface.PeekDataFromTunnel(data)
	}

	t.is_running = 0

	// 启动重连线程
	go t.ConnectToTunnelB()
}