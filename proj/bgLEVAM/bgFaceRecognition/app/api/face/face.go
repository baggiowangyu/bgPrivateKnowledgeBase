package face

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type Controller struct {

}

func init() {
	glog.Info("face controller init()")
}

func (c *Controller) PostCompareResult(r *ghttp.Request) {
	// 接收比对结果
	compare_result, err := r.GetJson()

	if err != nil {
		glog.Error(err.Error())

		err_result_map := make(map[string]interface{}, 3)
		err_result_map["code"] = -1
		err_result_map["status"] = "Recv request data failed."
		err_result_json := gjson.New(err_result_map)
		_ = r.Response.WriteJson(err_result_json)

		_ = err_result_map
		_ = err_result_json
		return
	}

	compare_result_map := compare_result.ToMap()
	_ = compare_result_map

	// 首先将结果发送给告警管理模块(dubbo或者qpid)

	// 发送QPID给告警管理部分，通知前端界面显示人脸布控告警

	glog.Info("recv compare result ...")
	r.Response.Write("{	\"code\": 0, \"status\": \"OK\"}")
}
