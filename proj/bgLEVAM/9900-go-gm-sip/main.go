package main

import (
	"fmt"
	"github.com/gogf/gf/net/gtcp"
)

func main() {
	fmt.Println("9900-go-gm-sip start up ...")

	_ = gtcp.NewServer("0.0.0.0:9900", func(conn *gtcp.Conn) {
		defer conn.Close()

		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					fmt.Println(err)
				}
			}

			if err != nil {
				break
			}
		}
	}).Run()
}
