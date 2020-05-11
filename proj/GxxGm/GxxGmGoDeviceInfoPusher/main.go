package main

import (
	_ "GxxGmGoDeviceInfoPusher/boot"
	_ "GxxGmGoDeviceInfoPusher/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	//glog.Debug("test...")
	//glog.Info("test...")
	//glog.Notice("test...")
	//glog.Warning("test...")
	//glog.Error("test...")
	//glog.Critical("test...")
	//glog.Panic("test...")
	//glog.Fatal("test...")

	g.Server().Run()
}
