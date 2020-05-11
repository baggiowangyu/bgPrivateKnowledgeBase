package extend

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type ExtendController struct {

}

///////////////////////////////////////////////////////////////////////////////
/*
读取子部门列表
GET /v3/suborg?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]&sjbmbh=[sjbmbh]
*/
func (c *ExtendController) Suborg(r *ghttp.Request) {

}

/*
获取指定部门下的用户列表
GET /v3/userinfo?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]&bmbh=[bmbh]
*/
func (c *ExtendController) Userinfo(r *ghttp.Request) {

}

/*
设备通知注册
接口说明 工作站通过该接口访问后台进行设备通知注册。
POST /v3/deviceinfo/notify_registed?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ExtendController) Notify_registed(r *ghttp.Request) {

}

/*
设备注册信息查询接口
接口说明 根据用户编号，设备序列号和imei号查询设备注册信息
GET /v3/deviceinfo/query_registed_dsj?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]&jybh=[jybh]&sn=[sn]&imei=[imei]
*/
func (c *ExtendController) Query_registed_dsj(r *ghttp.Request) {

}

/*
查询指定设备信息接口
接口说明 根据设备编号，查询指定设备的信息
GET /v3/deviceinfo/search_dsj?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]&cpxh=[cpxh]
*/
func (c *ExtendController) Search_dsj(r *ghttp.Request) {

}

/*
设备注册(采集站静默注册模式下使用)
接口说明 工作站通过该接口访问后台进行设备通册。
POST /v3/deviceinfo/registed?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ExtendController) Registed(r *ghttp.Request) {

}

/*
设备绑定使用人(采集站静默注册模式下使用)
接口说明 工作站通过该接口访问后台进行设备通册。
POST /v3/deviceinfo/binduser?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ExtendController) Binduser(r *ghttp.Request) {

}

/*
获取工作站未升级的升级包接口
接口说明 获取工作站未升级的升级包
GET /v3/upgradepatch/get_upgradepatch_list?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ExtendController) Get_upgradepatch_list(r *ghttp.Request) {

}

/*
提交工作站升级包完成状态
接口说明 提交工作站升级包完成状态
POST /v3/upgradepatch/update_upgradepatch_status?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ExtendController) Update_upgradepatch_status(r *ghttp.Request) {

}

/*
获取采集站公告信息接口
接口说明 获取指定采集站公告信息接口
GET /v3/announcement?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]&start=[start]&limit=[limit]
*/
func (c *ExtendController) Announcement(r *ghttp.Request) {

}

/*
Ping 连通性检查
接口说明 工作站通过该接口验证服务器是否正常。请求url上不要求附带gzz_xh、authkey、domain参数。
GET /v3/ping
*/
func (c *ExtendController) Ping(r *ghttp.Request) {
	result := make(map[string]interface{}, 0)
	result["code"] = 0
	result["message"] = "SUCCESS"

	json_result := gjson.New(result)

	err := r.Response.WriteJson(json_result)
	if err != nil {
		glog.Error(err)
	}
}
