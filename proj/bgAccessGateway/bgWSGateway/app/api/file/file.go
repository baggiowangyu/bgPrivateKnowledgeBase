package file

import "github.com/gogf/gf/net/ghttp"

type FileController struct {

}

///////////////////////////////////////////////////////////////////////////////
/*
文件记录上传
POST /v3/fileinfo?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *FileController) Fileinfo(r *ghttp.Request) {

}

/*
获取工作站可删除文件列表
GET /v3/can_delete_files?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *FileController) Can_delete_files(r *ghttp.Request) {

}

/*
通知后台文件已删除
POST /v3/notify_file_deleted?domain=[domain]&gzz_xh=[gzz_xh]&authkey=[authkey]
*/
func (c *FileController) Notify_file_deleted(r *ghttp.Request) {

}
