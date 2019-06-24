package main

import (
	"fmt"
	"runtime"
	"time"
)



func main() {
	var a [3]int
	for i := 0; i < 10; i++ {
		// 如果下面的匿名函数不传递参数i进去则报错：panic: runtime error: index out of range
		go func() {
			a[i]++
			runtime.Gosched()
		}() // 如果下面的匿名函数不传递参数i进去则报错：panic: runtime error: index out of range
	}
	time.Sleep(time.Millisecond)
	fmt.Println(a)
}

// output:
// panic: runtime error: index out of range
