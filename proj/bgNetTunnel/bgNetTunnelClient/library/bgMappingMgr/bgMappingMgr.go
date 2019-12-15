/*

映射管理器：

1、持有网络隧道服务端，接管隧道数据的收发
2、负责接收隧道客户端发过来的控制命令（主要用于管理映射）
3、负责管理映射对象，向映射对象发送业务数据，接收映射对象发回的业务数据

 */
package bgMappingMgr

import (
	"bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"
	"bgNetTunnelClient/library/bgNetTunnel"
	"errors"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/os/glog"
	"github.com/golang/protobuf/proto"
)

type MappingInfo struct {
	Mapping_id		int
	Mapping_object	*MappingObject
}

type MappingMgr struct {
	// 网络隧道服务端
	Net_tunnel_server bgNetTunnel.NetTunnelServer

	// 映射表
	MappingTable		map[int32]*MappingObject
}

func (m *MappingMgr) Initialize(tunnel_srv_ip string, tunnel_srv_port int, tunnel_user string, tunnel_pass string, tunnel_send_dir string, tunnel_recv_dir string, tunnel_net_type string) error {

	// 主要是启动隧道对象
	err := m.Net_tunnel_server.Initialize(tunnel_srv_ip, tunnel_srv_port, tunnel_user, tunnel_pass, tunnel_send_dir, tunnel_recv_dir, tunnel_net_type, m.TunnelRecvCallback)
	if err != nil {
		glog.Debug("NetTunnelServer::Initialize() failed.")
		glog.Error(err)
	} else {
		glog.Debug("NetTunnelServer::Initialize() succeed.")
	}

	// 初始化映射表对象
	glog.Debug("Initialize mapping table.")
	m.MappingTable = make(map[int32]*MappingObject, 0)

	return err
}

func (m *MappingMgr) TunnelRecvCallback(data []byte) error {
	// 还原协议对象
	msg_v1 := bgNetMessage.NetMessageV1{}
	err := proto.Unmarshal(data, &msg_v1)
	if err != nil {
		glog.Debug("MappingMgr::TunnelRecvCallback unmarshal NetMessageV1 failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("MappingMgr::TunnelRecvCallback unmarshal NetMessageV1 failed.")
	}

	// 如果是控制协议，扔给协议客户端处理
	if msg_v1.MsgType == 0 {
		// 控制协议，我们用Json定义吧，一般分为：
		// 1、更新映射表（启动新的映射、停止映射）
		// 2、
		glog.Debug("MappingMgr::TunnelRecvCallback recv control msg.")

		msg_ctrl_info_v1_bytes, err := gbase64.DecodeString(msg_v1.MessageBody)
		if err != nil {
			glog.Debug("MappingMgr::TunnelRecvCallback base64decode NetMessageV1::MessageBody failed.")
			glog.Error(err)
			return err
		} else {
			glog.Debug("MappingMgr::TunnelRecvCallback base64decode NetMessageV1::MessageBody succeed.")
		}

		msg_ctrl_info_v1 := bgNetMessage.MsgControlInfoV1{}
		err = proto.Unmarshal(msg_ctrl_info_v1_bytes, &msg_ctrl_info_v1)
		if err != nil {
			glog.Debug("MappingMgr::TunnelRecvCallback unmarshal MsgControlInfoV1 failed.")
			glog.Error(err)
			return err
		} else {
			glog.Debug("MappingMgr::TunnelRecvCallback unmarshal MsgControlInfoV1 succeed.")
		}

		if msg_ctrl_info_v1.MainType == bgNetMessage.MsgControlInfoV1_SyncMappingTable {

			// 反序列化出映射信息
			mapping_info_v1_data, err := gbase64.DecodeString(msg_ctrl_info_v1.CtrlCommand)

			if err != nil {
				glog.Debug("MappingMgr::TunnelRecvCallback base64decode MsgControlInfoV1::CtrlCommand failed.")
				glog.Error(err)
				return err
			} else {
				glog.Debug("MappingMgr::TunnelRecvCallback base64decode MsgControlInfoV1::CtrlCommand succeed.")
			}

			mapping_info_v1 := bgNetMessage.MsgMappingInfoV1{}
			err = proto.Unmarshal(mapping_info_v1_data, &mapping_info_v1)

			if err != nil {
				glog.Debug("MappingMgr::TunnelRecvCallback unmarshal MsgMappingInfoV1 failed.")
				glog.Error(err)
				return err
			} else {
				glog.Debug("MappingMgr::TunnelRecvCallback unmarshal MsgMappingInfoV1 succeed.")
			}

			// 这里是同步映射表操作
			if msg_ctrl_info_v1.SubType == bgNetMessage.MsgControlInfoV1_ADD {
				// 增加映射表
				glog.Debug("MappingMgr::TunnelRecvCallback Control command : Add mapping table")
				// 首先根据映射地址查找映射对象
				//source_address := mapping_info_v1.SourceIp + ":" + strconv.Itoa(int(mapping_info_v1.SourcePort))
				_, exist := m.MappingTable[mapping_info_v1.MappingId]
				if !exist {
					glog.Debug("MappingMgr::TunnelRecvCallback Control command : Add mapping table >>> Mapping info is not exist")

					// 创建映射对象，初始化好，加入到映射表中
					mapping_object := MappingObject{
						Mapping_id : int(mapping_info_v1.MappingId),
						Mapping_ip : mapping_info_v1.MappingIp,
						Mapping_port : int(mapping_info_v1.MappingPort),
						Source_ip : mapping_info_v1.SourceIp,
						Source_port : int(mapping_info_v1.SourcePort),
						Net_type : mapping_info_v1.NetType,
						Is_running : int(mapping_info_v1.IsRunning),
					}

					// 初始化映射对象
					err = mapping_object.Initialize(m.SendResponse)

					if err != nil {
						glog.Debug("MappingMgr::TunnelRecvCallback Control command : Add mapping table >>> MappingObject::Initialize failed.")
						glog.Error(err)
					} else {
						glog.Debug("MappingMgr::TunnelRecvCallback Control command : Add mapping table >>> MappingObject::Initialize succeed.")
						glog.Debugf("MappingObject::Mapping_id : %d\n" +
							"MappingObject::Mapping_ip : %s\n" +
							"MappingObject::Mapping_port : %d\n" +
							"MappingObject::Source_ip : %s\n" +
							"MappingObject::Source_port : %d\n" +
							"MappingObject::Net_type : %s\n" +
							"MappingObject::Is_running : %d", mapping_object.Mapping_id, mapping_object.Mapping_ip,
							mapping_object.Mapping_port, mapping_object.Source_ip, mapping_object.Source_port,
							mapping_object.Net_type, mapping_object.Is_running)
					}

					// 将映射对象加入映射表
					m.MappingTable[mapping_info_v1.MappingId] = &mapping_object
				} else {
					// 映射表中存在，我们认为错误
					glog.Debug("MappingMgr::TunnelRecvCallback Control command : Add mapping table >>> Mapping info exist")
					err = errors.New("Mapping info exist")
					return err
				}

			} else if msg_ctrl_info_v1.SubType == bgNetMessage.MsgControlInfoV1_REMOVE {
				// 移除映射表
				glog.Debug("Remove mapping table")
			} else if msg_ctrl_info_v1.SubType == bgNetMessage.MsgControlInfoV1_QUERY {
				// 查询映射表
				glog.Debug("Query mapping table")
			} else if msg_ctrl_info_v1.SubType == bgNetMessage.MsgControlInfoV1_UPDATE {
				// 更新映射表
				glog.Debug("Update mapping table")
			} else if msg_ctrl_info_v1.SubType == bgNetMessage.MsgControlInfoV1_UPDATEALL {
				// 全量更新映射表
				glog.Debug("Update all mapping table")
			}
		}
	} else if msg_v1.MsgType == 1 {
		// 数据协议，属于直接透传的，找出映射ID，直接扔给对应的映射对象去处理
		//
		object, exist := m.MappingTable[msg_v1.MappingId]
		if !exist {
			// 映射未建立，直接返回错误
		} else {
			// 将消息扔到这个对象中处理
			err = object.RecvData(msg_v1)
		}
	}

	return err
}

func (m *MappingMgr) SendResponse(msg bgNetMessage.NetMessageV1) error {
	err := m.Net_tunnel_server.SendMessage(msg)
	return err
}