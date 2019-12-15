package boot

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func init() {
	config := g.Config()
	http_view := g.View()
	http_server := g.Server()
	tcp_server := g.TCPServer()

	// 模板引擎配置
	_ = http_view.AddPath("template")
	http_view.SetDelimiters("${", "}")

	// glog配置
	logpath := config.GetString("setting.logpath")
	_ = glog.SetPath(logpath)
	glog.SetStdoutPrint(true)

	// Web Server配置
	http_server.SetServerRoot("public")
	http_server.SetLogPath(logpath)
	http_server.SetNameToUriType(ghttp.URI_TYPE_ALLLOWER)
	http_server.SetErrorLogEnabled(true)
	http_server.SetAccessLogEnabled(true)
	http_server.SetPort(config.GetInt("setting.port"))

	// TCP Server配置
	tcp_address := config.GetString("bridge.host") + config.GetString("bridge.port")
	tcp_server.SetAddress(tcp_address)

}