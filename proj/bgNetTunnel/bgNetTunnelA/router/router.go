package router

import (
    "bgNetTunnelA/app/api/mgr"
    "github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    //g.Server().BindHandler("/", hello.Handler)
    mgr_controller := new(mgr.Controller)
    g.Server().BindObject("GET:/MappingTableInfo", mgr_controller, "MappingTableInfo")
    g.Server().BindObject("GET:/MappingTableInfo", mgr_controller, "AddMappingTable")
    g.Server().BindObject("POST:/MappingTableInfo", mgr_controller, "AddMappingTableHandler")

    g.Server().BindObject("GET:/MappingObject", mgr_controller, "GetMappingObjectConnection")
}
