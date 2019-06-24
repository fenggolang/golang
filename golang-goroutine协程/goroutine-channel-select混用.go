package main

import (
	"fmt"
	"time"
)

// Golang并发实现：
/**
1. 程序并发执行(goroutine)
2. 多个goroutine间的数据同步和通信(channels)
3. select:多个channel选择数据读取或者写入(select)
*/

func goroutineTest() string {
	c := make(chan string) //创建一个channel
	go func() {
		time.Sleep(1 * time.Second)
		c <- "message from closure" // 发送数据到channel中
	}()
	msg := <-c

	// select(从多个channel中读取或写入数据)
	c1 := make(chan string)
	c2 := make(chan string)
	// 当c1有数据就从c1读取，c2有数据就从c2读取，c1,c2都有数据则随机读取
	select {
	case v := <-c1:
		fmt.Println("channel 1 sends", v)
	case v := <-c2:
		fmt.Println("channel 2 sends", v)
	default: //可选
		fmt.Println("neither channel was ready")
	}
	return msg
}
func main() {
	fmt.Println(goroutineTest())
}
