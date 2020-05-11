package router

import (
    "bgWSGateway/app/api/dsj"
    "bgWSGateway/app/api/extend"
    "bgWSGateway/app/api/file"
    "bgWSGateway/app/api/hello"
    "bgWSGateway/app/api/log"
    "bgWSGateway/app/api/workstation"
    "bgWSGateway/app/midware"
    "github.com/gogf/gf/frame/g"
)

// 统一路由注册.
func init() {
    g.Server().BindHandler("/", hello.Handler)

    workstation_controller := new(workstation.WorkStationController)
    g.Server().BindObject("POST:/v3/wsinfo/", workstation_controller, "Heartbeat")
    g.Server().BindMiddleware("POST:/v3/wsinfo/heartbeat", midware.MidwareAccessCheck)

    file_controller := new(file.FileController)
    g.Server().BindObject("POST:/v3/", file_controller, "Fileinfo")
    g.Server().BindMiddleware("POST:/v3/fileinfo", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", file_controller, "Can_delete_files")
    g.Server().BindMiddleware("GET:/v3/can_delete_files", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/", file_controller, "Notify_file_deleted")
    g.Server().BindMiddleware("POST:/v3/notify_file_deleted", midware.MidwareAccessCheck)

    dsj_controller := new(dsj.DsjController)
    g.Server().BindObject("POST:/v3/", dsj_controller, "Dsjinfo")
    g.Server().BindMiddleware("POST:/v3/dsjinfo", midware.MidwareAccessCheck)

    log_controller := new(log.LogController)
    g.Server().BindObject("POST:/v3/log/", log_controller, "Dsjlog")
    g.Server().BindMiddleware("POST:/v3/log/dsjlog", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/log/", log_controller, "Dsjlogfile")
    g.Server().BindMiddleware("POST:/v3/log/dsjlogfile", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/log/", log_controller, "Wslog")
    g.Server().BindMiddleware("POST:/v3/log/wslog", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/log/", log_controller, "Gpslog")
    g.Server().BindMiddleware("POST:/v3/log/gpslog", midware.MidwareAccessCheck)

    extend_controller := new(extend.ExtendController)
    g.Server().BindObject("GET:/v3/", extend_controller, "Suborg")
    g.Server().BindMiddleware("GET:/v3/suborg", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Userinfo")
    g.Server().BindMiddleware("GET:/v3/userinfo", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/", extend_controller, "Notify_registed")
    g.Server().BindMiddleware("POST:/v3/notify_registed", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Query_registed_dsj")
    g.Server().BindMiddleware("GET:/v3/query_registed_dsj", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Search_dsj")
    g.Server().BindMiddleware("GET:/v3/search_dsj", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/", extend_controller, "Registed")
    g.Server().BindMiddleware("POST:/v3/registed", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/", extend_controller, "Binduser")
    g.Server().BindMiddleware("POST:/v3/binduser", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Get_upgradepatch_list")
    g.Server().BindMiddleware("GET:/v3/get_upgradepatch_list", midware.MidwareAccessCheck)

    g.Server().BindObject("POST:/v3/", extend_controller, "Update_upgradepatch_status")
    g.Server().BindMiddleware("POST:/v3/update_upgradepatch_status", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Announcement")
    g.Server().BindMiddleware("GET:/v3/announcement", midware.MidwareAccessCheck)

    g.Server().BindObject("GET:/v3/", extend_controller, "Ping")
    //g.Server().BindMiddleware("GET:/v3/ping", midware.MidwareAccessCheck)
}
