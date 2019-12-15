/*
本地隧道服务端
本隧道服务端是一个本地服务端，持有一个隧道服务端观察者接口，用于向其反馈接收到的数据
本隧道服务端实现了数据发送接口，可以抽象为隧道服务端接口，供隧道服务端调用
本隧道不同于TCP、UDP、SIP、FTP等直接网络通信的应用场景
本隧道应用场景为：
1、网络通信隧道为 >>>FTP摆渡<<< ，即通信数据先以文本形式写入本地FTP服务中的发送目录，第三方摆渡程序将文件摆渡到目标服务器FTP中的接收目录
2、该摆渡通信应支持一路发送、一路接收。共两路FTP摆渡行为
*/
package tunnel_srv

import (
	"bgNetTunnelB/app/service/cli_mgr"
	"bgNetTunnelB/tunnel_protocol"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	tsgutils "github.com/typa01/go-utils"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type TunnelLocalServer struct {
	// 用于与TunnelA通信的两个目录
	recv_dir string
	send_dir string

	is_running int
	tunnel_proto *tunnel_protocol.TunnelProto
	enable_tunnel_crypto bool
	client_mgr *cli_mgr.ClientMgr

}

func (t *TunnelLocalServer) Initialize(tp *tunnel_protocol.TunnelProto, climgr *cli_mgr.ClientMgr) {

	t.enable_tunnel_crypto = g.Config().GetBool("tunnel.enable_crypto")
	t.client_mgr = climgr
	_ = t.client_mgr.Initialize(tp, t, t.enable_tunnel_crypto)

	// 从配置文件读取出接收目录和发送目录
	t.recv_dir = g.Config().GetString("tunnel.local_recv_dir")
	t.send_dir = g.Config().GetString("tunnel.local_send_dir")
	t.tunnel_proto = tp

	// 由于服务端在主线程启动会阻塞主线程，我们这里用一个协程来负责启动隧道服务端
	go t.RecvHandler()
}

func (t *TunnelLocalServer) RecvHandler() {
	// 这里用于扫描接收目录，读取对应的文件，并解析
	for {
		file_infos, err := ioutil.ReadDir(t.recv_dir)
		if err != nil {
			glog.Error(err)
			time.Sleep(time.Duration(1000) * time.Millisecond)
		}

		for _, file_info := range file_infos {
			if file_info.IsDir() {
				continue
			}

			full_path := t.recv_dir + "/" + file_info.Name()
			data, err := ioutil.ReadFile(full_path)
			if err != nil {
				glog.Error(err)
				continue
			}

			// 接收到数据，首先反序列化，得到整包，判断是业务数据还是异常数据，是业务数据则扔给客户端管理器
			tunnel_sec_protocol_data, err := t.tunnel_proto.Unmarshal(data)
			if err != nil {
				glog.Error(err)
				continue
			}

			if tunnel_sec_protocol_data.Main == tunnel_protocol.MainType_Business {
				// 向客户端管理器发送业务数据
				go t.client_mgr.AsynchronousRecvDataFromTunnel(tunnel_sec_protocol_data.Data)
			} else if tunnel_sec_protocol_data.Main == tunnel_protocol.MainType_Exception {
				// 处理异常，可能是A端客户端主动断开了，我们在这里一并通知客户端管理器，并发执行
				go t.client_mgr.RecvExceptionFromTunnel(tunnel_sec_protocol_data.Data)
			}
		}
	}
}

func (t *TunnelLocalServer) PostDataToTunnel(data []byte) error {
	// 根据当前时间，直接写入文件，等待后续文件被摆渡走
	total_file_name := t.send_dir + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + tsgutils.GUID() + ".bgdat"
	err := ioutil.WriteFile(total_file_name, data, os.ModeExclusive)
	if err != nil {
		glog.Error(err)
	}

	return err
}