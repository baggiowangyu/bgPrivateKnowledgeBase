package boot

import (
    "bgVODS/app/service/file_mgr"
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/net/ghttp"
    "github.com/gogf/gf/os/glog"
)

var Global_vods_root string = g.Config().GetString("vods.root_path")

func initConfig() {
    log_path := g.Config().GetString("app.log_path")
    log_level := g.Config().GetInt("app.log_level")

    glog.SetLevel(log_level)
    glog.SetPath(log_path)
    glog.SetStdoutPrint(true)

    http_port := g.Config().GetInt("app.http_port")
    g.Server().SetPort(http_port)
    g.Server().SetLogPath(log_path)
    g.Server().SetAccessLogEnabled(true)
    g.Server().SetLogStdout(true)
    g.Server().SetNameToUriType(ghttp.URI_TYPE_FULLNAME)

    if g.Config().GetBool("app.enable_log_stdout") {
        glog.SetStdoutPrint(true)
        g.Server().SetLogStdout(true)
    } else {
        glog.SetStdoutPrint(false)
        g.Server().SetLogStdout(false)
    }

    if g.Config().GetBool("app.enable_https") {
        g.Server().EnableHTTPS("cert/server.crt", "cert/server.key")
    }
}

func init() {
    // 初始化配置信息
    initConfig()

    // 初始化文件管理模块
    err := file_mgr.VodsFileMgr.Initialize(Global_vods_root)
    if err != nil {
        glog.Error(err)
    }
}

