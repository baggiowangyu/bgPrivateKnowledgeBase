package service

import (
	"bgFTPExchanger/library/ftpexchanger"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"time"
)

var Ftp_exchange_service = FTPExchangeService{}

func init() {

}

func exchange_A2B(s *FTPExchangeService) {
	glog.Info("goruntime exchange_A2B started...")

	err := s.A2B_exchanger.Initialize(s.Ftp_A_address, s.Ftp_A_login_name, s.Ftp_A_login_pass, s.Ftp_A_send_dir, s.Ftp_B_address, s.Ftp_B_login_name, s.Ftp_B_login_pass, s.Ftp_B_recv_dir)
	if err != nil {
		glog.Error(err)
	}

	for i := 0;; i++ {
		// 等待指定的时间后，开始摆渡
		time.Sleep(time.Duration(s.A2B_interval) * time.Millisecond)

		glog.Debugf("Start Ferry from %s to %s", s.Ftp_A_address, s.Ftp_B_address)
		err = s.A2B_exchanger.Ferry()
		if err != nil {
			glog.Error(err)
		}
	}
}

func exchange_B2A(s *FTPExchangeService) {
	glog.Info("goruntime exchange_B2A started...")

	err := s.B2A_exchanger.Initialize(s.Ftp_B_address, s.Ftp_B_login_name, s.Ftp_B_login_pass, s.Ftp_B_send_dir, s.Ftp_A_address, s.Ftp_A_login_name, s.Ftp_A_login_pass, s.Ftp_A_recv_dir)
	if err != nil {
		glog.Error(err)
	}

	for i := 0;; i++ {
		// 等待指定的时间后，开始摆渡
		time.Sleep(time.Duration(s.A2B_interval) * time.Millisecond)

		glog.Debugf("Start Ferry from %s to %s", s.Ftp_B_address, s.Ftp_A_address)
		err = s.B2A_exchanger.Ferry()
		if err != nil {
			glog.Error(err)
		}
	}
}

type FTPExchangeService struct {
	Ftp_A_address 		string
	Ftp_A_port 			int
	Ftp_A_login_name 	string
	Ftp_A_login_pass 	string
	Ftp_A_send_dir 		string
	Ftp_A_recv_dir 		string

	Ftp_B_address 		string
	Ftp_B_port 			int
	Ftp_B_login_name 	string
	Ftp_B_login_pass 	string
	Ftp_B_send_dir 		string
	Ftp_B_recv_dir 		string

	A2B_interval		int
	B2A_interval		int

	A2B_exchanger		ftpexchanger.FTPExchanger
	B2A_exchanger		ftpexchanger.FTPExchanger
}

func (f *FTPExchangeService) Initialize() error {
	// 读取配置文件
	f.Ftp_A_address = g.Config().GetString("ftp-A-config.address")
	f.Ftp_A_port = g.Config().GetInt("ftp-A-config.port")
	f.Ftp_A_login_name = g.Config().GetString("ftp-A-config.login-name")
	f.Ftp_A_login_pass = g.Config().GetString("ftp-A-config.login-pass")
	f.Ftp_A_send_dir = g.Config().GetString("ftp-A-config.send-dir")
	f.Ftp_A_recv_dir = g.Config().GetString("ftp-A-config.recv-dir")

	f.Ftp_B_address = g.Config().GetString("ftp-B-config.address")
	f.Ftp_B_port = g.Config().GetInt("ftp-B-config.port")
	f.Ftp_B_login_name = g.Config().GetString("ftp-B-config.login-name")
	f.Ftp_B_login_pass = g.Config().GetString("ftp-B-config.login-pass")
	f.Ftp_B_send_dir = g.Config().GetString("ftp-B-config.send-dir")
	f.Ftp_B_recv_dir = g.Config().GetString("ftp-B-config.recv-dir")

	f.A2B_interval = g.Config().GetInt("A2B.interval")
	f.B2A_interval = g.Config().GetInt("B2A.interval")

	// 读取完毕后，启动两个同步线程
	go exchange_A2B(f)
	go exchange_B2A(f)

	return nil
}
