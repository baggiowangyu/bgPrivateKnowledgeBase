package dao

// 组织架构信息结构
type OrgInfo struct {
	Unique_id 	string	// 唯一ID
	Org_code  	string	// 组织编码
	Org_name  	string	// 组织名称
	Org_path	string	// 组织路径
	Parent_id 	string	// 父组织编码
	Create_time int64	// 创建时间，采用int64，排除时区问题
	Update_time int64	// 更新时间，采用int64，排除时区问题
	Source		string	// 数据来源
	Order_index	int		// 排序索引
	Extend		string	// 扩展信息
}

var Orgs map[string]*OrgInfo

func init() {
	// 直连数据库
}