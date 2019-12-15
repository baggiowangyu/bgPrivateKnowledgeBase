package info

import (
	"bgNetTunnelB/app/service/tunnel_srv"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/net/ghttp"
)

type Controller struct {

}

func (c *Controller) GetAllClients(r *ghttp.Request) {
	// 拿到所有客户端列表，做成Json返回
	objects_array := tunnel_srv.Tunnel_server.Client_mgr.GetAllClients()
	data := gjson.New(objects_array)
	r.Response.WriteJson(data)
}
