package main

import (
	_ "./boot"
	_ "./router"
	"fmt"
	"github.com/gogf/gf/frame/g"
)

func main() {
	fmt.Println("bgDeviceMgr started ...")
	g.Server().Run()
}