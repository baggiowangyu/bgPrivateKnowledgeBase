package main

import (
	_ "bgOrgMgr/boot"
	_ "bgOrgMgr/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	//wait_seconds := 5 * time.Second
	//time.Sleep(wait_seconds)
	g.Server().Run()
}
