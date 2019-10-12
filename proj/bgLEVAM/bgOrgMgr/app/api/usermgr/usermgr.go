package usermgr

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
)

type Controller struct {

}

// 添加用户
func (c *Controller) AddUser(r *ghttp.Request) {

}

// 修改用户
func (c *Controller) ModifyUser(r *ghttp.Request) {

}

// 删除用户
func (c *Controller) DeleteUser(r *ghttp.Request) {

}

// 查询用户
// 查询条件：{"org_name": "", "contain_sub_orgs": True, }
func (c *Controller) QueryUser(r *ghttp.Request) {
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
