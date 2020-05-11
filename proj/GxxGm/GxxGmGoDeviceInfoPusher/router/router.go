package router

import (
	"GxxGmGoDeviceInfoPusher/app/api/subscribe"
	"GxxGmGoDeviceInfoPusher/app/api/test"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	//s.Group("/", func(group *ghttp.RouterGroup) {
	//	group.ALL("/", hello.Hello)
	//})

	s.Group("/subscribe", func(group *ghttp.RouterGroup) {
		group.POST("/Add", subscribe.AddSubScribe)
		group.POST("/Remove", subscribe.RemoveSubScribe)
	})

	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.POST("/TestRecvGps", test.TestRecvGps)
	})
}
