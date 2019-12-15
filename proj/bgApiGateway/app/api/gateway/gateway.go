package gateway

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func Interface(r *ghttp.Request) {
	result := make(map[string]string, 0)
	result["Method"] = r.Request.Method
	result["RequestURI"] = r.Request.RequestURI

	// 映射标签
	//api_tag := r.Request.Method + ":" + r.Request.RequestURI

	// 根据映射标签查找映射表

	//result["Method"] = r.Request.URL.Path
	//result["Method"] = r.Request.Method
	//result["Method"] = r.Request.Method
	j := gjson.New(result)
	//r.Response.Write(r.Router)
	err := r.Response.WriteJson(j)
	if err != nil {
		glog.Error(err)
	}
}
