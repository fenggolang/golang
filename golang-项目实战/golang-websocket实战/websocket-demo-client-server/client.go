package main

import (
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

// https://gist.github.com/legendtkl/1922db71553c849ef0029429f737aadb
// 用法：启动server.go，再启动client.go,clienty.go会每5秒输出：2018/05/14 10:01:32 recv: 2018-05-14 10:01:32 AM

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	time.Sleep(100 * time.Second)
}
