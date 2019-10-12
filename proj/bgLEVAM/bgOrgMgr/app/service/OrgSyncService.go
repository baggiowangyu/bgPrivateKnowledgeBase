package service

import (
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"time"
)

type OrgSyncService struct {

}

var (
	redis_client = g.Redis()
	mysql_client = g.DB("default").Table("bg_org_info")
)

// 初始化当前的
func init() {
	println("OrgSyncService::init()")
}

func (s *OrgSyncService) Thread() error {
	// 得到线程的休眠时长，时间是毫秒，转为Duration类型/对象
	duration := g.Config().GetDuration("orgsync.duration") * time.Second

	for true {
		// 从数据库中取出所有部门数据，并写入Redis
		orgs_result, err := mysql_client.Select()
		if err == nil {
			orgs_result_json_string := orgs_result.ToJson()
			println(orgs_result_json_string)
			orgs_result_json_string_base64 := gbase64.EncodeString(orgs_result_json_string)

			redis_value, err := redis_client.DoVar("GET", "org_info")
			if err != nil {
				// 不存在，直接SET进去
				glog.Info("OrgSyncService::Thread() " + err.Error())

				_, err = redis_client.Do("SET", "org_info", orgs_result_json_string_base64)
				if err != nil {
					glog.Error("OrgSyncService::Thread() " + err.Error())
				}
			}

			// 存在，则比较结果，相同则不处理，不同则重新设置
			if orgs_result_json_string_base64 != redis_value.String() {
				_, err = redis_client.Do("SET", "org_info", orgs_result_json_string_base64)
				if err != nil {
					glog.Error("OrgSyncService::Thread() " + err.Error())
				}
			}
		} else {
			// 相关错误信息
			glog.Error("OrgSyncService::Thread() " + err.Error())
		}

		// 等待指定时长后再次同步
		time.Sleep(duration)
	}
	return nil
}

// 启动组织架构同步服务
func (s *OrgSyncService) Start() {
	// 这里主要工作就是启动一个线程，定期执行读取数据库，更新Redis的工作
	println("OrgSyncService::Start()")
	go s.Thread()
}