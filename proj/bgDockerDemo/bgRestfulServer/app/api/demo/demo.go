package demo

import "github.com/gogf/gf/net/ghttp"

type Controller struct {

}

func (c *Controller) Hello(r *ghttp.Request) {
	r.Response.Write("Hello! Welcome to docker service...")
}