package main

import (
	_ "bgEncryptTool/boot"
	_ "bgEncryptTool/router"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rc4"
	"errors"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/glog"
	"io/ioutil"
	"os"
)

func PKCS5Padding(data []byte, block_size int) []byte {
	// PKCS7填充规则：
	// 1、填充结果必须为块大小的整数倍，如果填充数据正好等于块大小整数倍，则后面附加一个附加块
	// 2、超出原始数据的部分，每个字节填入填充块长度
	padding := block_size - len(data) % block_size
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func PKCS5Unpadding(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length - 1])
	return data[:length-unpadding]
}

func AesGcmEncrypt(source_path, target_path, key, nonce string) error {
	if len(key) != 16 {
		err := errors.New("Key lenth error.")
		return err
	}

	if len(nonce) != 12 {
		err := errors.New("Key lenth error.")
		return err
	}

	key_bytes := []byte(key)
	nonce_bytes := []byte(nonce)

	block, err := aes.NewCipher([]byte(key_bytes))
	if err != nil {
		glog.Error(err)
		return err
	}

	// 读取文件
	plain, err := ioutil.ReadFile(source_path)
	if err != nil {
		glog.Error(err)
		return err
	}

	block_size := block.BlockSize()
	padding_plain_data := PKCS5Padding(plain, block_size)
	cipher_buffer := make([]byte, len(padding_plain_data))

	aesgcm, err := cipher.NewGCM(block)
	if err == nil {
		cipher_buffer = aesgcm.Seal(nil, nonce_bytes, padding_plain_data, nil)
		err = ioutil.WriteFile(target_path, cipher_buffer, os.ModeExclusive)
	}

	return err
}

func AesGcmDecrypt(source_path, target_path, key, nonce string) error {
	if len(key) != 16 {
		err := errors.New("Key lenth error.")
		return err
	}

	if len(nonce) != 12 {
		err := errors.New("Key lenth error.")
		return err
	}

	key_bytes := []byte(key)
	nonce_bytes := []byte(nonce)

	block, err := aes.NewCipher([]byte(key_bytes))
	if err != nil {
		glog.Error(err)
		return err
	}

	// 读取文件
	cipher_data, err := ioutil.ReadFile(source_path)
	if err != nil {
		glog.Error(err)
		return err
	}

	cipher_len := len(cipher_data)
	padding_plain_buffer := make([]byte, cipher_len)

	aesgcm, err := cipher.NewGCM(block)
	if err == nil {
		padding_plain_buffer, err = aesgcm.Open(nil, nonce_bytes, cipher_data, nil)
		err = ioutil.WriteFile(target_path, padding_plain_buffer, os.ModeExclusive)
	}

	return err
}

func Rc4Crypt(source_path, target_path, key string) error {

	cipher, err := rc4.NewCipher([]byte(key))
	if err != nil {
		glog.Error(err)
		return err
	}

	// 读取文件
	plain, err := ioutil.ReadFile(source_path)
	if err != nil {
		glog.Error(err)
		return err
	}

	dst := make([]byte, len(plain))
	cipher.XORKeyStream(dst, plain)

	err = ioutil.WriteFile(target_path, dst, os.ModeExclusive)
	return err
}

func main() {
	//g.Server().Run()

	// 根据参数对文件进行加密
	//argc := len(gcmd.GetArgAll())
	//if argc < 5 {
	//	fmt.Print("bgEncryptpTool --s=[source_file] --t=[target_file] --m=[Encrypt|Decrypt] --c=[AES_GCM|RC4] --k=[KEY] --n=[NONCE]")
	//	fmt.Print("   KEY >>> 16bytes")
	//	fmt.Print("   NONCE >>> 12bytes")
	//	return
	//}

	source_file := gcmd.GetOpt("s")
	target_file := gcmd.GetOpt("t")
	method := gcmd.GetOpt("m")
	crypto := gcmd.GetOpt("c")
	key := gcmd.GetOpt("k")
	nonce := gcmd.GetOpt("n")

	var err error
	if method == "Encrypt" {
		if crypto == "AES_GCM" {
			// AES-GCM加密
			err = AesGcmEncrypt(source_file, target_file, key, nonce)
		} else if crypto == "RC4" {
			// RC4加密
			err = Rc4Crypt(source_file, target_file, key)
		} else {
			err = errors.New("Not support crypto")
		}
	} else if method == "Decrypt" {
		if crypto == "AES_GCM" {
			// AES-GCM解密
			err = AesGcmDecrypt(source_file, target_file, key, nonce)
		} else if crypto == "RC4" {
			// RC4解密
			err = Rc4Crypt(source_file, target_file, key)
		} else {
			err = errors.New("Not support crypto")
		}
	} else {
		err = errors.New("Not support method")
	}

	if err != nil {
		glog.Error(err)
	}
}
