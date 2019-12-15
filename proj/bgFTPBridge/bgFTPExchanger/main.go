package main

import (
	"bgFTPExchanger/app/service"
	_ "bgFTPExchanger/boot"
	_ "bgFTPExchanger/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	println("FTP摆渡工具")
	_ = service.Ftp_exchange_service.Initialize()
	g.Server().Run()
}
