package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//isSymbol表示有无符号
func BytesToInt(b []byte, isSymbol bool)  (int, error){
	if isSymbol {
		return bytesToIntS(b)
	}
	return bytesToIntU(b)
}


//字节数(大端)组转成int(无符号的)
func bytesToIntU(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0},b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp uint8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp uint16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp uint32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0,fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}



//字节数(大端)组转成int(有符号)
func bytesToIntS(b []byte) (int, error) {
	if len(b) == 3 {
		b = append([]byte{0},b...)
	}
	bytesBuffer := bytes.NewBuffer(b)
	switch len(b) {
	case 1:
		var tmp int8
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 2:
		var tmp int16
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	case 4:
		var tmp int32
		err := binary.Read(bytesBuffer, binary.BigEndian, &tmp)
		return int(tmp), err
	default:
		return 0,fmt.Errorf("%s", "BytesToInt bytes lenth is invaild!")
	}
}


//整形转换成字节
func IntToBytes(n int,b byte) ([]byte,error) {
	switch b {
	case 1:
		tmp := int8(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 2:
		tmp := int16(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	case 3,4:
		tmp := int32(n)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, &tmp)
		return bytesBuffer.Bytes(),nil
	}
	return nil,fmt.Errorf("IntToBytes b param is invaild")
}

func main() {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)
	if listener, err := kcp.ListenWithOptions("127.0.0.1:12345", block, 10, 3); err == nil {
		// spin-up the client
		go client()
		for {
			s, err := listener.AcceptKCP()
			if err != nil {
				log.Fatal(err)
			}
			go handleEcho(s)
		}
	} else {
		log.Fatal(err)
	}
}

// handleEcho send back everything it received
func handleEcho(conn *kcp.UDPSession) {
	buf := make([]byte, 4096)
	for {
		read_byte, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("Read length : %d\n", read_byte)

		read_buf_data_len, err := BytesToInt(buf, true)
		data_buf := make([]byte, read_buf_data_len)

		read_byte, err = conn.Read(data_buf)
		if err != nil {
			log.Println(err)
			return
		}

		err = ioutil.WriteFile("E:/opensource/bgPrivateKnowledgeBase/proj/bgNetTunnel/TestKCP/1.exe", data_buf, os.ModeExclusive)
		if err != nil {
			log.Println(err)
			return
		}

		//n, err = conn.Write(buf[:read_byte])
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
	}
}

func client() {
	key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
	block, _ := kcp.NewAESBlockCrypt(key)

	// wait for server to become ready
	time.Sleep(time.Second)

	// 读取一个比较大的数据，最好超过64MB，然后一次性扔给服务端，服务端存下来
	data_buf, err := ioutil.ReadFile("E:/opensource/bgPrivateKnowledgeBase/proj/bgNetTunnel/TestKCP/nosqlbooster4mongo-5.1.5.exe")
	if err != nil {
		println(err)
		return
	}

	data_buf_len := len(data_buf)

	data_buf_len_bytes, err := IntToBytes(data_buf_len, 4)

	// dial to the echo server
	if sess, err := kcp.DialWithOptions("127.0.0.1:12345", block, 10, 3); err == nil {
		send_bytes, err := sess.Write(data_buf_len_bytes)
		if err == nil {
			// 发送长度完成，接下来发送实际数据
			log.Printf("Send length : %d\n", send_bytes)
			send_bytes, err = sess.Write(data_buf)
			if err == nil {
				log.Println("Send data length : %d", send_bytes)
			} else {
				log.Fatal(err)
			}
		} else {
			log.Fatal(err)
		}
		time.Sleep(time.Second)
	} else {
		log.Fatal(err)
	}
}