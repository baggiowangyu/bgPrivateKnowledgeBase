package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func init() {
	v := g.View()
	c := g.Config()
	s := g.Server()

	// 模板引擎配置
	v.AddPath("template")
	v.SetDelimiters("${", "}")

	// glog配置
	logpath := c.GetString("setting.logpath")
	glog.SetPath(logpath)
	glog.SetStdoutPrint(true)
	glog.SetLevel(c.GetInt("setting.loglevel"))
	//level := glog.LEVEL_PROD|glog.LEVEL_INFO
	//println(level)


	// Web Server配置
	s.SetServerRoot("public")
	s.SetLogPath(logpath)
	s.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	s.SetErrorLogEnabled(true)
	s.SetAccessLogEnabled(true)
	s.SetPort(c.GetInt("setting.port"))
}