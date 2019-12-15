package api_mgr

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

type ApiMgrService struct {
	/*
	api列表
	key 	: Api标签 >>> [Method:URI]
	value	: Api对象
	*/
	Api_list map[string]*ApiObject
}

var ApiMgrController ApiMgrService

func (a *ApiMgrService) Initialize() error {
	var err error

	a.Api_list = make(map[string]*ApiObject, 0)

	return err
}

type Controller struct {

}

func (c *Controller) AddApi(r *ghttp.Request) {
	// 数据校验
	request_json, err := r.GetJson()
	if err != nil {
		glog.Error(err)
		r.Response.WriteStatus(500)
	}

	println(request_json.Export())
}

func (a *ApiMgrService) RemoveApi(r *ghttp.Request) {
	// 数据校验
}

func (a *ApiMgrService) QueryAllApi(r *ghttp.Request) {
	// 数据校验
}

func (a *ApiMgrService) ModifyApi(r *ghttp.Request) {
	// 数据校验
}
