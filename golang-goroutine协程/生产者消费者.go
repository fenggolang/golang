package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	ch := make(chan int)
	// 一个生产者
	wg.Add(1)
	dataProducer(ch, &wg)

	// 多个消费者
	wg.Add(1)
	dataReceiver("消费者1", ch, &wg)
	wg.Add(1)
	dataReceiver("消费者2", ch, &wg)

	wg.Wait()
}

// 生产者
func dataProducer(ch chan int, wg *sync.WaitGroup) {
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			fmt.Printf("生产了数字:%d\n", i)
		}
		close(ch) // 写入十个数字之后关闭通道，这个类似于广播的形式，所有接受者会受到这个关闭信号之后不再接收
		wg.Done()
	}()
}

// 消费者
func dataReceiver(receiver string, ch chan int, wg *sync.WaitGroup) {
	go func(receiver string) {
		for {
			if data, ok := <-ch; ok { //ok 为true表示channel还没有关闭
				fmt.Printf("%s 消费了数字:%d\n", receiver, data)
			} else {
				break
			}
		}
		wg.Done()
	}(receiver)
}
