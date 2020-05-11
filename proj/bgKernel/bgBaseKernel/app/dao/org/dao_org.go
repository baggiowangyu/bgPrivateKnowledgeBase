package org

type InsertOrgRequest struct {
	Org_code	string `v:"required#部门编码不能为空"`
	Org_name	string `v:"required#部门编码不能为空"`
	Org_path	string `v:"required#部门编码不能为空"`
	Parent_id	string `v:"required#部门编码不能为空"`
	Order_index	int
	Is_hide		int
	Is_disable	int
	Extend		string
}

type OrgInfo struct {
	Unique_id	string
	Org_code	string
	Org_name	string
	Org_path	string
	Parent_id	string
	Create_time	int64
	Update_time	int64
	Source		string
	Order_index	int
	Is_hide		int
	Is_disable	int
	Extend		string
}


