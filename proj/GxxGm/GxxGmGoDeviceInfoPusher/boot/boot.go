package boot

import (
	_ "GxxGmGoDeviceInfoPusher/app/service/device"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

func init() {
	glog.SetLevelStr(g.Config().GetString("logger.Level"))
}

