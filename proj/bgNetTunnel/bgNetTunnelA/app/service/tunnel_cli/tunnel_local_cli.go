package tunnel_cli

import (
	"github.com/gogf/gf/os/glog"
	tsgutils "github.com/typa01/go-utils"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type TunnelLocalClient struct {
	recv_dir string
	send_dir string
	is_running int
	tunnel_client_obsever_interface TunnelClientObserverInterface
}

func (t *TunnelLocalClient) Initialize(rdir string, sdir string, inter TunnelClientObserverInterface) error {
	t.recv_dir = rdir
	t.send_dir = sdir
	t.tunnel_client_obsever_interface = inter

	return t.ConnectToTunnelB()
}

func (t *TunnelLocalClient) ConnectToTunnelB() error {
	var err error

	// 连接成功，启动接收线程开始接收
	go t.RecvThread()

	return err
}

func (t *TunnelLocalClient) PostDataToTunnel(data []byte) error {
	// 根据当前时间，直接写入文件，等待后续文件被摆渡走
	total_file_name := t.send_dir + "/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + tsgutils.GUID() + ".bgdat"
	err := ioutil.WriteFile(total_file_name, data, os.ModeExclusive)
	if err != nil {
		glog.Error(err)
	}

	return err
}

func (t *TunnelLocalClient) RecvThread() {
	// 接收线程
	t.is_running = 1
	for {
		file_infos, err := ioutil.ReadDir(t.recv_dir)
		if err != nil {
			glog.Error(err)
			time.Sleep(time.Duration(1000) * time.Millisecond)
			continue
		}

		for _, file_info := range file_infos{
			if file_info.IsDir() {
				continue
			}

			full_path := t.recv_dir + "/" + file_info.Name()
			data, err := ioutil.ReadFile(full_path)
			if err != nil {
				glog.Error(err)
				continue
			}

			err = os.RemoveAll(full_path)
			if err != nil {
				glog.Error(err)
			}

			// 这里应该往上层观察者扔了
			err = t.tunnel_client_obsever_interface.PeekDataFromTunnel(data)
		}
	}

	t.is_running = 0
}
