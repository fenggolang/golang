package main

import (
	"fmt"
	"time"
)

/**
关于管道：Channel
Channels用来同步并发执行的函数并提供它们某种传值交流的机制
Channels的一些特性：通过channel传递的元素类型、容器(或缓冲区)和传递的方向由"<-"操作符指定
*/

// 演示无缓冲和有缓冲的channel的样子
func test0() {
	done := make(chan bool)     //无缓冲
	done1 := make(chan bool, 1) //有缓冲
	fmt.Println(done, done1)
}

// 演示无缓冲在同一个main里面的死锁的例子
func test1() {
	// 编译错误deadlock(死锁),阻死main进程
	done := make(chan bool)
	done <- true //这句是输入值，它会一直阻塞，等待读取,因为创建chan的时候没有指定缓冲区，所以输入的值的同时必须有读取这，这可以通过2个goroutine协程来实现同时输入和读取
	<-done       // 这句是读取，但是在上面已经阻死了，永远走不到这里
	fmt.Println("完成")
}

// 演示仅有输入语句，但没读取语句的死锁例子
func test2() {
	done := make(chan bool)
	done <- true
	fmt.Println("完成")
}

// 演示仅有读取语句，但没有输入语句的死锁例子
func test3() {
	done := make(chan bool)
	<-done //读取输出，但前面没有输入语句，done是empty的，所以一直等待输入
	fmt.Println("完成")
}

// 演示：协程的阻死，不会影响main
func test4() {
	/** 编译通过 */
	// 演示：协程的阻死，不会影响main
	done := make(chan bool)
	go func() {
		<-done //一直等待
	}()
	fmt.Println("完成")
	// 控制台输出
}

// 在test4的基础上，无缓冲channel在协程go routine里面阻塞死
func test5() {
	done := make(chan bool)
	go func() {
		fmt.Println("我可能会输出") // 阻塞前的语句
		done <- true          //这里阻塞死，因为没有缓冲区。则chan必须同时输入和读取
		fmt.Println("我永远不会输出")
		<-done //这句也不会走到，除非在别的协程里面读取，或者在main
	}()
	fmt.Println("完成")
}

// 在test5的基础上，延时main的跑完
func test6() {
	done := make(chan bool)
	go func() {
		fmt.Println("我可能会输出")
		done <- true
		fmt.Println("我永远不会输出")
		<-done //这句也不会走到
	}()
	time.Sleep(time.Second * 1) //加入延时1秒
	fmt.Println("完成")
	/**
	控制台输出：
	我可能会输出
	完成
	*/
	/**
	结论：如果go routine中阻塞死，也可能不会把阻塞语句前的内容输出
	因为main已经跑完了，所以延时一会，等待go routine
	*/
}

// 演示无缓冲channel在不同的位置里面接收填充和接收
func test7() {
	// 编译通过，演示无缓冲channel在不同的位置里面接收填充和接收
	done := make(chan bool)
	go func() {
		done <- true // 直到<-done执行，否则这里阻塞死
		fmt.Println("我永远不会输出，除非<-done执行")
	}()
	<-done //这里接收，在输出完成之前，那么上面的语句将会走通
	fmt.Println("完成")

}

// 演示无缓冲channel在不同地方接收的影响
func test8() {
	// 编译通过，演示无缓冲channel在不同地方接收的影响
	done := make(chan bool)
	go func() {
		done <- true // 直到,<-done执行，否则这里阻塞死
		fmt.Println("我永远不会输出，除非<-done执行")
	}()
	fmt.Println("完成")
	<-done
	time.Sleep(time.Second)

}

// 没缓冲的channel使用close后，不会阻塞
func test9() {
	// 编译通过
	// 演示：没缓存的channel使用close后，不会阻塞
	done := make(chan bool)
	close(done)
	//done <- true //关闭了的，不能再往里面输入值
	<-done //这句是读取，但是上面已经关闭channel了，不会阻塞
	fmt.Println("完成")
}

// 没缓存的channel,在go routine里面使用close后，不会阻塞
func test10() {
	// 编译通过
	// 演示：没缓存的channel,在go routine里面使用close后，不会阻塞
	done := make(chan bool)
	go func() {
		close(done)
	}()
	//done<-true //关闭了的，不能再往里面输入值
	<-done //这句是读取，但是在上面已经关闭channel,不会阻塞
	fmt.Println("完成")
}

// 有缓冲的channel不会阻塞的例子
func test11() {
	// 编译通过
	// 有缓冲的channel不会阻塞的例子
	done := make(chan bool, 1)
	done <- true
	<-done
	fmt.Println("完成")
}

// 有缓冲的channel会阻塞的例子
func test12() {
	//编译通过
	// 有缓冲的channel会阻塞的例子
	done := make(chan bool, 1)
	//done <- true //注释这句
	<-done //虽然是有缓冲的，但是在没输入的情况下，读取，会阻塞
	fmt.Println("完成")
}

// 有缓冲的channel会阻塞的例子
func test13() {
	//编译不通过
	//有缓冲的channel会阻塞的例子
	done := make(chan bool, 1)
	done <- true
	done <- false //放第二个值的时候，第一个还没被人拿走，这时候才会阻塞，根据缓冲值而定
	fmt.Println("完成")
}

// 有缓冲的channel不会阻塞的例子
func test14() {
	//编译通过
	//有缓冲的channel不会阻塞的例子
	done := make(chan bool, 1)
	done <- true //不会阻塞在这里，等待读取
	fmt.Println("完成")

}

// 有缓冲的channel,如果在go routine中使用，一定要做适当的延时，否则会输出来不及，因为
// main已经跑完了，所以延时一会，等待go routine
func test15() {
	// 编译通过
	// 有缓冲的channel在go routine里面的例子
	done := make(chan bool, 1)
	go func() {
		//不会阻塞
		fmt.Println("我可能会输出")
		done <- true //如果把这个注释，也会导致<-done阻塞
		fmt.Println("我也可能会输出")
		<-done
		fmt.Println("别注释done<-true,不然我就输出不了了")
	}()
	// 1秒延时，去掉就可能上面的都不会输出也有可以输出，routine
	fmt.Println("完成")
	time.Sleep(time.Second * 1)
}

func main() {
	//test0()
	test1()
	//test2()
	//test3()
	//test4()
	//test5()
	//test6()
	//test7()
	//test8()
	//test9()
	//test10()
}
