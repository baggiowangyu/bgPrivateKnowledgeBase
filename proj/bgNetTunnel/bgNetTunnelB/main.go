package main

import (
	"bgNetTunnelB/app/service/tunnel_srv"
	_ "bgNetTunnelB/boot"
	_ "bgNetTunnelB/router"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	//tunnel_server := new(tunnel_srv.TunnelServer)
	err := tunnel_srv.Tunnel_server.Initialize()
	if err != nil {
		glog.Error(err)
		return
	}

	g.Server().Run()
}
