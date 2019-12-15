package router

import (
	"bgRestfulServer/app/api/demo"
	"github.com/gogf/gf/frame/g"
)

func init() {
	g.Server().BindObject("GET:/demo", new(demo.Controller), "Hello")
}