package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// https://gist.github.com/legendtkl/1922db71553c849ef0029429f737aadb
// 用法：启动server.go，再启动client.go

var upgrader = websocket.Upgrader{} // use default options
var messageChan = make(chan string)

func updateMsg() {
	for {
		time.Sleep(5 * time.Second)
		messageChan <- time.Now().Format("2006-01-02 03:04:05 PM")
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	// http协议升级为websocket协议
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		c.WriteMessage(websocket.TextMessage, []byte(<-messageChan))
	}
}

func main() {
	go updateMsg()
	http.HandleFunc("/echo", echo)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
