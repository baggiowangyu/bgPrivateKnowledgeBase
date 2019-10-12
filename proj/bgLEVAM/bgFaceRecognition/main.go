package main

import (
	_ "./boot"
	_ "./router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	g.Server().Run()
}
