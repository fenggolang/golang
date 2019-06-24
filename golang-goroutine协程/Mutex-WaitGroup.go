package main

import (
	"fmt"
	"sync"
	"time"
)

func TestCounter() {
	counter := 0
	for i := 0; i < 5000; i++ {
		go func() {
			counter++
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Printf("counter = %d\n", counter)
}
func TestCounterMutexThreadSafe() {
	var mut sync.Mutex
	counter := 0
	for i := 0; i < 5000; i++ {
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			counter++
		}()
	}
	time.Sleep(1 * time.Second) // 这一步不严谨
	fmt.Printf("Mutex:counter = %d\n", counter)
}

func TestCounterWaitGroupThreadSafe() {
	var mut sync.Mutex
	var wg sync.WaitGroup
	counter := 0
	for i := 0; i < 5000; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				mut.Unlock()
			}()
			mut.Lock()
			counter++
			wg.Done()
		}()
	}
	wg.Wait() // 把time.Sleep(1 * time.Second)换成了wg.Wait()
	fmt.Printf("WaitGroup:counter = %d\n", counter)
}

// golang中通过通信来共享内存，而不是通过共享内存来通信
//Go 语言的并发模型准确的说是２种都支持：通过共享内存来通信和通过通信来共享内存
//通过共享内存来通信：比如sync.Mutex和sync.WaitGroup结合使用
//通过通信来共享内存就是协程之间通过channel来传递信息
/**
嗯, chan做同步的确是替代了一些mutex和wg的工作. 目前go里面chan做同步和其他一些并发结构可以相互转换, 比如semaphore即可以用mutex实现, 也可以用chan实现.
还有通过chan来传递数据, 接收端一个协程做数据处理, 把数据做到线程封闭, 或者解耦两端处理. 这种send, consume的pipeline的模式, 以前看zk的源码, java的用blocking queue做, 日志分处理多个阶段, 不断往下传递. 放大来看, 消息队列也是这样, 服务处理(tcp)也是一样.
目前mutex, wg还是有很多地方用的. 全局的cache用mutex. 并发调用时, add任务和wait用wg,  当然你用chan也可以做, 不过其实底层也是锁, 用起来感觉还别扭一些.  然后对于计数来说, atomic还是效率高一些.

https://draveness.me/golang-channel
*/
func main() {
	TestCounter()
	TestCounterMutexThreadSafe()
	TestCounterWaitGroupThreadSafe()
}
