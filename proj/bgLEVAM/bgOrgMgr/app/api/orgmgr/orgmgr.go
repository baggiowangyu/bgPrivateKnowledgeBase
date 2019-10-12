package orgmgr

import (
	"../../service"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

var (
	org_sync_service = service.OrgSyncService{}
	bg_org_info = g.DB().Table("bg_org_info")
)

func init() {
	// 同步服务开始工作
	org_sync_service.Start()
}

type Controller struct {

}


// 添加组织（部门）
func (c *Controller) AddOrg(r *ghttp.Request) {
	// 获取POST上来的数据
	//var post_data_json *gjson.Json
	post_data_json, err := r.GetJson()
	if err != nil {
		println(err.Error())
		return
	}
	println(post_data_json)

	// 转换成数组
	json_array_element := post_data_json.ToArray()

	// 开启数据库事务
	tx, err := g.DB().Begin()
	for _, element := range json_array_element {
		// 插入
		_, err := bg_org_info.Data(element).Insert()
		if err != nil {
			// 出现错误了，跳出
			break
		}
	}

	if err != nil {
		// 前面插入数据的时候出错了，这里回滚，记录日志
		glog.Error(err.Error())
		err = tx.Rollback()
		r.Response.Write("{\"status\": \"" + err.Error() + "\"}")
	} else {
		// 提交事务
		err = tx.Commit()
		r.Response.Write("{\"status\": \"OK\"}")
	}


}

// 修改组织（部门）
func (c *Controller) ModifyOrg(r *ghttp.Request) {
	//// 获取POST上来的数据
	//var post_data_json *gjson.Json
	//post_data_json, err := r.GetJson()
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
}

// 删除部门（部门）
func (c *Controller) DeleteOrg(r *ghttp.Request) {
	//// 获取POST上来的数据
	//var post_data_json *gjson.Json
	//post_data_json, err := r.GetJson()
	//if err != nil {
	//	println(err.Error())
	//	return
	//}
}

// 查询组织（部门）
// 查询条件：{"org_name": "", "contain_sub_orgs": True, }
func (c *Controller) QueryOrg(r *ghttp.Request) {
	// 获取POST上来的数据
	var post_data_json *gjson.Json
	post_data_json, err := r.GetJson()
	if err != nil {
		println(err.Error())
		return
	}

	org_name := post_data_json.GetString("org_name")
	contain_sub_orgs := post_data_json.GetBool("contain_sub_orgs")
	if (org_name == "ROOT") && (contain_sub_orgs == true) {
		// 从根部门开始，查询所有部门信息
		// 先尝试从缓存查，如果缓存没有，则查询数据库，然后将结果写入缓存
	} else if (org_name == "ROOT") && (contain_sub_orgs == false) {
		// 只查询根部门信息
	}
}
