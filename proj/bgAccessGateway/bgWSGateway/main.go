package main

import (
	_ "bgWSGateway/boot"
	_ "bgWSGateway/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
