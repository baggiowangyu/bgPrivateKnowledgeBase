package ftpexchanger

import (
	"bytes"
	"github.com/gogf/gf/os/glog"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	//"sync"
	"time"
)

type FTPExchanger struct {
	Source_ftp_client	*ftp.ServerConn
	Source_address		string
	Source_username		string
	Source_password		string
	Source_directory	string

	Target_ftp_client 	*ftp.ServerConn
	Target_address		string
	Target_username		string
	Target_password 	string
	Target_directory	string

	// 等待数组，用于同步
	//wait_group			sync.WaitGroup

}

// 初始化FTP摆渡器，主要是连接到两个FTP服务器，并转到相应的目录下
func (f *FTPExchanger) Initialize(src_addr string, src_user string, src_pass string, src_cwd string, tar_addr string, tar_user string, tar_pass string, tar_cwd string) error {
	var err error

	err = f.connect_and_login_source(src_addr, src_user, src_pass, src_cwd)
	if err != nil {
		return err
	}

	err = f.connect_and_login_target(tar_addr, tar_user, tar_pass, tar_cwd)
	if err != nil {
		return err
	}

	return nil
}

func (f *FTPExchanger) Ferry() error {
	files, err := f.Source_ftp_client.List("/" + f.Source_directory)
	if err != nil {
		glog.Error("Refresh source ftp server (ftp://" + f.Target_address + "/" + f.Source_directory +") failed. \n" + err.Error())

		// 这里先断开连接，再重新登录，然后等待下一次扫描
		err = f.ReconnectSourceFTPServer()
		if err != nil {
			glog.Error("FTPExchanger::Ferry ReconnectSourceFTPServer failed.")
		} else {
			glog.Info("FTPExchanger::Ferry ReconnectSourceFTPServer succeed.")
		}

		return err
	} else {
		glog.Debug("Refresh source ftp server (ftp://" + f.Target_address + "/" + f.Source_directory +") succeed.")
	}

	for _, entry_element := range files {

		if entry_element.Type == ftp.EntryTypeFile {

			glog.Info("Find file (ftp://" + f.Target_address + "/" + f.Source_directory + "/" + entry_element.Name + ")")

			// 读取文件数据
			response, err := f.Source_ftp_client.Retr("/" + f.Source_directory + "/" + entry_element.Name)
			if err != nil {
				glog.Error("Retr (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") failed. " + err.Error())
				continue
			} else {
				glog.Debug("Retr (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") succeed.")
			}

			// 转换为字节切片
			buf, err := ioutil.ReadAll(response)
			if err != nil {
				glog.Error("Read source file (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") failed. " + err.Error())
				_ = response.Close()
				continue
			} else {
				glog.Info("Read source file (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") succeed.")
			}

			// 构建新的缓冲区
			data := bytes.NewBuffer(buf)

			// 写入目标服务器，如果出错，尝试3次
			retry_count := 0
			for {
				err = f.Target_ftp_client.Stor("/" + f.Target_directory + "/" + entry_element.Name, data)
				_ = response.Close()
				if err != nil {
					glog.Error("Stor (ftp://" + f.Target_address + "/" + f.Target_directory + "/" + entry_element.Name + ") failed. " + err.Error())
					if retry_count < 3 {
						// 尝试重连
						err = f.ReconnectTargetFTPServer()
						if err != nil {
							// 重连失败
							glog.Debug("FTPExchanger::Ferry ReconnectTargetFTPServer falied.")
							glog.Error(err)
						} else {
							glog.Info("FTPExchanger::Ferry ReconnectTargetFTPServer succeed.")
						}

						retry_count++
						continue
					} else {
						// 重试三次不行，尝试重连客户端
						break
					}

				} else {
					glog.Info("Stor (ftp://" + f.Target_address + "/" + f.Target_directory + "/" + entry_element.Name + ") succeed.")
					break
				}
			}

			if err != nil {
				continue
			}

			// 成功后，删除源服务器上的文件
			err = f.Source_ftp_client.Delete("/" + f.Source_directory + "/" + entry_element.Name)
			if err != nil {
				glog.Error("Delete (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") failed. " + err.Error())
				continue
			} else {
				glog.Info("Delete (ftp://" + f.Source_address + "/" + f.Source_directory + "/" + entry_element.Name + ") succeed.")
			}
		}
	}

	return nil
}

// 私有函数，连接到服务器，登录，转到指定目录
func (f *FTPExchanger) connect_and_login_source(src_addr string, src_user string, src_pass string, src_cwd string) error {
	var err error

	// 连接到源FTP服务器
	f.Source_ftp_client, err = ftp.Dial(src_addr, ftp.DialWithTimeout(30 * time.Second))
	if err != nil {
		glog.Error("Connect to source ftp server (ftp://" + src_addr + ") failed. " + err.Error())
		return err
	} else {
		glog.Debug("Connect to source ftp server (ftp://" + src_addr + ") succeed. ")
	}

	// 登录源FTP服务器
	err = f.Source_ftp_client.Login(src_user, src_pass)
	if err != nil {
		glog.Error("Login to souorce ftp server (ftp://" + src_addr + "@" + src_user + ":" + src_pass + ")" + err.Error())
		return err
	} else {
		glog.Debug("Login to souorce ftp server (ftp://" + src_addr + "@" + src_user + ":" + src_pass + ") succeed.")
	}

	//// 将源FTP服务器工作目录转移到指定目录
	//err = f.Source_ftp_client.ChangeDir(src_cwd)
	//if err != nil {
	//	glog.Error("Set source ftp server's current working directory (ftp://" + src_addr + "/" + src_cwd + ") failed" + err.Error())
	//	return err
	//} else {
	//	glog.Debug("Set source ftp server's current working directory (ftp://" + src_addr + "/" + src_cwd + ") succeed.")
	//}

	f.Source_address = src_addr
	f.Source_username = src_user
	f.Source_password = src_pass
	f.Source_directory = src_cwd

	return nil
}

func (f *FTPExchanger) connect_and_login_target(tar_addr string, tar_user string, tar_pass string, tar_cwd string) error {
	var err error

	// 连接到目标FTP服务器
	f.Target_ftp_client, err = ftp.Dial(tar_addr, ftp.DialWithTimeout(30 * time.Second))
	if err != nil {
		glog.Error("Connect to target ftp server (ftp://" + tar_addr + ") failed. " + err.Error())
		return err
	} else {
		glog.Debug("Connect to target ftp server (ftp://" + tar_addr + ") succeed. ")
	}

	// 登录目标FTP服务器
	err = f.Target_ftp_client.Login(tar_user, tar_pass)
	if err != nil {
		glog.Error("Login to target ftp server (ftp://" + tar_addr + "@" + tar_user + ":" + tar_pass + ")" + err.Error())
		return err
	} else {
		glog.Debug("Login to target ftp server (ftp://" + tar_addr + "@" + tar_user + ":" + tar_pass + ") succeed.")
	}

	//// 将目标FTP服务器工作目录转移到指定目录
	//err = f.Target_ftp_client.ChangeDir(tar_cwd)
	//if err != nil {
	//	glog.Error("Set target ftp server's current working directory (ftp://" + tar_addr + "/" + tar_cwd + ") failed" + err.Error())
	//	return err
	//} else {
	//	glog.Debug("Set target ftp server's current working directory (ftp://" + tar_addr + "/" + tar_cwd + ") succeed.")
	//}

	f.Target_address = tar_addr
	f.Target_username = tar_user
	f.Target_password = tar_pass
	f.Target_directory = tar_cwd

	return nil
}

func (f *FTPExchanger) ReconnectSourceFTPServer() error {
	// 首先断开连接
	_ = f.Source_ftp_client.Logout()

	// 尝试重新连接
	err := f.connect_and_login_source(f.Source_address, f.Source_username, f.Source_password, f.Source_directory)
	return err
}

func (f *FTPExchanger) ReconnectTargetFTPServer() error {
	// 首先断开连接
	_ = f.Target_ftp_client.Logout()

	// 尝试重新连接
	err := f.connect_and_login_target(f.Target_address, f.Target_username, f.Target_password, f.Target_directory)
	return err
}