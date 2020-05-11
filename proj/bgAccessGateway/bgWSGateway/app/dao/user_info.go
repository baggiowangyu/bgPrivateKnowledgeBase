package dao

// 用户信息结构
type UserInfo struct {
	Unique_id 		string	// 唯一ID
	User_code  		string	// 用户编码（相当于工号，警察为警号，辅警为辅警号，其他人为相对应工号），也作为系统登录账号
	User_pin		string	// 用户口令
	User_name  		string	// 用户名称
	Org_id			string	// 所在组织ID
	Role_id			string	// 系统角色（系统角色默认有三个：系统管理员、安全管理员、审计管理员）
	Is_disable		bool	// 是否被标记为禁用
	Is_deleted		bool	// 是否被标记为删除
	Is_multi_login	bool	// 是否允许多重登录
	Gender			int		// 性别：0-女；1-男；2-其他
	Create_time		int64	// 创建时间
	Update_time		int64	// 更新时间
}
