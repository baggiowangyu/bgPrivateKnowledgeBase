package main

import (
	"bgFTPBridgeServer/app/service/bridgemgr"
	_ "bgFTPBridgeServer/boot"
	_ "bgFTPBridgeServer/router"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	// 启动TCP服务器
	err := bridgemgr.Bridge_Object.StartUp()
	if err != nil {
		glog.Error(err.Error())
	}

	// 启动WEB服务器
	http_server := g.Server()
	http_server.Run()
}
