package person

import (
	//"encoding/json"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"strconv"
)

var (
	key_person_lib_table = g.DB().Table("uap_tb_gmvcs_key_person_type_db")
	key_person_info_table = g.DB().Table("uap_tb_gmvcs_key_person_recognition")
	faceserver_prefix = g.Config().GetString("faceserver.prefix")
)

type Controller struct {

}

func init() {
	glog.Info("face controller init()")
}

func (c *Controller) List(r *ghttp.Request) {
	// 先拿到想要查询的数据库信息
	key_person_lib_id := r.GetQueryInt("id")
	glog.Debug("收到查询重点人口库\"" + strconv.Itoa(key_person_lib_id) + "\"的请求")

	// 查库，或者将所有人脸库缓存到内存中，gdb链式操作
	// 这里需要考虑一下，似乎一个表对象的链式操作是一直延续下去的，其中的条件部分会不停的累加
	// 每次重新获取Table就可以了
	key_person_results, err := g.DB().Table("uap_tb_gmvcs_key_person_recognition").
		Where("key_person_type_db_id=?", key_person_lib_id).Select()
	//key_person_results, err := key_person_info_table.Where("key_person_type_db_id=?", key_person_lib_id).Select()

	if err != nil {
		// 发生错误了，这里处理一下
		glog.Error(err.Error())

		err_result_map := make(map[string]interface{}, 3)
		err_result_map["code"] = -1
		err_result_map["status"] = "Query database failed."
		err_result_json := gjson.New(err_result_map)
		_ = r.Response.WriteJson(err_result_json)

		// 牛逼的内存释放，Go的GC真是屌！
		_ = err_result_map
		_ = err_result_json
		return
	}

	json_element_array := make([]interface{}, 0)
	for _, record := range key_person_results {
		record["person_reg_img"].Set(faceserver_prefix + record["person_reg_img"].String())
		json_element_string := record.ToJson()
		json_element := gjson.New(json_element_string)
		json_element_array = append(json_element_array, json_element)
		_ = json_element
	}

	// 先组map
	result_map := make(map[string]interface{}, 0)
	result_map["code"] = 0
	result_map["status"] = "OK"
	result_map["data"] = json_element_array

	result_json := gjson.New(result_map)

	// 组织结果，返回
	_ = r.Response.WriteJson(result_json)

	_ = result_map
	_ = result_json
	_ = key_person_info_table
	return
}

//func (c *Controller) GetImages(r *ghttp.Request) {
//	// 根据传上来的人员ID，查找人脸路径
//	request_json, err := r.GetJson()
//	if err != nil {
//		// 获取上传上来的Json出错了
//		glog.Error(err.Error())
//
//		err_result_map := make(map[string]interface{}, 3)
//		err_result_map["code"] = -1
//		err_result_map["status"] = "OK"
//		err_result_json := gjson.New(err_result_map)
//		r.Response.WriteJson(err_result_json)
//
//		// 牛逼的内存释放，Go的GC真是屌！
//		_ = err_result_map
//		return
//	}
//
//
//	face_img_array := make([]string, 0)
//	for element := range request_json.Get("person_list").ToArray() {
//		// 取出的是人脸ID，那么我们去
//	}
//
//	// Base64Encode
//
//	// 组织Json
//
//	// 返回
//	glog.Debug("")
//}
