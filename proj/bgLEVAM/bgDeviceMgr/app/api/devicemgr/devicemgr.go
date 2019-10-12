package devicemgr

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
)

// 定义一个控制器结构
type Controller struct { }


// 添加设备信息
func (c *Controller) AddDevice(r *ghttp.Request) {
	// 获取POST上来的数据
	//post_data := r.GetRawString()
	//println(post_data)

	// 直接取Json对象
	var post_data_json *gjson.Json
	post_data_json, err := r.GetJson()
	if err != nil {
		println(err.Error())
		return
	}

	//// 从Json对象取出对应的值
	//device_id := post_data_json.GetString("device_id")
	//device_name := post_data_json.GetString("device_name")
	//println("device_id:" + device_id)
	//println("device_name:" + device_name)

	//post_data_json_string, err := post_data_json.ToJsonString()
	//println(post_data_json_string)

	// 转换为对象，执行ORM
}

// 修改设备信息
func (c *Controller) ModifyDeivce(r *ghttp.Request) {
	// 获取POST上来的数据
	var post_data_json *gjson.Json
	post_data_json, err := r.GetJson()
	if err != nil {
		println(err.Error())
		return
	}
}

// 查询设备信息
// 提供几种查询条件
func (c *Controller) QueryDevice(r *ghttp.Request) {
	// 获取POST上来的数据
	var post_data_json *gjson.Json
	post_data_json, err := r.GetJson()
	if err != nil {
		println(err.Error())
		return
	}
}
