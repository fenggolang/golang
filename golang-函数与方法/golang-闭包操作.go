package main

import (
	"fmt"
)

/**
闭包定义：内层函数引用了外层函数中的变量(或称为引用了自由变量的函数),
其返回值也是一个函数
*/

func main() {
	//f := outer(10) //返回的是一个函数
	//fmt.Println(f(100))

	testRange()
}

// golang中的闭包，(注意golang中不能嵌套函数，但是可以在一个函数中包含匿名函数)
func outer(x int) func(int) int {
	return func(y int) int {
		return x + y
	}
}

// 闭包用在哪些地方
// 1. for range中使用闭包
func testRange() {
	s := []string{"a", "b", "c"}
	for _, v := range s {
		go func(v string) {
			fmt.Println(v) // 由于使用了 go 协程，并非顺序输出。
		}(v) //每次将变量v的拷贝传进函数
	}

	select {} //阻塞模式:fatal error: all goroutines are asleep - deadlock!
}
