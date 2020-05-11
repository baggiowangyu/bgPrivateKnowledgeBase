package test

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func TestRecvGps(r *ghttp.Request) {
	request, err := r.GetJson()
	if err != nil {
		glog.Error(err)
		return
	}

	glog.Debugf("TestRecvGps() >>> \n%s", request.Export())

	r.Response.WriteStatus(200)
	return
}
