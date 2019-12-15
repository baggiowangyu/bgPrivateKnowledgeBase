package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"bgCMS/module/component/started"
)

// 管理初始化顺序.
func init() {
	// 初始化配置信息
	initConfig()

	// 初始化路由信息
	initRouter()

	// 初始化数据信息
	started.Start()
}

// 用于配置初始化.
func initConfig() {
	glog.Info("########service start...")

	v := g.View()
	c := g.Config()
	s := g.Server()

	// 配置对象及视图对象配置
	c.AddPath("config")

	v.SetDelimiters("${", "}")
	v.AddPath("template")

	// glog配置
	logPath := c.GetString("log-path")
	glog.SetPath(logPath)
	glog.SetStdoutPrint(true)

	s.SetServerRoot("public")
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.SetLogPath(logPath)
	s.SetErrorLogEnabled(true)
	//s.SetAccessLogEnabled(true)
	s.SetPort(c.GetInt("http-port"))

	glog.Info("########service finish.")

}
