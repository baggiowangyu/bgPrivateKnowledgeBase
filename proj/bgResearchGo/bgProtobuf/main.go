package main

import (
	"fmt"
	"io/ioutil"
	"os"
)
import pb "bgProtobuf/pb"

func main() {
	fmt.Println("Protobuf demo ...")

	person := pb.Person{}
	person.Rid = "000001"
	person.Name = "阿狗"
	person.Orgid = "100000"

	// 序列化
	var b []byte
	result, err := person.XXX_Marshal(b, false)
	if err != nil {
		println(err)
		return
	}

	// 写文件
	err = ioutil.WriteFile("personinfo.dat", result, os.ModeAppend)

	fmt.Println(result)

	// 反序列化
	person_new := pb.Person{}
	err = person_new.XXX_Unmarshal(result)
	if err != nil {
		println(err)
		return
	}

	println(person_new)
}
