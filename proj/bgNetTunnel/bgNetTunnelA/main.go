package main

import (
	"bgNetTunnelA/app/service/mapping_srv"
	_ "bgNetTunnelA/boot"
	_ "bgNetTunnelA/router"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	// 创建一个映射服务，并初始化
	err := mapping_srv.Mapping_service_instance.Initialize()
	if err != nil {
		glog.Debug("Initialize mapping service instance failed.")
		glog.Error(err)
		return
	}

	// 启动Web服务
	g.Server().Run()
}
