package common

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"bgCMS/module/constants"
)

// Login 登录页面
func Welcome(r *ghttp.Request) {
	err := r.Response.WriteTpl("pages/welcome.html", g.Map{
		"app_name" : "执法视音频一体化管理平台",
	})

	if err != nil {
		glog.Error(err)
	}
}

// Login 登录页面
func Debug(r *ghttp.Request) {
	if constants.DEBUG {
		constants.DEBUG = false
		g.DB().SetDebug(false)
		r.Response.Writeln("debug close ~!~ ")
	} else {
		constants.DEBUG = true
		g.DB().SetDebug(true)
		r.Response.Writeln("debug open ~!~ ")
	}
}
