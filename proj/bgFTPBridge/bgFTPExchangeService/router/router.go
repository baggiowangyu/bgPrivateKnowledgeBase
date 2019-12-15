package router

import (
	"bgFTPExchangeService/app/api/info"
	"github.com/gogf/gf/frame/g"
)

func init() {
	g.Server().BindObject("/info", new(info.Controller))
}
