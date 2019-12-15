package bgMappingMgr

import "bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"

// 定义一个回调函数
type VitualClientRecvCallback func(msg bgNetMessage.NetMessageV1) error
