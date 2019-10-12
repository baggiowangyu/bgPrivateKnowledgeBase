package main

import (
	_ "./boot"
	_ "./router"
	"github.com/gogf/gf/frame/g"
)

func main() {
	//wait_seconds := 5 * time.Second
	//time.Sleep(wait_seconds)
	g.Server().Run()
}
