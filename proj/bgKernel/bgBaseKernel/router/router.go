package router

import (
	"bgBaseKernel/app/api/auth"
	"bgBaseKernel/app/api/consul"
	"bgBaseKernel/app/api/middleware"
	"bgBaseKernel/app/api/org"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func init() {
	s := g.Server()
	s.Group("/api.v4", func(group *ghttp.RouterGroup) {

		// 这里增加权限处理以及中间件过滤
		group.Middleware(auth.ChechAuth, middleware.MiddlewareCORS)

		// 定义组织架构相关API路由
		s.Group("/org", func(group *ghttp.RouterGroup) {
			group.POST("/InsertOrgInfo", org.InsertOrgInfo)

			group.GET("/GetAllOrgs", org.GetAllOrgInfo)
		})
	})

	// 设置连接到Consul
	s.BindHandler("GET:/consul/health", consul.Health)

}
