package router

import (
	"../app/api/devicemgr"
	"github.com/gogf/gf/frame/g"
)

func init() {
	// 设备管理模块 路由注册 - 使用执行对象注册方式
	g.Server().BindObject("/devicemgr", new(devicemgr.Controller))
}