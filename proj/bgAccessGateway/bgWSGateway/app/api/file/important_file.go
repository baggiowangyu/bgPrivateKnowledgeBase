package file

import "github.com/gogf/gf/net/ghttp"

type ImportantFileController struct {

}

///////////////////////////////////////////////////////////////////////////////
/*
获取需要上传的重要文件任务
GET /v3/need_upload_files?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ImportantFileController) need_upload_files(r *ghttp.Request) {

}

/*
提交重要文件上传任务状态
POST /v3/notify_upload_status?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *ImportantFileController) notify_upload_status(r *ghttp.Request) {

}
