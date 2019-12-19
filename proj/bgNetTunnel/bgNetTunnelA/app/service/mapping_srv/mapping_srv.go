/*
映射服务
本服务用于管理映射入口，映射入口会根据传输层协议（TCP|UDP）提供对应的端口映射
端口映射表存在数据库中，数据库支持（MySQL|SQLite）
映射表应包含以下字段：
| 映射端IP(默认为0.0.0.0) | 映射端端口 | 源端IP | 源端端口 | 网络类型 | 映射状态(启用、禁用) |
本服务持有一个映射表，映射表的Key是源地址(源IP:源端口)
本服务持有一个隧道客户端(隧道A端)
本服务负责将数据传递给通道，以及从通道中取出数据
 */
package mapping_srv

import (
	"bgNetTunnelA/app/service/tunnel_cli"
	"bgNetTunnelA/tunnel_protocol"
	"crypto/md5"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/mitchellh/mapstructure"
)

type MappingService struct {
	TunnelProtocolObject tunnel_protocol.TunnelProto
	TunnelA tunnel_cli.TunnelClient
	MappingTable map[int]*MappingObject
	Enable_crypto bool
}

// 全局变量，映射服务
var Mapping_service_instance MappingService

/*
初始化映射服务
 */
func (m *MappingService) Initialize() error {

	//////////////////////////////////////////////////////////////////////////////////
	// 首先初始化隧道协议对象
	m.Enable_crypto = g.Config().GetBool("tunnel.enable_crypto")
	algorithm := g.Config().GetString("tunnel.crypto")
	factor := g.Config().GetString("tunnel.crypto_key_factor")

	// 生成Key，生成算法：md5(key + factor + 加密算法 + key)，结果截取前后8字节，得到16字节数据转换为byte
	md5 := md5.New()
	md5.Write([]byte("key" + factor + algorithm + "key"))
	key := md5.Sum(nil)

	// 生成IV，生成算法：md5(IV + factor + 加密算法 + IV)，结果截取前后8字节，得到16字节数据转换为byte
	md5.Reset()
	md5.Write([]byte("IV" + factor + algorithm + "IV"))
	iv := md5.Sum(nil)

	err := m.TunnelProtocolObject.Initialize(algorithm, key, iv)
	if err != nil {
		glog.Debug("Initialize tunnel protocol object failed.")
		return err
	} else {
		glog.Debug("Initialize tunnel protocol object succeed.")
	}

	//////////////////////////////////////////////////////////////////////////////////
	// 然后初始化隧道A端
	err = m.TunnelA.Initialize(m, key)
	if err != nil {
		glog.Debug("[MappingService::Initialize] Initialize tunnel A object failed.")
		return err
	} else {
		glog.Debug("[MappingService::Initialize] Initialize tunnel A object succeed.")
	}

	//////////////////////////////////////////////////////////////////////////////////
	// 从数据库读取映射表，创建映射对象，根据对应的映射状态启动映射对象
	m.MappingTable = make(map[int]*MappingObject, 16)
	result, err := g.DB("default").Table("bg_mapping_table").Select()
	if err != nil {
		glog.Debug("[MappingService::Initialize] Read mapping info from database failed.")
		glog.Error(err)
		return err
	} else {
		glog.Debug("[MappingService::Initialize] Read mapping info from database succeed.")
	}

	for _, record := range result {
		// 将取出的map转换成结构体
		record_map := record.Map()
		mapping_base_info := new(MappingBaseInfo)
		err = mapstructure.Decode(record_map, mapping_base_info)
		if err != nil {
			glog.Debug("[MappingService::Initialize] database map convert to struct failed.")
			glog.Error(err)
		} else {
			glog.Debug("[MappingService::Initialize] database map convert to struct succeed.")
		}

		// 初始化映射对象后，将映射对象加入映射表
		mapping_object := new(MappingObject)
		err = mapping_object.Initialize(mapping_base_info, m, &m.TunnelProtocolObject, m.Enable_crypto)
		if err != nil {
			glog.Debug("[MappingService::Initialize] Initialize mapping object failed.")
			glog.Error(err)
		} else {
			glog.Debug("[MappingService::Initialize] Initialize mapping object succeed.")
		}

		m.MappingTable[mapping_object.Info.Mapping_id] = mapping_object

		// 如果数据库标记为运行状态，则启动映射对象
		if mapping_object.Info.Is_running == 1 {
			mapping_object.Start()
		}
	}

	return err
}

/*
MappingService::PostDataToTunnel()
向通道发送数据
 */
func (m *MappingService) PostDataToTunnel(data []byte) error {
	err := m.TunnelA.PostDataToTunnel(data)
	if err != nil {
		glog.Debugf("[MappingService::PeekDataFromTunnel] Post data to tunnel-B failed.")
		glog.Error(err)
	}

	return err
}

/*
MappingService::PeekDataFromTunnel()
从通道读取数据，此函数应为通道观察者接口，这里会接收到隧道A端接收到的数据
*/
func (m *MappingService) PeekDataFromTunnel(data []byte) error {
	var err error

	// 这里反序列化，取出对应的映射key(即Address)，将消息分发到对应的映射管理器中
	tunnel_sec_protocol, err := m.TunnelProtocolObject.Unmarshal(data)
	if err != nil {
		glog.Debug("[MappingService::PeekDataFromTunnel] Unmarshal tunnel data failed.")
		glog.Error(err)
		return err
	}

	// 区分属于哪种类型的请求
	if tunnel_sec_protocol.Main == tunnel_protocol.MainType_Business {
		mapping_id := int(tunnel_sec_protocol.Data.MappingID)
		mapping_object, exists := m.MappingTable[mapping_id]
		if exists {
			glog.Debugf("[MappingService::PeekDataFromTunnel] Send bussiness data to client :\n%s", tunnel_sec_protocol.String())
			err = mapping_object.SendDataToClient(tunnel_sec_protocol.Data)
		} else {
			glog.Warningf("MappingService::PeekDataFromTunnel() not found mapping object with key : %d", mapping_id)
		}
	} else if tunnel_sec_protocol.Main == tunnel_protocol.MainType_Exception {
		glog.Debugf("[MappingService::PeekDataFromTunnel] Send exception data to client :\n%s", tunnel_sec_protocol.String())
		//mapping_object, exists := m.MappingTable[tunnel_sec_protocol.Data.DstSrvAddr]
		//if exists {
			// 这里应当有一个接口处理异常情况，
			//err = mapping_object.SendDataToClient(tunnel_sec_protocol.Data)
		//}
	}

	return err
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////
//
// 管理类接口
//
////////////////////////////////////////////////////////////////////////////////////////////////////////

func (m *MappingService) AddMappingInfoHandle(mapping_info *gjson.Json) error {
	mapping_port := mapping_info.GetInt("mapping_port")
	source_ip := mapping_info.GetString("source_ip")
	source_port := mapping_info.GetInt("source_port")
	net_type := mapping_info.GetString("net_type")

	// 首先添加到数据库
	result, err := g.DB("default").Insert("bg_mapping_table", gdb.Map{
		"Mapping_port" : mapping_port,
		"Source_ip" : source_ip,
		"Source_port" : source_port,
		"Net_type" : net_type,
		"Is_running" : 1,
	})

	if err != nil {
		glog.Error(err)
		return err
	}

	id, _ := result.LastInsertId()
	row, _ := result.LastInsertId()
	glog.Debugf("[MappingService::AddMappingInfoHandle] Insert mapping info into database succeed. Id : %d, Row : %d",
		id, row)

	// 然后创建映射对象、添加到内存映射表
	mapping_base_info := new(MappingBaseInfo)
	mapping_base_info.Mapping_id = int(id)
	mapping_base_info.Mapping_port = mapping_port
	mapping_base_info.Source_ip = source_ip
	mapping_base_info.Source_port = source_port
	mapping_base_info.Net_type = net_type
	mapping_base_info.Is_running = 1
	mapping_object := new(MappingObject)
	err = mapping_object.Initialize(mapping_base_info, m, &m.TunnelProtocolObject, m.Enable_crypto)
	if err != nil {
		glog.Debug("[MappingService::AddMappingInfoHandle] Initialize mapping object failed.")
		glog.Error(err)
	} else {
		glog.Debug("[MappingService::AddMappingInfoHandle] Initialize mapping object succeed.")
	}

	m.MappingTable[mapping_object.Info.Mapping_id] = mapping_object

	// 如果标记为运行状态，则启动映射对象
	if mapping_object.Info.Is_running == 1 {
		mapping_object.Start()
	}

	return nil
}