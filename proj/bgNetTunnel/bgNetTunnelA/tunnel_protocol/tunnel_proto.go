/*
隧道协议

本隧道协议基于两个结构体来实现：
1、协议头，包含了路由信息
2、协议体，包含了实际需要透传的用户数据

本隧道协议用于组装、加密透传信息
本隧道协议用于解密、拆装加密的透传信息
 */
package tunnel_protocol

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"github.com/gogf/gf/os/glog"
	"github.com/golang/protobuf/proto"
	"github.com/tjfoc/gmsm/sm4"
)

const (
	AES_ECB_128 = "AES-ECB-128"
	AES_CBC_128 = "AES-CBC-128"
	AES_GCM_128 = "AES-GCM-128"
	SM4_128		= "SM4-128"
	MAGIC		= "bg"
)

type TunnelProto struct {
	crypt_algorithm string
	cipher_block 	cipher.Block
	// 在AES-GCM下充当为随机数
	iv 				[]byte
}

/*
初始化隧道协议

参数：
@algorithm	[入参] 加密算法名称("NONE" | "AES-ECB-128" | "AES-ECB-256" ...)
@key		[入参] 加密密钥
@iv			[入参] 密钥向量

返回值：
@error		[] 错误信息
 */
func (t *TunnelProto) Initialize(algorithm string, key []byte, iv []byte) error {
	key_len := len(key)

	if algorithm == AES_ECB_128 {
		// AES_ECB_128
		if key_len != 16 {
			err := errors.New("Key lenth error.")
			return err
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		t.cipher_block = block

	} else if algorithm == AES_CBC_128 {
		// AES_CBC_128
		if key_len != 16 {
			err := errors.New("Key lenth error.")
			return err
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		t.cipher_block = block
		t.iv = iv

	} else if algorithm == AES_GCM_128 {
		// AES-GCM-128
		if key_len != 16 {
			err := errors.New("Key lenth error.")
			return err
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			return err
		}

		t.cipher_block = block
		t.iv = iv
	} else if algorithm == SM4_128 {
		// SM4_128
		if key_len != 16 {
			err := errors.New("Key lenth error.")
			return err
		}

		block, err := sm4.NewCipher(key)
		if err != nil {
			return err
		}

		t.cipher_block = block
		t.iv = iv
	} else {
		// 不支持
	}

	t.crypt_algorithm = algorithm

	return nil
}

/*
数据加密

参数：
@plain		[入参] 明文数据

返回值：
@[]byte		[] 加密结果
@error		[] 错误信息
 */
func (t *TunnelProto) EncryptData(plain []byte) ([]byte, error) {
	var err error = nil

	plain_len := len(plain)
	block_size := t.cipher_block.BlockSize()
	padding_plain_data := t.PKCS5Padding(plain, block_size)
	cipher_buffer := make([]byte, len(padding_plain_data))

	if t.crypt_algorithm == AES_ECB_128{
		// AES_ECB_128
		for bs, be := 0, block_size; bs < plain_len; bs, be = bs + block_size, be + block_size {
			t.cipher_block.Encrypt(cipher_buffer[bs:be], padding_plain_data[bs:be])
		}

		err = nil
	} else if t.crypt_algorithm == AES_CBC_128{
		// AES_CBC_128 || AES_CBC_256
		block_mode := cipher.NewCBCEncrypter(t.cipher_block, t.iv)
		block_mode.CryptBlocks(cipher_buffer, padding_plain_data)

		err = nil
	} else if t.crypt_algorithm == AES_GCM_128 {
		// AES_GCM_128
		aesgcm, err := cipher.NewGCM(t.cipher_block)
		if err == nil {
			cipher_buffer = aesgcm.Seal(nil, t.iv, padding_plain_data, nil)
		}
	} else if t.crypt_algorithm == SM4_128 {
		// SM4_128
		block_mode := cipher.NewCBCEncrypter(t.cipher_block, t.iv)
		block_mode.CryptBlocks(cipher_buffer, padding_plain_data)

		err = nil
	} else {
		// 不支持
		err = errors.New("Not support")
	}

	return cipher_buffer, err
}

/*
数据解密

参数：
@cipher_data [入参] 密文数据

返回值：
@[]byte		 [] 解密结果
@error		 [] 错误信息
*/
func (t *TunnelProto) DecryptData(cipher_data []byte) ([]byte, error) {
	var err error = nil

	cipher_len := len(cipher_data)
	block_size := t.cipher_block.BlockSize()
	padding_plain_buffer := make([]byte, cipher_len)

	if t.crypt_algorithm == AES_ECB_128 {
		// AES_ECB_128 || AES_ECB_256
		for bs, be := 0, block_size; bs < cipher_len; bs, be = bs + block_size, be + block_size {
			t.cipher_block.Decrypt(padding_plain_buffer[bs:be], cipher_data[bs:be])
		}

		err = nil
	} else if t.crypt_algorithm == AES_CBC_128 {
		// AES_CBC_128 || AES_CBC_256
		block_mode := cipher.NewCBCDecrypter(t.cipher_block, t.iv)
		block_mode.CryptBlocks(padding_plain_buffer, cipher_data)

		err = nil
	} else if t.crypt_algorithm == AES_GCM_128 {
		// AES_GCM_128
		aesgcm, err := cipher.NewGCM(t.cipher_block)
		if err == nil {
			padding_plain_buffer, err = aesgcm.Open(nil, t.iv, cipher_data, nil)
			if err != nil {
				glog.Error(err)
			}
		}
	} else if t.crypt_algorithm == SM4_128 {
		// SM4_128
		block_mode := cipher.NewCBCDecrypter(t.cipher_block, t.iv)
		block_mode.CryptBlocks(padding_plain_buffer, cipher_data)
	} else {
		// 不支持
		err = errors.New("Not support")
	}

	var plain_data []byte
	if err == nil {
		plain_data = t.PKCS5Unpadding(padding_plain_buffer)
	}

	return plain_data, err
}

/*
PKCS7数据填充算法

参数：
@data		[入参] 原始数据


返回值：
@[]byte		[] 填充数据
*/
func (t *TunnelProto) PKCS5Padding(data []byte, block_size int) []byte {
	// PKCS7填充规则：
	// 1、填充结果必须为块大小的整数倍，如果填充数据正好等于块大小整数倍，则后面附加一个附加块
	// 2、超出原始数据的部分，每个字节填入填充块长度
	padding := block_size - len(data) % block_size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

/*
PKCS7数据反填充算法

参数：
@data		[入参] 填充数据


返回值：
@[]byte		[] 原始数据
*/
func (t *TunnelProto) PKCS5Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length - 1])
	return data[:length-unpadding]
}

/*
序列化隧道协议数据体

参数：
@data		[入参] 隧道数据体
@sec		[入参] 是否加密

返回值：
 */
func (t *TunnelProto) Marshal(version string, main_type MainType, sub_type SubType, sec bool, data *TunnelProtocol) ([]byte, error) {
	var tunnel_sec_protocol TunnelSecProtocol
	if sec {
		// 将实际数据序列化后加密
		tunnel_protocol_data, err := proto.Marshal(data)
		if err != nil {
			return nil, err
		}

		tunnel_protocol_data_cipher, err := t.EncryptData(tunnel_protocol_data)
		if err != nil {
			return nil, err
		}

		// 构建完整协议体
		tunnel_sec_protocol = TunnelSecProtocol{
			Magic : MAGIC,
			Version : version,
			Main : main_type,
			Sub : sub_type,
			IsSec : sec,
			SecData : tunnel_protocol_data_cipher,
		}
	} else {
		// 直接将协议数据放到加密隧道包内，然后序列化
		tunnel_sec_protocol = TunnelSecProtocol{
			Magic : MAGIC,
			Version : version,
			Main : main_type,
			Sub : sub_type,
			IsSec : sec,
			Data : data,
		}
	}

	tunnel_sec_protocol_data, err := proto.Marshal(&tunnel_sec_protocol)
	return tunnel_sec_protocol_data, err
}

/*
反序列化隧道协议数据体

参数：
@data		[入参] 隧道协议数据

返回值：
*/
func (t *TunnelProto) Unmarshal(data []byte) (*TunnelSecProtocol, error) {
	tunnel_sec_protocol := new(TunnelSecProtocol)
	err := proto.Unmarshal(data, tunnel_sec_protocol)
	if err != nil {
		glog.Error(err)
		return nil, err
	}

	if tunnel_sec_protocol.IsSec {
		// 是加密的，那么我们就解密它
		tunnel_protocol_marshal_data, err := t.DecryptData(tunnel_sec_protocol.SecData)
		if err != nil {
			glog.Error(err)
			return nil, err
		}

		tunnel_protocol := new(TunnelProtocol)
		err = proto.Unmarshal(tunnel_protocol_marshal_data, tunnel_protocol)
		if err != nil {
			glog.Error(err)
			return nil, err
		}

		tunnel_sec_protocol.IsSec = false
		tunnel_sec_protocol.Data = tunnel_protocol
		_ = tunnel_sec_protocol.SecData
	}

	return tunnel_sec_protocol, err
}

/**
组装数据
*/
func (t *TunnelProto) BuildDataV1_0(data []byte, mid int32, src_cli_addr string, dst_srv_addr string, nettype NetType, mtype MainType, stype SubType, sec bool) ([]byte, error) {
	tunnel_data := TunnelProtocol{
		MappingID : mid,
		SrcCliAddr : src_cli_addr,
		DstSrvAddr : dst_srv_addr,
		Net : nettype,
		Data : data,
	}

	tunnel_sec_data, err := t.Marshal("1.0", mtype, stype, sec, &tunnel_data)
	if err != nil {
		glog.Error(err)
	}

	return tunnel_sec_data, err
}