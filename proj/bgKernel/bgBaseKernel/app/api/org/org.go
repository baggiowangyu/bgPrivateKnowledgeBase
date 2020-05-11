/*
本接口用于处理所有组织架构相关的操作
主要操作包括：
- 增加组织信息（单条、批量）
- 移除组织信息（单条、批量）
- 查询组织信息（根据）
- 修改组织信息（单条、批量）
*/
package org

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gvalid"
)

//
// @Summary 增加组织架构接口
// @Description 增加组织架构接口
// @Tags InsertOrgInfo
// @Success 200 {string} string	"ok"
// @router / [get]
func InsertOrgInfo(r *ghttp.Request) {
	rules := map[string]string {
		"Org_code"  : "required|length:6,16",
		"Org_name"  : "required|length:6,16",
		"Parent_id"	: "required|length:6,16",
	}

	msgs  := map[string]interface{} {
		"Org_code"	: "部门编码不能为空|账号长度应当在:min到:max之间",
		"Org_name"	: "部门名称不能为空|账号长度应当在:min到:max之间",
		"Parent_id"	: "父部门名称不能为空|账号长度应当在:min到:max之间",
	}

	request_object, err := r.GetJson()
	if err != nil {
		glog.Error(err)
		return
	}

	data_array := request_object.ToArray()
	for _, data := range data_array {
		perr := gvalid.CheckMap(data, rules, msgs)
		if perr != nil {
			glog.Info(perr)
			return
		}

		// 校验完成，开始补齐字段
		// 1、创建rid，规则rid_o_[YYYYMMDDHHmmSS]
	}


	r.Response.WriteStatus(200)
}

func GetAllOrgInfo(r *ghttp.Request) {
	r.Response.WriteStatus(200)
}
