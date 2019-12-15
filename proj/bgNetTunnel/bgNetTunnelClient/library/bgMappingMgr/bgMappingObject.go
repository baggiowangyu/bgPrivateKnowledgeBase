/*

映射对象

1、持有虚拟客户端管理器
2、负责向虚拟客户端管理器发送数据
3、接收虚拟客户端管理器发来的数据

 */
package bgMappingMgr

import (
	"bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/os/glog"
	"strconv"
)

type MappingObject struct {
	Mapping_id		int
	Mapping_ip		string
	Mapping_port	int
	Source_ip 		string
	Source_port 	int
	Net_type		string
	Is_running		int

	// 客户端管理器
	ClientMgr 		VirtualClientMgr

	// 这里还需要一个回调函数，用于向上层发送消息
	Callback		VitualClientRecvCallback
}

func (m *MappingObject) Initialize(callback VitualClientRecvCallback) error {
	m.Callback = callback

	err := m.ClientMgr.Initialize(m.SendData)
	return err
}

func (m *MappingObject) RecvData(msg bgNetMessage.NetMessageV1) error {
	// 收到数据后，在客户端表中查询对应客户端
	client_address := msg.ClientId
	client_object, err := m.ClientMgr.FindClient(m.Mapping_id, client_address)

	// 没有查询到则创建客户端，并建立与目标服务器的连接
	if err != nil {
		//没找到
		target_address := m.Source_ip + ":" + strconv.Itoa(m.Source_port)
		err = m.ClientMgr.CreateClient(m.Mapping_id, client_address, target_address, m.Net_type)
		if err != nil {
			glog.Error(err)
			return err
		}

		// 成功创建客户端，再次查询，还是没查询到的话就非常不合理了
		client_object, err = m.ClientMgr.FindClient(m.Mapping_id, client_address)
		if err != nil {
			glog.Error(err)
			return err
		}
	}

	// 解码真实请求
	data, err := gbase64.DecodeString(msg.MessageBody)
	if err != nil {
		glog.Error(err)
		return err
	}

	// 调用客户端发送接口
	err = client_object.SendData(data)
	return err
}

func (m *MappingObject) SendData(msg bgNetMessage.NetMessageV1) error {
	// 这里接收到的是目标服务端返回的真实数据, 序列化后继续往上扔
	err := m.Callback(msg)
	return err
}
