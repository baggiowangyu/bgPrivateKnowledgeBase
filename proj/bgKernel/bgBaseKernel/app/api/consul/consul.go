/*

*/
package consul

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	consulapi "github.com/hashicorp/consul/api"
	"time"
)

func InitializeConsul() {
	consul_address := g.Config().GetString("consul.address")

	consul_config := consulapi.DefaultConfig()
	consul_config.Address = consul_address

	consul_register := consulapi.AgentServiceRegistration{}
	consul_register.Name = g.Config().GetString("server.ServerAgent")
	consul_register.Address = g.Config().GetString("consul.local_ip")
	consul_register.Port = g.Config().GetInt("consul.local_port")
	consul_register.Tags = []string{"primary"}

	consul_check := consulapi.AgentServiceCheck{}
	consul_check.Interval = g.Config().GetString("consul.check_interval")
	consul_check.HTTP = g.Config().GetString("consul.check_url")
	consul_register.Check = &consul_check

	// 创建Consul客户端
	consul_client, err := consulapi.NewClient(consul_config)
	if err != nil {
		glog.Error(err)
		return
	}

	// 注册到Consul，如果注册不上就尝试重复注册
	for {
		err = Register(consul_client, consul_register)
		if err != nil {
			glog.Error(err)
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

}

func Register(client *consulapi.Client, register consulapi.AgentServiceRegistration) error {
	return client.Agent().ServiceRegister(&register)
}

func Health(r *ghttp.Request) {
	json_object := gjson.New("{\"status\": \"ok\"}")
	err := r.Response.WriteJson(json_object)
	if err != nil {
		glog.Error(err)
	}
}