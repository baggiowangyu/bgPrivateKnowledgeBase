/*

网络隧道：

主要实现以下职责：

1、启动后根据配置连接到指定的隧道服务端
2、提供接口，向隧道服务端发送数据
3、提供接口，接收隧道服务端发来的数据

 */

package bgNetTunnel

import (
	"github.com/gogf/gf/os/glog"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)


// FTP通道服务端
// 用于向目标路径写入一个请求文件，数据通过FTP摆渡机摆渡出去
// 同时等待请求结果摆渡回来
type FTPTunnelClient struct {
	SendDirectory	string
	RecvDirectory	string
	Callback		NetTunnelServerRecvCallback
}

func (f *FTPTunnelClient) Initialize(send_dir string, recv_dir string) {
	f.SendDirectory = send_dir
	f.RecvDirectory = recv_dir

	glog.Debugf("FTPTunnelClient::Initialize Set send directory (%s)", f.SendDirectory)
	glog.Debugf("FTPTunnelClient::Initialize Set recv directory (%s)", f.RecvDirectory)
}

// 发送数据
func (f *FTPTunnelClient) SendData(data []byte) error {
	// 获得当前时间的纳秒计数
	total_file_name := f.SendDirectory + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".bgdat"
	err := ioutil.WriteFile(total_file_name, data, os.ModeExclusive)
	if err != nil {
		glog.Debugf("FTPTunnelClient::SendData ioutil.WriteFile>>>%s failed.", total_file_name)
		glog.Error(err)
	} else {
		glog.Debugf("FTPTunnelClient::SendData ioutil.WriteFile>>>%s succeed.", total_file_name)
	}

	return err
}

// 接收数据
func (f *FTPTunnelClient) RecvData(callback NetTunnelServerRecvCallback) {
	// 保存回调函数
	f.Callback = callback

	// 启动协程处理
	go f.RecvThread()
}

func (f *FTPTunnelClient) RecvThread() {

	glog.Debug("FTPTunnelClient::RecvThread Started.")
	glog.Debugf("FTPTunnelClient::RecvThread Enum recv directory (%s) begin.", f.RecvDirectory)

	for index := 0; ; index++  {

		// 遍历指定目录的文件
		resluts, err := ioutil.ReadDir(f.RecvDirectory)
		if err != nil {
			glog.Debugf("FTPTunnelServer::RecvThread ioutil.ReadDir(%s) failed.", f.RecvDirectory)
			glog.Error(err)
			time.Sleep(1 * time.Microsecond)
			continue
		} else {
			//glog.Debugf("FTPTunnelServer::RecvThread ioutil.ReadDir(%s) succeed.", f.RecvDirectory)
		}

		for _, file_info := range resluts {

			file_path := f.RecvDirectory + "/" + file_info.Name()

			if file_info.IsDir() {
				// 是目录，直接跳过
				glog.Debugf("FTPTunnelServer::RecvThread path(%s) is directory, skip...", file_path)
				continue
			} else {
				// 是文件
				glog.Debugf("FTPTunnelServer::RecvThread path(%s) is file, start handle...", file_path)
				data, err := ioutil.ReadFile(file_path)
				if err != nil {
					glog.Error("FTPTunnelClient::RecvThread ")
					time.Sleep(1 * time.Microsecond)
					continue
				}

				err = f.Callback(data)
				if err != nil {
					glog.Error("FTPTunnelClient::RecvThread Callback failed. " + err.Error())
					//time.Sleep(1 * time.Microsecond)
					//continue
				}

				// 删除这个文件
				err = os.Remove(file_path)
				if err != nil {
					glog.Error("FTPTunnelClient::RecvThread os.Remove(" + file_path + ") failed. " + err.Error())
					time.Sleep(1 * time.Microsecond)
					continue
				}
			}
		}

		time.Sleep(1 * time.Microsecond)
	}
}

