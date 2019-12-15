package bgMappingMgr

import (
	"bgNetTunnelClient/library/bgNetProtocol/bgNetMessage"
	"errors"
	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/net/gudp"
	"github.com/gogf/gf/os/glog"
)

type VirtualClientObject struct {
	Mapping_id		int
	Client_address 	string
	Target_address	string
	Net_type		string
	Tcp_conn		*gtcp.Conn
	Udp_conn		*gudp.Conn
	Is_running		bool
	RecvCallback 	VitualClientRecvCallback
}

func (v *VirtualClientObject) Initialize() error {
	var err error
	// 创建连接
	if v.Net_type == "TCP" {
		v.Tcp_conn, err = gtcp.NewConn(v.Target_address)
	} else if v.Net_type == "UDP" {
		v.Udp_conn, err = gudp.NewConn(v.Target_address)
	} else {
		err = errors.New("Not supported net type")
	}

	if err != nil {
		glog.Debugf("VirtualClientObject::Initialize create new connection failed. Net type is %s", v.Net_type)
		glog.Error(err)
	} else {
		glog.Debugf("VirtualClientObject::Initialize create new connection succeed. Net type is %s", v.Net_type)
	}

	// 创建接收线程，接收到数据后进行回调
	v.StartRecv(v.RecvCallback)

	return err
}

func (v *VirtualClientObject) SendData(data []byte) error {
	var err error

	SEND:
	if v.Net_type == "TCP" {
		err = v.Tcp_conn.Send(data)
	} else if v.Net_type == "UDP" {
		err = v.Udp_conn.Send(data)
	} else {
		err = errors.New("Not supported net type")
	}

	if err != nil {
		if err.Error() == "EOF" {
			// 连接已经中断了
		} else {
			glog.Debugf("VirtualClientObject::SendData send data failed. Net type is %s", v.Net_type)
			// 尝试重连
			if v.Net_type == "TCP" {
				v.Tcp_conn, err = gtcp.NewConn(v.Target_address)
				goto SEND
			} else if v.Net_type == "UDP" {
				v.Udp_conn, err = gudp.NewConn(v.Target_address)
				goto SEND
			} else {
				err = errors.New("Not supported net type")
			}
		}

		glog.Error(err)
	} else {
		glog.Debugf("VirtualClientObject::SendData send data succeed. Net type is %s", v.Net_type)
	}

	return err
}

func (v *VirtualClientObject) StartRecv(callback VitualClientRecvCallback) {
	v.RecvCallback = callback
	go v.RecvThread()
}

func (v *VirtualClientObject) RecvThread() {
	var err error
	var data []byte

	glog.Debug("VirtualClientObject::RecvThread Started.")

	for {

		// 接收数据
		if v.Net_type == "TCP" {
			data, err = v.Tcp_conn.Recv(-1)
		} else if v.Net_type == "UDP" {
			data, err = v.Udp_conn.Recv(-1)
		} else {
			err = errors.New("Not supported net type")
		}

		// 出错的话跳出
		if err != nil {
			glog.Debugf("VirtualClientObject::RecvThread Recv data failed. Net type is %s", v.Net_type)
			break
		} else {
			glog.Debugf("VirtualClientObject::RecvThread Recv data succeed. Net type is %s", v.Net_type)
		}

		// 数据先Base64编码
		data_base64encode := gbase64.EncodeToString(data)

		// 这里可能要打包一下
		msg_v1 := bgNetMessage.NetMessageV1{
			MsgType: 1,
			MappingId: int32(v.Mapping_id),
			ClientId: v.Client_address,
			MessageBody: data_base64encode,
		}

		err = v.RecvCallback(msg_v1)
		if err != nil {
			glog.Debug("VirtualClientObject::RecvThread Call recv callback failed.")
			glog.Error(err)
		} else {
			glog.Debug("VirtualClientObject::RecvThread Call recv callback succeed.")
		}
	}

	glog.Error(err)
	return
}
