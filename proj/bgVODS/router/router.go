package router

import (
    "bgVODS/app/api/hello"
    "github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    g.Server().BindHandler("/", hello.Handler)
}
