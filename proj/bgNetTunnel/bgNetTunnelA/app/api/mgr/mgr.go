package mgr

import (
	"bgNetTunnelA/app/service/mapping_srv"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"net/http"
)

type Controller struct {
	
}

/**
展示映射表信息
 */
func (c *Controller) MappingTableInfo(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
	mapping_list := g.ArrayAny{}
	for _, mapping_object := range mapping_srv.Mapping_service_instance.MappingTable {
		mapping_list = append(mapping_list, mapping_object.Info)
	}

	mapping_info_json := gjson.New(mapping_list).ToArray()

	err := r.Response.WriteTpl("index.html", g.Map{
		"title" : "网络隧道A端",
		"table_name" : "映射表信息",
		"Mapping_id_title" : "映射ID",
		"Mapping_port_title" : "映射端口",
		"Source_ip_title" : "目标IP",
		"Source_port_title" : "目标端口",
		"Net_type_title" : "网络类型",
		"Is_running_title" : "运行状态",
		"data" : mapping_info_json,
	})

	if err != nil {
		glog.Error(err)
	}
}

/**
添加映射表
 */
func (c *Controller) AddMappingTable(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
}

/**
移除映射表
 */
func (c *Controller) RemoveMappingTable(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
}

/**
启动映射
 */
func (c *Controller) StartMappingObject(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
}

/**
停止映射
 */
func (c *Controller) StopMappingObject(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
}

/**
查看映射表连接信息
 */
func (c *Controller) GetMappingObjectConnection(r *ghttp.Request) {
	// 查找所有映射表信息，通过模板返回给页面
	mapping_id := r.GetQuery("id").(int)
	mapping_object, exists := mapping_srv.Mapping_service_instance.MappingTable[mapping_id]
	if !exists {
		r.Response.WriteStatus(http.StatusNotFound)
	} else {
		if mapping_object.Info.Net_type == "TCP" {
			// 将所有客户端信息打印出来
			//mapping_object.Tcp_client_table
		} else if mapping_object.Info.Net_type == "UDP" {
			// 将所有UDP客户端打印出来
		}
	}

	r.Response.Write(mapping_id)
}