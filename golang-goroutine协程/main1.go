package main

import (
	"fmt"
	"time"
)



func main() {
	// 开10个人去打印printf
	for i := 0; i < 10; i++ {
		// 不加go关键字将会一直打印:Hello from goroutine 0
		// 加了go关键字之后表明是并发去执行这个函数,但是主程序函数还是在往下跑，如果goroutine还没来得及打印东西，main函数
		// 就退出的话，就不会打印东西
		go func(i int) {
			for {
				fmt.Printf("Hello from "+"goroutine %d\n", i)
			}
		}(i)
	}

	//让主程序main休息1毫秒，使得上面可以打印输出
	time.Sleep(time.Millisecond)
}

// 结论：协程goroutine,是轻量级的线程
// 协程是非抢占式多任务处理，由协程主动交出控制权
// 线程是没有控制权，操作系统回收线程也只能释放资源

// output:
...
Hello from goroutine 4
Hello from goroutine 4
Hello from goroutine 9
Hello from goroutine 9
Hello from goroutine 6
Hello from goroutine 9
Hello from goroutine 9
Hello from goroutine 9