/*

映射管理接口：

1、查询所有映射信息
2、添加映射
3、删除映射
4、启用映射
5、禁用映射

 */

package mapping

import (
	"bgNetTunnelServer/library/bgMappingMgr"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type Controller struct {

}

func (c *Controller) QueryAllMappings(r *ghttp.Request) {
	// 先组列表
	mapping_info_array := make([]bgMappingMgr.MappingInfo, 0)
	for _, mapping_info := range bgMappingMgr.MappingTable {
		mapping_info_array = append(mapping_info_array, mapping_info)
	}

	// 再组map
	result_map := make(map[string]interface{}, 0)
	result_map["code"] = 0
	result_map["status"] = "OK"
	result_map["data"] = mapping_info_array

	// 最后Json化
	json_object := gjson.New(result_map)
	//result, err := json_object.ToJsonString()

	//if err != nil {
	//	r.Response.WriteHeader(404)
	//
	//	//_ = json_object
	//	//_ = result_map
	//	//_ = mapping_info_array
	//	return
	//}

	err := r.Response.WriteJson(json_object)
	if err != nil {
		glog.Error(err)
	}
	//_ = json_object
	//_ = result_map
	//_ = mapping_info_array
	return
}

func (c *Controller) AddMapping(r *ghttp.Request) {
	return
}

func (c *Controller) RemoveMapping(r *ghttp.Request) {
	return
}

func (c *Controller) StartMapping(r *ghttp.Request) {
	return
}

func (c *Controller) StopMapping(r *ghttp.Request) {
	return
}