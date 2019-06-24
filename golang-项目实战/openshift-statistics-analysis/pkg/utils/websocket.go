package utils

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)



// 保证websocket超时不断开
func KeepAlive(c *websocket.Conn, timeout time.Duration) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Now().Sub(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

//func aa(cc net.Conn){
//
//}
//长连接入口
func HandleConnection(conn net.Conn, timeout int) {

	buffer := make([]byte, 2048)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			//LogErr(conn.RemoteAddr().String(), " connection error: ", err)
			fmt.Println(conn.RemoteAddr().String(), " connection error: ", err)
			return
		}
		Data := (buffer[:n])
		messnager := make(chan byte)
		postda := make(chan byte)
		fmt.Println(postda)
		//心跳计时
		go HeartBeating(conn, messnager, timeout)
		//检测每次Client是否有数据传来
		go GravelChannel(Data, messnager)
		Log("receive data length:", n)
		Log(conn.RemoteAddr().String(), "receive data string:", string(Data))

	}
}

//心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, readerChannel chan byte, timeout int) {
	select {
	case fk := <-readerChannel:
		Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		//conn.SetReadDeadline(time.Now().Add(time.Duration(5) * time.Second))
		break
	case <-time.After(time.Second * 5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}

}

func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}

func Log(v ...interface{}) {
	//log.Println(v...)
	fmt.Println(v...)
}

func sender(conn net.Conn) {
	for i := 0; i < 5; i++ {
		words := strconv.Itoa(i) + "This is a test for long conn"
		conn.Write([]byte(words))
		time.Sleep(2 * time.Second)

	}
	fmt.Println("send over")

}
