/*
本对象为抽象的隧道服务端对象
本对象初始化时，根据配置文件config.toml中的隧道类型来创建对应协议的隧道服务端对象
本对象初始化时，需要传入一个观察者接口，用于向对应的观察者反馈接收到的信息
 */
package tunnel_srv

import (
	"bgNetTunnelB/app/service/cli_mgr"
	"bgNetTunnelB/tunnel_protocol"
	"crypto/md5"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

type TunnelServer struct {
	TunnelProtocolObject tunnel_protocol.TunnelProto
	Client_mgr cli_mgr.ClientMgr
}

var Tunnel_server TunnelServer

func (t *TunnelServer) Initialize() error {
	var err error

	//////////////////////////////////////////////////////////////////////////////////
	// 然后初始化隧道协议对象
	algorithm := g.Config().GetString("tunnel.crypto")
	factor := g.Config().GetString("tunnel.crypto_key_factor")

	// 生成Key，生成算法：md5(key + factor + 加密算法 + key)，结果截取前后8字节，得到16字节数据转换为byte
	md5 := md5.New()
	md5.Write([]byte("key" + factor + algorithm + "key"))
	key := md5.Sum(nil)

	// 生成IV，生成算法：md5(IV + factor + 加密算法 + IV)，结果截取前后8字节，得到16字节数据转换为byte
	md5.Reset()
	md5.Write([]byte("IV" + factor + algorithm + "IV"))
	iv := md5.Sum(nil)

	err = t.TunnelProtocolObject.Initialize(algorithm, key, iv)
	if err != nil {
		return err
	}

	///////////////////////////////////////////////////////////////////////////////
	// 根据配置文件生成对应的通道对象

	tunnel_type := g.Config().GetString("tunnel.type")
	if tunnel_type == "TCP" {
		tunnel_tcp_server := new(TunnelTcpServer)
		tunnel_tcp_server.Initialize(&t.TunnelProtocolObject, &t.Client_mgr)
		glog.Debug("TunnelServer::Initialize() Create a TCP tunnel server...")

	} else if tunnel_type == "UDP" {

	} else if tunnel_type == "LOCAL" {
		tunnel_local_server := new(TunnelLocalServer)
		tunnel_local_server.Initialize(&t.TunnelProtocolObject, &t.Client_mgr)
		glog.Debug("TunnelServer::Initialize() Create a Local tunnel server...")

	} else if tunnel_type == "FTP" {
		//
	} else if tunnel_type == "SIP" {
		//
	} else {
		//
	}

	return err
}
