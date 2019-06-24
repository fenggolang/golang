package main

import (
	"fmt"
	"runtime"
	"time"
)



func main() {
	var a [10]int
	for i := 0; i < 10; i++ {
		go func(i int) {
			for {
				a[i]++
				runtime.Gosched() //由协程主动交出控制权,否则程序一直处于卡死状态
			}
		}(i)
	}
	time.Sleep(time.Millisecond)
	fmt.Println(a)
}

// output
[233 238 246 239 117 135 107 145 109 151]
