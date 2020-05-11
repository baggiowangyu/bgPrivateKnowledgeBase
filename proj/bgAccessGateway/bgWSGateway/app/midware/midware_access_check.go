package midware

import (
	"github.com/gogf/gf/net/ghttp"
)

func MidwareAccessCheck(r *ghttp.Request) {
	//request_domain := r.GetQueryString("domain")
	//request_workstation_id := r.GetQueryString("gzz_xh")
	//request_authkey := r.GetQueryString("authkey")

	// 主要是校验workstation_id对应设备的domain和authkey
	// 这些数据应当都缓存在内存中
	// 如果校验失败，则直接在这里返回了
	r.Middleware.Next()
}
