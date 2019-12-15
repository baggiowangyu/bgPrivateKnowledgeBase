package login

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type Controller struct {

}

func (c *Controller) Entry(r *ghttp.Request) {
	// 首先检查用户Session
	err := r.Response.WriteTpl("login/index.tpl")
	if err != nil {
		glog.Error(err)
	}
}

func (c *Controller) Login(r *ghttp.Request) {
	// 获得表单数据
	err := r.PostForm
	if err != nil {
		glog.Error(err)
	}
	url_values := r.Request.Form
	println(url_values)
}
