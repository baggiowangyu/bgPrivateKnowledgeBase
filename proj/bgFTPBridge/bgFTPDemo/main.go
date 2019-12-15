package main

import "bytes"

func TestCallback(func a(data string) error) {
	for index := 0; ; index++ {

	}
}

func Callback(data string) error {
	println(data)

	return nil
}

func main() {
	//// 时间戳
	//for index := 0; ; index++ {
	//	println(time.Now().UnixNano())
	//	time.Sleep(100 * time.Microsecond)
	//}



	//clientA, err := ftp.Dial("127.0.0.1:21", ftp.DialWithTimeout(5 * time.Second))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = clientA.Login("userA", "123456")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = clientA.ChangeDir("send")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//entry_listA, err := clientA.List("/send")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//clientB, err := ftp.Dial("127.0.0.1:21", ftp.DialWithTimeout(5 * time.Second))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = clientB.Login("userB", "123456")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//err = clientB.ChangeDir("recv")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, entry_element := range entry_listA {
	//	// 读取文件
	//	println("name : ", entry_element.Name)
	//	println("size : ", entry_element.Size)
	//	response, err := clientA.Retr("/send/" + entry_element.Name)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	buf, err := ioutil.ReadAll(response)
	//	println(string(buf))
	//
	//	data := bytes.NewBuffer(buf)
	//	err = clientB.Stor("/recv/" + entry_element.Name, data)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//
	//	// 必须关闭Response才能正常读取下一个文件
	//	response.Close()
	//	data.Reset()
	//}
	//
	//_ = clientB.Quit()
	//_ = clientA.Quit()
}
