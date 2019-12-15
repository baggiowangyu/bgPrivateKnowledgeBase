package main

import (
	_ "bgRestfulServer/boot"
	_ "bgRestfulServer/router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
