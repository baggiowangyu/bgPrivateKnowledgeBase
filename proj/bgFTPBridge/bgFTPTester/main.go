package main

import (
	"github.com/jlaffaye/ftp"
	"log"
	"time"
)

func main() {

	client, err := ftp.Dial("127.0.0.1:21", ftp.DialWithTimeout(5 * time.Second))
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("connect succeed.")
	}

	err = client.Login("userA", "123456")
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("login succeed.")
	}

	for {
		entries, err := client.List("/send")
		if err != nil {
			log.Println(err)
			break
		} else {
			log.Println("list succeed.")
		}

		for index, entry := range entries {
			log.Println("%d %s", index, entry.Name)
		}

		time.Sleep(1 * time.Second)
	}
}
