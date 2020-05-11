package main

import (
	_ "bgGB28181SignalGateway/boot"
	_ "bgGB28181SignalGateway/router"
	"github.com/gogf/gf/frame/g"
)

func main() {

	g.Server().Run()
}
