package main

import (
	_ "bgFTPExchangeSrv/boot"
	_ "bgFTPExchangeSrv/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
