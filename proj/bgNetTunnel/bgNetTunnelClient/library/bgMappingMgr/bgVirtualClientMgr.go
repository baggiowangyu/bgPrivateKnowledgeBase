package bgMappingMgr

import (
	"bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"
	"errors"
	"github.com/gogf/gf/os/glog"
)

type VirtualClientMgr struct {
	Clients map[string]*VirtualClientObject
	RecvCallback 	VitualClientRecvCallback
}

func (v *VirtualClientMgr) Initialize(callback 	VitualClientRecvCallback) error {
	glog.Debug("VirtualClientMgr::Initialize")
	v.Clients = make(map[string]*VirtualClientObject, 0)
	v.RecvCallback = callback
	return nil
}

func (v *VirtualClientMgr) FindClient(mapping_id int, client_address string) (*VirtualClientObject, error) {
	client, exist := v.Clients[client_address]
	if exist {
		glog.Debugf("VirtualClientMgr::FindClient Find client(mapping_id : %d, client_address : %s)", mapping_id, client_address)
		return client, nil
	} else {
		glog.Debugf("VirtualClientMgr::FindClient Not find client(mapping_id : %d, client_address : %s)", mapping_id, client_address)
		return client, errors.New("Client not found")
	}
}

func (v *VirtualClientMgr) CreateClient(mapping_id int, client_address string, target_address string, net_type string) error {
	_, exist := v.Clients[client_address]
	if exist {
		glog.Debugf("VirtualClientMgr::CreateClient Client(mapping_id : %d, client_address : %s, target_address : %s, net_type : %s) already exist.", mapping_id, client_address, target_address, net_type)
		return errors.New("Client exist")
	} else {
		glog.Debugf("VirtualClientMgr::CreateClient Client(mapping_id : %d, client_address : %s, target_address : %s, net_type : %s) not exist.", mapping_id, client_address, target_address, net_type)
	}
	
	virtual_client := VirtualClientObject{
		Mapping_id		: mapping_id,
		Client_address 	: client_address,
		Target_address 	: target_address,
		Net_type		: net_type,
		RecvCallback	: v.RecvData,
	}

	v.Clients[client_address] = &virtual_client

	// 初始化
	err := virtual_client.Initialize()
	if err != nil {
		glog.Debug("VirtualClientMgr::CreateClient Initialize virtual client object failed.")
		glog.Error(err)
	} else {
		glog.Debug("VirtualClientMgr::CreateClient Initialize virtual client object succeed.")
	}

	return err
}

func (v *VirtualClientMgr) SendData(client_address string, data []byte) error {

	client, exist := v.Clients[client_address]
	if !exist {
		glog.Debug("VirtualClientMgr::SendData find client")
		return errors.New("Client not exist")
	}

	err := client.SendData(data)
	return err
}

func (v *VirtualClientMgr) RecvData(msg bgNetMessage.NetMessageV1) error {

	// 这是某个客户端返回上来的数据，继续往上层传吧，传到映射对象那边
	err := v.RecvCallback(msg)
	return err
}