package info

import "github.com/gogf/gf/net/ghttp"

type Controller struct {

}

func (c *Controller) info(r *ghttp.Request) {
	r.Response.Write("bgFTPExchanger")
}