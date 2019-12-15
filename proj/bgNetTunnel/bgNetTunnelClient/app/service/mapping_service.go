package service

import (
	"bgNetTunnelClient/library/bgMappingMgr"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

/*
映射服务

此服务主要负责启动隧道服务端
隧道服务端管理映射管理器
 */
type MappingService struct {
	Mapping_mgr bgMappingMgr.MappingMgr
}

var MappingSrv	MappingService

func (m *MappingService) Initialize() error {
	glog.Infof("MappingService::Initialize()")

	tunnel_srv_ip := g.Config().GetString("tunnel.client")
	tunnel_srv_port := g.Config().GetInt("tunnel.cli_port")
	tunnel_user := g.Config().GetString("tunnel.user")
	tunnel_pass := g.Config().GetString("tunnel.pass")
	tunnel_send_dir := g.Config().GetString("tunnel.senddir")
	tunnel_recv_dir := g.Config().GetString("tunnel.recvdir")
	tunnel_net_type := g.Config().GetString("tunnel.proto")

	glog.Infof("tunnel_srv_ip : %s", tunnel_srv_ip)
	glog.Infof("tunnel_srv_port : %d", tunnel_srv_port)
	glog.Infof("tunnel_user : %s", tunnel_user)
	glog.Infof("tunnel_pass : %s", tunnel_pass)
	glog.Infof("tunnel_send_dir : %s", tunnel_send_dir)
	glog.Infof("tunnel_recv_dir : %s", tunnel_recv_dir)
	glog.Infof("tunnel_net_type : %s", tunnel_net_type)

	// 初始化映射管理器
	err := m.Mapping_mgr.Initialize(tunnel_srv_ip, tunnel_srv_port, tunnel_user, tunnel_pass, tunnel_send_dir, tunnel_recv_dir, tunnel_net_type)
	return err
}
