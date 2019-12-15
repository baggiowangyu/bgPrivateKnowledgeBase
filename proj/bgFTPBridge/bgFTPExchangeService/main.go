package main

import (
	_ "bgFTPExchangeService/boot"
	_ "bgFTPExchangeService/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
