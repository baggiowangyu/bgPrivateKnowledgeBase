/*
本接口用于处理登录请求、登出请求、登录验证
*/
package auth

import "github.com/gogf/gf/net/ghttp"

type LoginRequest struct {
	Username	string	`v:"required|length:1,16#账号不能为空|账号长度应当在:min到:max之间"`
	Password	string	`v:"required|"`
}

func Login(r *ghttp.Request) {
	// 检查输入参数的有效性
	var login_request *LoginRequest
	err := r.Parse(&login_request)
	if err != nil {
		//
	}
}

func Logout(r *ghttp.Request) {

}

func ChechAuth(r *ghttp.Request) {

}