package router

import (
    "bgApiGateway/app/api/gateway"
    "bgApiGateway/app/service/api_mgr"
    "github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    // 网关路由
    g.Server("Gateway").BindHandler("GET:/*any", gateway.Interface)

    // 管理器路由
    apimgrcontroller := new(api_mgr.Controller)
    g.Server("Manager").BindObject("POST:/Manager", apimgrcontroller, "AddApi")
}
