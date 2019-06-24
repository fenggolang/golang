#### golang指针示例1

```go
package main

import "fmt"

// 指针
// 定义一个普通变量和一个指针变量，把普通变量的地址赋值给指针变量
// 然后操作指针变量，间接修改普通变量的值

func main(){
	var a int = 2
	// c语言中是int*代表指针，go语言中是*int代表指针
	var p *int =&a
	*p = 3
	fmt.Println("普通变量的值是",a)
}
// output: 3
```

---

#### golang指针的特点

- 指针不能运算(c语言的指针有运算)

---

#### golang的参数传递

##### 值传递

​	是对值做了一份拷贝

##### 引用传递

​	操作值其实是操作变量的地址

​	Go语言使用值传递还是引用传递？

注意**Go语言只有值传递一种方式**，也就是说只要涉及到函数的参数传递，都是拷贝一份。，go语言通过指针达到其他语言中的一种引用传递的效果。

​	下面用go实现交换2个变量的值(通过指针传递)

```go
package main

import "fmt"

func main(){
	//交换2个变量的值
	a1,b1:=2,3
	fmt.Println("原始a=",a1,"原始b=",b1)
	swap(&a1,&b1)
	fmt.Println("现在a=",a1,"现在b=",b1)
}

// go语言实现交换2个变量的值
func swap(a1,b1 *int){
	*a1,*b1=*b1,*a1
}
```

​	下面用go实现交换2个变量的值(通过返回值实现)

```go
package main

import "fmt"

func main(){
	//交换2个变量的值
	a1,b1 :=2,3
	fmt.Println("原始a=",a1,"原始b=",b1)
	a1,b1=swap(a1,b1)
	fmt.Println("现在a=",a1,"现在b=",b1)
}

// go语言实现交换2个变量的值
func swap(a1,b1 int) (a,b int){
	return b1,a1
}
```

