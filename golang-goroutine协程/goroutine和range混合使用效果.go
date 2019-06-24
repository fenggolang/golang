package main

import (
	"fmt"
	"time"
)

func test1() {
	a := []int{1, 2, 3}
	// 输出都是3是因为：在for range循环里面i变量是复用的，所以每次循环都会
	// 修改i的值，而goroutine执行的速度比for 循环慢，所以最终打印输出都是3！！！
	for _, i := range a {
		go func() {
			fmt.Println(i) // i都是a切片的最后一个值
		}()
	}

	// 下面一行代码不能保证运行结果正确
	time.Sleep(time.Second)
}

func test2() {
	a := []int{1, 2, 3}
	for _, i := range a {
		b := i
		go func() {
			fmt.Println(b) // 这样会随机打印1,2，3，顺序不一定
		}()

		// 下面在for循环中添加休眠，，让for循环暂停一下,goroutine
		// 可以=在每个for循环内有足够的时间执行，就可以得到这样的顺序：1,2,3
		//time.Sleep(time.Millisecond * 100)
	}

	// 下面一行代码不能保证运行结果正确
	time.Sleep(time.Second)
}

func test3() {
	a := []int{1, 2, 3}
	for _, i := range a {
		go func(b int) {
			fmt.Println(b)
		}(i)
	}

	// 下面一行代码不能保证运行结果正确
	time.Sleep(time.Second)
}
func main() {
	test1() // 输出3行3， goroutine内部实现决定的，可以看下Go语言读书笔记

	test2() // 随机打印输出1,2,3，顺序不定

	test3()
}
