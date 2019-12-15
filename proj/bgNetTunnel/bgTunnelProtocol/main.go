package main

import (
	_ "bgTunnelProtocol/boot"
	_ "bgTunnelProtocol/router"
	"bgTunnelProtocol/tunnel_protocol"
	"bytes"
	"github.com/gogf/gf/os/glog"
)

var origin_data = []byte("This is a test data")

func test_crypto(a string,  key []byte, iv []byte) {
	tunnel_proto := new(tunnel_protocol.TunnelProto)
	err := tunnel_proto.Initialize(a, key, iv)
	if err != nil {
		glog.Error(err)
		return
	}

	cipher_data, err := tunnel_proto.EncryptData(origin_data)
	if err != nil {
		glog.Error(err)
		return
	}

	plain_data, err := tunnel_proto.DecryptData(cipher_data)
	if err != nil {
		glog.Error(err)
		return
	}

	if bytes.Equal(origin_data, plain_data) {
		glog.Infof("%s OK", a)
	}
}

func test_proto(a string,  key []byte, iv []byte) {
	// 初始化隧道协议对象
	tunnel_proto := new(tunnel_protocol.TunnelProto)
	err := tunnel_proto.Initialize(a, key, iv)
	if err != nil {
		glog.Error(err)
		return
	}

	// 创建隧道数据
	data := new(tunnel_protocol.TunnelProtocol)
	data.SrcCliAddr = "127.0.0.1:43210"
	data.DstSrvAddr = "127.0.0.1:80"
	data.Net = tunnel_protocol.NetType_TCP
	data.MappingID = 1
	data.Data = []byte("This is a test data ...")

	// 序列化隧道协议
	marshal_data, err := tunnel_proto.Marshal("1.0", tunnel_protocol.MainType_Business, tunnel_protocol.SubType_Request, true, data)
	if err != nil {
		glog.Error(err)
		return
	}

	// 反序列化隧道协议
	tunnel_sec_protocol, err := tunnel_proto.Unmarshal(marshal_data)
	if err != nil {
		glog.Error(err)
		return
	}

	glog.Info(tunnel_sec_protocol.String())
}

func main() {
	//g.Server().Run()

	key_128 := []byte("1234567890123456")
	iv_128 := []byte("1234567890123456")
	nonce := []byte("123456789012")

	test_crypto(tunnel_protocol.AES_ECB_128, key_128, iv_128)
	test_crypto(tunnel_protocol.AES_CBC_128, key_128, iv_128)
	test_crypto(tunnel_protocol.AES_GCM_128, key_128, nonce)
	test_crypto(tunnel_protocol.SM4_128, key_128, iv_128)

	test_proto(tunnel_protocol.AES_ECB_128, key_128, iv_128)
	test_proto(tunnel_protocol.AES_CBC_128, key_128, iv_128)
	test_proto(tunnel_protocol.AES_GCM_128, key_128, nonce)
	test_proto(tunnel_protocol.SM4_128, key_128, iv_128)
}
