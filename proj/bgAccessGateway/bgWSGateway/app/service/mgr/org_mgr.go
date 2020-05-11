// 组织架构同步管理
//
// 本模块用于定期从MySQL或MariaDB或PostgreSQL等关系型数据库中获取最新的组织架构情况
//
package mgr

import "github.com/gogf/gf/frame/g"

type OrgMgr struct {

}

// 在这里根据配置创建一个定时任务，用于执行组织架构信息同步
func (o *OrgMgr) Initialize() error {
	g.DB("auth-cas").A
}