package router

import (
	"bgVODS/app/api/vods"
	"github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    //g.Server().BindHandler("/", hello.Handler)
	g.Server().BindHandler("GET:/*any", vods.PlayHandlerEx)
}
