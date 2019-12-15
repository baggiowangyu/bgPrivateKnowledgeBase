package router

import (
	"bgNetTunnelB/app/api/info"
	"github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    //g.Server().BindHandler("/", hello.Handler)
    info_controller := new(info.Controller)
    g.Server().BindObject("GET:/ClientList/All", info_controller, "GetAllClients")
}
