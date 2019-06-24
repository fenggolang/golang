

#### go函数

一般的函数定义叫做函数，定义在结构体上面的函数叫做该结构体的方法。

#### go结构体

go语言中没有继承的说法，只有组合的说法，go中没有class的概念，那么struct就承担了一个面向对象的功能。struct承担了class的一个角色

#### go方法

go语言中没有方法重载的概念，如果方法所绑定的类型不同

一般的函数定义叫做函数，定义在结构体上面的函数叫做该结构体的方法。

从某种意义上说，方法是函数的“语法糖”。当函数与某个特定的类型绑定，那么它就是一个方法。也证因为如此，我们可以将方法“还原”成函数。

```go
package main

import "fmt"



type A struct {
	Name string
}
type B struct {
	Name string
}

func main() {
	a := A{}
	a.Print()
	fmt.Println(a.Name)

	b := B{}
	b.Print()
	fmt.Println(b.Name)
}

//方法
func (a *A) Print() {
	a.Name = "AA"
	fmt.Println("A")
}

//func (a A) Print(x, y int) { //和上面是同一个方法，报错
//	fmt.Println("A")
//}

func (b B) Print() {
	b.Name = "BB"
	fmt.Println("B")
}
```

#### go接口

go语言中所有的类型都实现了空接口,相当于void

 //接口是一个方法签名的集合//所谓方法签名，就是指方法的声明，而不包括实现。

```go
//定义一个空接口
type empty interface {
}
```



```GO
package main

import "fmt"



//定义一个USB接口，接口里面含有2个没有实现的方法,go中的接口只包含方法
type USB interface {
	Name() string
	//Connect()
	Connector //也可以使用嵌入接口，类似于接口继承接口，但是go语言中应该叫组合接口
}

type Connector interface {
	Connect()
}

//定义一个结构体
type PhoneConnector struct {
	name string
}

//定义实现接口的方法
func (pc PhoneConnector) Name() string {
	return pc.name
}

func (pc PhoneConnector) Connect() {
	fmt.Println("Connected:", pc.name)
}
func main() {
	var usb USB
	usb = PhoneConnector{"PhoneConnector"}
	usb.Connect()
	Disconnect(usb)

}

//func Disconnect(usb USB) {
func Disconnect(usb interface{}) {
	switch v := usb.(type) {
	case PhoneConnector:
		fmt.Println("Disconnected:", v.name)
	default:
		fmt.Println("Unknown decive.")
	}
}
```

```go
package main

import (
	"fmt"
	"math"
)

//接口是一个方法签名的集合
//所谓方法签名，就是指方法的声明，而不包括实现。

//这里定义了一个最基本的标识集合形状的方法的接口
type geometry interface {
	area() float64
	perim() float64
}

//这里我们要让正方形square和圆形circle实现这个接口
type square struct {
	width, height float64
}
type circle struct {
	//半径
	radius float64
}

//在Go中实现一个接口，只要实现该接口定义的所有方法即可
//下面是正方形实现的接口
func (s square) area() float64 {
	return s.width * s.height
}
func (s square) perim() float64 {
	return 2*s.width + 2*s.height
}

//圆形实现的接口
func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

//如果一个函数的参数是接口类型，那么我们可以使用命名接口
//来调用这个函数
//比如这里的正方形square和圆形circle都实现了接口geometry,
//那么它们都可以作为这个参数为geometry类型的函数的参数。
//在measure函数内部，Go知道调用哪个结构体实现的接口方法。
func measure(g geometry) {
	fmt.Println(g)
	fmt.Println(g.area())
	fmt.Print(g.perim())
}
func main() {
	s := square{width: 3, height: 4}
	c := circle{radius: 5}

	//这里circle和square都实现了geometry接口，所以
	//circle类型变量和square类型变量都可以作为measure
	//函数的参数
	measure(s) //12 14
	measure(c) //
}
```

