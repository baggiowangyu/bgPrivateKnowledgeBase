package main

import (
	_ "bgApiGateway/boot"
	_ "bgApiGateway/router"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func main() {
	err := g.Server("Manager").Start()
	if err != nil {
		glog.Error(err)
	}

	g.Server("Gateway").Run()
}
