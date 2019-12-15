/*

隧道协议定义：

1、隧道协议分为两种

- 控制信息(MsgType:0)，隧道两端控制数控同步所用的数据
- 数据信息(MsgType:1)，客户的交互信息

2、协议序列化规约

- 协议序列化采用protobuf v3序列化方案进行序列化与反序列化

3、

 */

package bgNetProtocol

type NetProtocol struct {

}



func init() {

}

/*
初始化隧道协议对象，这里主要是进行数据加密与解密
目前加密算法只支持AES或者DES，后面可以考虑SM4
 */
func (n *NetProtocol) Initialize(crypto_method string, key []byte, iv []byte) error {
	return nil
}

// 加密数据
func (n *NetProtocol) EncryptData(plain_data []byte) ([]byte, error) {
	return plain_data, nil
}

// 解密数据
func (n *NetProtocol) DecryptData(cipher []byte) ([]byte, error) {
	return cipher, nil
}
