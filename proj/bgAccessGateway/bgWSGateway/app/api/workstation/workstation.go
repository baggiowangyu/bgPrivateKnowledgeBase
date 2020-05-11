package workstation

import (
	"bgWSGateway/app/dao"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

// 根路径 /openapi/workstation

type WorkStationController struct {

}

///////////////////////////////////////////////////////////////////////////////
/*
工作站心跳
POST /v3/wsinfo/heartbeat?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *WorkStationController) Heartbeat(r *ghttp.Request) {

	// 获取发上来的心跳
	request, err := r.GetJson()
	if err != nil {
		glog.Error(err)
		r.Response.WriteStatus(500)
		return
	}

	// 解析心跳数据
	heartbeat := new(dao.WorkStationHeartBeat)
	if request.Contains("gzz_xh") {
		heartbeat.Ws_gbcode = request.GetString("gzz_xh")
	} else if request.Contains("gzz_ipdz") {
		heartbeat.Ws_ip = request.GetString("gzz_ipdz")
	} else if request.Contains("zxzt") {
		heartbeat.Ws_online = request.GetInt("zxzt")
	} else if request.Contains("qizt") {
		heartbeat.Ws_enable = request.GetInt("qizt")
	} else if request.Contains("cczrl") {
		heartbeat.Ws_total_storage_size = request.GetInt64("cczrl")
	} else if request.Contains("syzrl") {
		heartbeat.Ws_used_storage_size = request.GetInt64("syzrl")
	} else if request.Contains("cpu") {
		heartbeat.Ws_cpu_usage = request.GetInt("cpu")
	} else if request.Contains("ram") {
		heartbeat.Ws_ram_usage = request.GetInt("ram")
	} else if request.Contains("raid_zt") {
		heartbeat.Ws_raid_state = request.GetInt("raid_zt")
	} else if request.Contains("bjlx") {
		heartbeat.Ws_alarm_type = request.GetInt("bjlx")
	} else if request.Contains("version") {
		heartbeat.Ws_version = request.GetString("version")
	}

	// 检查当前设备是否为注册设备

	// 检测通过了，从内存取出对应的
}

