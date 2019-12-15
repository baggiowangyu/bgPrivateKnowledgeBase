/*

映射服务：

主要实现以下几个功能：

1、映射表的管理（增加、删除、查询）
2、映射服务的管理（开启、停止）

下属一个客户端管理对象，实现以下功能：

1、管理客户端的连接、断开事件
2、管理客户端接入映射服务情况

下属一个隧道协议对象，实现以下功能：

1、实现单通道复用多映射连接
2、实现单通道复用多客户端连接
3、实现数据序列化算法
4、实现隧道数据机密性传输（对称加密算法[AES|SM4]，密钥管理采用PKI方式管理或固定密钥管理）

下属一个隧道对象，实现以下功能：

1、创建隧道连接（TCP、UDP、FTP、SIP...）
2、向隧道发送数据
3、从隧道接收数据

 */

package service

import (
	"bgNetTunnelServer/library/bgMappingMgr"
	"github.com/gogf/gf/frame/g"
)

var (
	// 映射服务对象实例，单例
	MappingSrv 		MappingService
)

type MappingService struct {
	Mapping_mgr bgMappingMgr.MappingMgr
}

/*

初始化映射管理器

1、读数据库，将映射表载入内存
2、启动隧道对象，用于穿越“网闸”，同时同步映射表
3、启动本地的映射服务

*/
func (m *MappingService) Initialize() error {

	tunnel_srv_ip := g.Config().GetString("tunnel.client")
	tunnel_srv_port := g.Config().GetInt("tunnel.cli_port")
	tunnel_user := g.Config().GetString("tunnel.user")
	tunnel_pass := g.Config().GetString("tunnel.pass")
	tunnel_send_dir := g.Config().GetString("tunnel.senddir")
	tunnel_recv_dir := g.Config().GetString("tunnel.recvdir")
	tunnel_net_type := g.Config().GetString("tunnel.proto")

	// 初始化映射管理器
	err := m.Mapping_mgr.Initialize(tunnel_srv_ip, tunnel_srv_port, tunnel_user, tunnel_pass, tunnel_send_dir, tunnel_recv_dir, tunnel_net_type)
	return err
}

/*
增加映射信息
 */
func (m *MappingService) AddMapping(mapping_ip string, mapping_port int, source_ip string, source_port int, net_type string) error {
	// 查找是否已经存在
	return nil
}

func (m *MappingService) RemoveMapping(source_ip string, source_port int) error {
	// 根据映射目标地址，从数据库和缓存中移除映射
	return nil
}

//func (m *MappingService) QueryMappings() (*[]MappingInfo, error) {
//	// 根据映射目标地址，查询完整映射信息
//}

