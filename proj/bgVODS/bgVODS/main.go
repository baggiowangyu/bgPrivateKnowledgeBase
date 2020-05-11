package main

import (
	_ "bgVODS/boot"
	_ "bgVODS/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	//println(time.Now().Format())
	g.Server().Run()
}
