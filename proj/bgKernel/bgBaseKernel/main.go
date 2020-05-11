package main

import (
	_ "bgBaseKernel/boot"
	_ "bgBaseKernel/docs"
	_ "bgBaseKernel/router"
	"github.com/gogf/gf/frame/g"

	//"github.com/zhwei820/gogf-swagger"
	//"github.com/zhwei820/gogf-swagger/swaggerFiles"
	//
	//_ "github.com/zhwei820/gogf-swagger/example/docs"
)

// @title bgBaseKernel API 说明
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1
// @BasePath ""
func main() {
	g.Server().Run()
}
