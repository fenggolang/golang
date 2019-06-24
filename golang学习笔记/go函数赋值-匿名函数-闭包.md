#### go函数赋值

函数赋值就是把函数名给一个变量，然后变量使用的时候加个()

```go
package main

import "fmt"



//go语言中一切皆类型，我们也可以把一个函数赋值给一个变量
func main() {
	a:=A
	a()
}

func A(){
	fmt.Println("Func A")
}
```

#### go匿名函数

匿名函数就是没有名称的函数

```go
package main

import "fmt"



//go语言中一切皆类型，我们也可以把一个函数赋值给一个变量
func main() {
	a:=func (){
		fmt.Println("Func A")
	}
	//变量赋值得到一个匿名函数，那么a()就是可以当成一个函数使用
	a()
}
```

#### go闭包

闭包函数的作用就是返回一个匿名函数

```go
package main

import "fmt"



//go语言中一切皆类型，我们也可以把一个函数赋值给一个变量
func main() {
	f:=closure(10) //调用闭包函数，因为闭包函数返回的是一个匿名函数，所以f相当于得到了一个匿名函数的赋值
	fmt.Println(f(1)) //f(1)中的参数值1就是传递给闭包函数的里面的参数y的值 10+1=11
	fmt.Println(f(2)) //f(2)中的参数值1就是传递给闭包函数的里面的参数y的值 10+2=12
}

// 闭包函数
func closure(x int) func(int) int{
	//闭包函数中打印x的内存地址：0xc04200a240
	fmt.Println("%p\n",&x)
	return func (y int) int{
		//匿名函数中打印x的内存地址：0xc04200a240，地址一样说明使用的是同一个x的地址引用
		fmt.Println("%p\n",&x)
		return x+y
	}
}
```

#### go 析构函数

go中的defer函数执行方式类似其它语言中的析构函数，在函数体执行结束后按照调用顺序的相反顺序逐个执行

常用于资源清理、文件关闭、解锁以及记录时间等操作

```go
package main

import "fmt"



//go语言中一切皆类型，我们也可以把一个函数赋值给一个变量
func main() {
	fmt.Println("a")
	defer fmt.Println("b")
	defer fmt.Println("c")
}

// output a \n c \n b
```

```go
package main

import "fmt"


func main() {
	for i:=0;i<3;i++{
		defer fmt.Println(i)
	}
}

// 输出 2,1,0
```

```go
func main() {
	for i:=0;i<3;i++{
		defer func(){
			fmt.Println(i)
		}()
	}

}

// 输出3 3 3，因为defer里面的打印的i是引用的外面for循环的i的地址，当i++变为3的时候for循环结束开始输出defer,此时defer里面的输出的i已经变为了3所以输出3 3 3
```

#### go异常机制

Go 没有异常机制，但有 panic/recover 模式来处理错误，Panic 可以在任何地方引发，但recover只有在defer调用的函数中有效。

```go
package main

import "fmt"



func main() {
	A()
	B()
	C()
}

func A(){
	fmt.Println("Func A")
}
func B(){
	defer func(){
		if err:=recover();err!=nil{
			fmt.Println("Recover in B")
		}
	}()
	panic("Panic in B")
}
func C(){
	fmt.Println("Func C")
}

// 输出 Func A,Recover in B,Func C
```

