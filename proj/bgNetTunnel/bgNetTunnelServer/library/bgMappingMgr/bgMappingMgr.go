/*
映射管理服务

此服务主要功能：

1、管理映射表（内存表与数据库表的增、删、查、改）
2、管理映射对象（创建、删除、启动、停止）
3、发送、接收消息
 */

package bgMappingMgr

import (
	"bgNetTunnelServer/library/bgNetProtocol/bgNetMessage"
	"bgNetTunnelServer/library/bgNetTunnel"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
)

/*
映射信息
 */
type MappingInfo struct {
	Mapping_id		int
	Mapping_ip		string
	Mapping_port	int
	Source_ip 		string
	Source_port 	int
	Net_type		string
	Is_running		int


}

type MappingMgr struct {
	// 网络隧道客户端
	Net_tunnel_client bgNetTunnel.NetTunnelClient
}

var (
	// 内存映射表
	MappingTable  		map[int]MappingInfo
	MappingObjectTable  map[int]*MappingObject

	// 数据表
	mysql_client 	= g.DB("default").Table("bg_mapping_table")
)

func (m *MappingMgr) Initialize(tunnel_srv_ip string, tunnel_srv_port int, tunnel_user string, tunnel_pass string, tunnel_send_dir string, tunnel_recv_dir string, tunnel_net_type string) error {
	// 连接到数据库，读取已经保存的映射数据
	MappingTable = make(map[int]MappingInfo, 0)
	MappingObjectTable = make(map[int]*MappingObject, 0)
	db_mapping_result, err := mysql_client.Select()
	if err == nil {

		// 获取数据，缓存到内存映射表中
		for _, record := range db_mapping_result {

			mapping_info := new(MappingInfo)
			record_map := record.Map()
			err = mapstructure.Decode(record_map, mapping_info)
			if err != nil {
				glog.Error(err)
			}

			// 2019-11-17 00:10
			// 这里的映射表标KEY设计的有点问题
			// 实际上应当直接使用PORT来作为KEY，因为一个映射网关不可能使用同一个PORT构建两个端口
			key := mapping_info.Mapping_port
			MappingTable[key] = *mapping_info
		}
	}

	// 启动隧道对象，创建隧道连接，同步映射表
	err = m.Net_tunnel_client.Initialize(tunnel_srv_ip, tunnel_srv_port, tunnel_user, tunnel_pass, tunnel_send_dir, tunnel_recv_dir, tunnel_net_type, m.TunnelRecvCallback)
	if err != nil {
		return err
	}

	// 检查哪些映射处于启动状态的，启动对应的映射端口服务
	// 同时同步映射表，原先是批量的，似乎有点问题，这里改成单个的吧
	for _, mapping_info := range MappingTable {

		// 组装映射信息，并序列化
		mapping_info_v1 := &bgNetMessage.MsgMappingInfoV1{
			MappingId : int32(mapping_info.Mapping_id),
			MappingIp : mapping_info.Mapping_ip,
			MappingPort: int32(mapping_info.Mapping_port),
			SourceIp : mapping_info.Source_ip,
			SourcePort : int32(mapping_info.Source_port),
			NetType : mapping_info.Net_type,
			IsRunning : int32(mapping_info.Is_running),
		}

		mapping_info_v1_bytes, err := proto.Marshal(mapping_info_v1)
		if err != nil {
			glog.Error(err)
			continue
		}
		mapping_info_v1_bytes_base64encode := gbase64.EncodeToString(mapping_info_v1_bytes)

		// 组装控制命令，并序列化
		ctl_cmd := &bgNetMessage.MsgControlInfoV1{
			MainType : bgNetMessage.MsgControlInfoV1_SyncMappingTable,
			SubType : bgNetMessage.MsgControlInfoV1_ADD,
			CtrlCommand : mapping_info_v1_bytes_base64encode,
		}

		ctl_cmd_bytes, err := proto.Marshal(ctl_cmd)
		if err != nil {
			glog.Error(err)
			continue
		}
		ctl_cmd_bytes_base64encode := gbase64.EncodeToString(ctl_cmd_bytes)

		// 组装隧道协议，并序列化
		msg_v1 := &bgNetMessage.NetMessageV1{
			MsgType : 0,
			MappingId : 0,
			ClientId : "control",
			MessageBody : ctl_cmd_bytes_base64encode,
		}

		// 发送映射表
		err = m.Net_tunnel_client.SendMessage(*msg_v1)
		if err != nil {
			glog.Error(err)
			continue
		}

		// 这里是测试内容，我们希望能反序列化回来
		ctl_cmd_bytes_base64decode_data, err := gbase64.DecodeString(msg_v1.MessageBody)
		ctl_cmd_unmashal := bgNetMessage.MsgControlInfoV1{}
		err = proto.Unmarshal(ctl_cmd_bytes_base64decode_data, &ctl_cmd_unmashal)

		mapping_info_v1_bytes_base64decode_data, err := gbase64.DecodeString(ctl_cmd_unmashal.CtrlCommand)
		mapping_info_unmashal := bgNetMessage.MsgMappingInfoV1{}
		err = proto.Unmarshal(mapping_info_v1_bytes_base64decode_data, &mapping_info_unmashal)

		println(mapping_info_unmashal.String())
	}

	// 最后根据映射表状态，启动对应的映射对象
	// 这里创建的映射对象也应当缓存起来，实际上找的应该是映射对象
	for _, mapping_info := range MappingTable {
		mapping_object_ptr := new(MappingObject)
		if mapping_info.Is_running == 1 {
			err = mapping_object_ptr.Startup(mapping_info.Mapping_id, mapping_info.Mapping_ip, mapping_info.Mapping_port, mapping_info.Net_type, &m.Net_tunnel_client)
		}
		MappingObjectTable[mapping_info.Mapping_id] = mapping_object_ptr
	}

	return nil
}

func (m *MappingMgr) TunnelRecvCallback(date []byte) error {
	// 首先协议栈解密，反序列化
	msg_v1 := bgNetMessage.NetMessageV1{}
	err := proto.Unmarshal(date, &msg_v1)
	if err != nil {
		glog.Error(err)
		return err
	}

	mapping_object, exist := MappingObjectTable[int(msg_v1.MappingId)]
	if !exist {
		glog.Error("Not Found Mapping Object")
		return err
	} else {
		response_data, err := gbase64.DecodeString(msg_v1.MessageBody)
		if err != nil {
			return err
		}

		err = mapping_object.SendMsgToClient(msg_v1.ClientId, response_data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MappingMgr) AddMapping(mapping_ip string, mapping_port int, source_ip string, source_port int, net_type string) error {
	// 首先尝试向数据库映射表添加映射信息
	
	// 成功后向内存映射表添加映射信息
	
	return nil
}

func (m *MappingMgr) RemmoveMapping(mapping_port int) error {
	// 首先检查映射对象处于停止状态，若不是，则返回失败

	// 然后从数据库映射表移除映射信息

	// 最后从内存映射表移除映射信息

	return nil
}