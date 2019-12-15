package router

import (
	"bgFTPExchanger/app/api/info"
	"github.com/gogf/gf/frame/g"
)

func init() {
	g.Server().BindObject("GET:/info", new(info.Controller), "info")
}