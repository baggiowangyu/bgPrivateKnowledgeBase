package sip

import (
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	"time"
)

/*
本模块定义了注册超期检查任务
本模块定义了心跳超期检查任务
*/

func init() {
	gcron.Add("0-59/5 * * * * *", RegisterExpiredCheckTask, "RegisterExpiredCheckTask")
	time.Sleep(1 * time.Millisecond)

	gcron.Add("0-59/5 * * * * *", HeartbeatExpiredCheckTask, "HeartbeatExpiredCheckTask")
	time.Sleep(1 * time.Millisecond)
}

func RegisterExpiredCheckTask() {
	for client_gbcode, sip_client := range sip_server.Sip_clients {
		// 计算注册有效期是否过期
		// 有效期计算方法：当前时间-注册时间 > 注册有效期，则将此对象删除
		current_time := gtime.Now().Unix()
		interval := current_time - sip_client.register_time
		if interval > sip_client.expired_time {
			glog.Debugf("[%s] 注册超期，强制踢掉", client_gbcode)
			// 已经超期了，这里将对象从map移除
			delete(sip_server.Sip_clients, client_gbcode)

			// 如有必要，向观察者推送下线通知
		}
	}
}

func HeartbeatExpiredCheckTask() {
	for client_gbcode, sip_client := range sip_server.Sip_clients {
		current_time := gtime.Now().Unix()
		interval := current_time - sip_client.latest_heartbeat_time
		if interval > Heartbeat_check_timeout {
			glog.Debugf("[%s] 心跳超期，强制踢掉", client_gbcode)
			// 已经超期了，这里将对象从map移除
			delete(sip_server.Sip_clients, client_gbcode)

			// 如有必要，向观察者推送下线通知
		}
	}
}