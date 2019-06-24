### 查看godoc文档

​	在dos命令行里面执行

```bash
godoc -http=:8080
```

go语言中首字母大写是导出字段，首字母小写是私有字段，但是对同一个包来讲都是内部的私有字段在任何包内都是可见的。

方法访问权限，字段访问权限与package的关系

数组是值类型，将一个数组赋值给另一个数组时将复制一份新的元素；

切片是引用类型，因此在当传递切片时将引用同一指针，修改值将会影响其他的对象。

go定义接口方式是：

```bash
type Interface_name interface{
  方法签名集。。。
}
```

go语言中不需要显示的说我实现了哪一个接口，只要你拥有了接口所定义的方法，也就是说你的方法集是某一个接口的超集，那么这个时候你就默认实现了对应的接口

go语言中不存在继承的概念

go语言中没有方法重载的概念,但是方法是和某一个类型相绑定的，如果你这个方法所绑定的对象不同，那么方法名称和方法参数即使一样，这个方法整体也是不同的，因为绑定的对象类型不同，

go语言中没有class的概念，但是可以通过结构体和方法来模拟class

```go
package main

import (
	"fmt"
)

type A struct {
	Name string
}

func main(){
	a := A{}
	a.Print()
}

func (a A) Print(){
	fmt.Println("A")
}
```



#### 可见性规则

​	我们在其他语言中会有私有成员，公有成员变量等概念，一般是通过private,public等关键字来识别的，在Go语言中，使用的则是大小写来角色该常量、变量、类型、接口、结构或函数是否可以被外部包所调用。

​	根据约定，函数名首字母小写即为private私有，函数名首字母大写即为public公有，比如我们导入fmt包，我是调用包中的函数Println(),这个就是首字母大写的，因为他要被外部使用，需要public公有

### golang中是否存在继承

​	我们知道其他语言有继承，比如相同的属性，我们不必重复去写，只需继承父类的公共属性即可。遗憾的是Go没有继承，但Go有`组合`.

​	golang不支持完整的面向对象思想，它没有继承，多态则完全依赖接口实现。golang只能模拟继承，其本质是组合，只不过golang语言为我们提供了一些语法糖使其看起来达到了继承的效果。

​	golang的设计理念是大道至简，传统的继承概念在golang中已经显得不是那么必要，golang通过接口去实现多态.golang中的模拟继承并不等价于面向对象中的继承关系。



### golang中的接口

​	什么是接口。其实通俗地讲，接口就是一个协议，规定了一组成员,相当于是一份契约，它规定了一个对象所能提供的一组操作。

​	要理解golang中接口的概念我们最好还是先来看看别的现代语言是如何实现接口的。C++没有提供interface这样的关键字，它通过纯虚基类实现接口，而java则通过interface关键字声明接口。它们有个共同特征就是一个类要实现该接口必须进行显示的声明，如下是java方式：

```java
interface IFoo {
    void Bar();
}
class Foo implements IFoo { 
    void Bar(){}
}
```

这种必须明确声明自己实现了 某个接口的方式我们称为侵入式接口。关于侵入式接口的坏处我们这里就不再详细讨论，看java庞大的继承体系及其繁复的接口类型我们就可以窥之一二了。

golang则采取了完全不同的设计理念，在Go语言中，一个类只需要实现了接口要求的所有函数，我们就说这个类实现了该接口， 例如：

```go
type IWriter interface {
    Write(buf [] byte) (n int, err error)
}
type File struct {
    // ...
}
func (f *File) Write(buf [] byte) (n int, err error) {
    // ...
}
```

非侵入式接口一个很重要的好处就是去掉了繁杂的继承体系.Go语言的标准库，再也不需要绘制类库的继承树图。在Go中，类的继承树并无意义，你只需要知道这个类实现了哪些方法，每个方法是啥含义就足够了。 

其二，实现类的时候，只需要关心自己应该提供哪些方法，不用再纠结接口需要拆得多细才 合理。接口由使用方按需定义，而不用事前规划。

其三，不用为了实现一个接口而导入一个包，因为多引用一个外部的包，就意味着更多的耦 合。接口由使用方按自身需求来定义，使用方无需关心是否有其他模块定义过类似的接口。

​	go的精髓是interface和channel

### golang中的错误处理机制

​	golang错误处理机制"比较原始",在Go语言中处理错误的基本模式是：函数通常返回多个值，其中最后一个值是error类型，用于表示错误类型极其描述；调用者每次调用完一个函数，都需要检查这个error并进行相应的错误处理：if err != nil { /*xxx*/ }。

### golang的特点

- C-like 语法风格
- 强一致类型（静态语言）
- struct组合复杂类型
- function和Method
- 没有异常处理（Error is value）
- 基于首字母的访问控制
- 多余的import和变量会引起编译错误
- 内置GC
- 完备的标准库包（网络编程、系统编程、互联网应用）
- 支持各种编程范式：过程式编程、面向对象编程、函数式编程

### golang中的特殊函数

###### defer函数

```
//在函数返回之前，先进后出执行defer语句，一般在需要资源回收的时候使用
```

###### main函数

```
func main(){ //... } //package main中必须包含一个main函数，是程序的入口函数
```

###### init函数

```
func init(){ //some init logic... } /*一个package中最好只写一个init 1.如果package中需要import其他包，则先导入其他包 2.所有包import完毕后，依次初始化cons->var->init() */
```

### golang中的函数与方法

​	在golang中，函数function与方法method有什么区别？

官方解释：

```go
// go里一个type的function即是method。
A method is a function with a special receiver argument.
```

**method和function的关系**:

​	method是特殊的function，定义在某一特定的类型上，通过类型的实例来进行调用，这个实例被叫receiver。

golang的函数定义如下：

```go
func function_name( [parameter list] ) [return_types] {
   函数体
}
```

golang的方法定义如下：

​	Golang 没有类，只有结构体。不过Golang可以在结构体类型上定义方法，其实就是配合结构体的函数。方法和之前讲过的函数是有些小区别的——对应的结构体信息（也叫“方法接受者”），出现在方法定义中。

```go
// 我们只需要在普通函数前面加个接受者（receiver，写在函数名前面的括号里面），这样编译器就知道这个函数（方法）属于哪个struct了。
// method是附属在一个给定的类型上，语法和函数的声明语法几乎一样，只是再func后面增加了一个recevier（也就是method所依从的主体）
// 如果需要改变结构体中的值的话下面的receive需要指针传递:(r *ReceiverType)
// 那么有哪些情况下必须是使用指针接收者
// 首先，避免在每个方法调用中拷贝值（如果值类型是大的结构体的话会更有效率）；
// 其次，方法可以修改接收者指向的值。
func (r ReceiverType) funcName(parameters) (results){
  方法体
}
```

```go
package main

import (
	"math"
	"fmt"
)

type Vertex struct{
	X,Y float64
}
// 在func和函数名之间声明的的参数即receiver, 下面例子中类型为Vertex的v即receiver。
func (v Vertex) Abs() float64{
	return math.Sqrt(v.X*v.X+v.Y*v.Y)
}

func main(){
	vertex:=Vertex{1.0,2.0}
	fmt.Println(vertex.Abs())
}
```



参考：[Go基础学习四之函数function、结构struct、方法method](https://segmentfault.com/a/1190000011446643#articleHeader11)

​	[golang 语法常识——function and method](https://radrupt.com/golang-yu-fa-chang-shi/)



### golang中的结构体与数组

​	结构体可以包含多种数据类型，数组只能是单一类型的数据集合。如果要访问结构体成员，需要使用点号 (.) 操作符，格式为：`"结构体.成员名"`。

结构体定义和使用如下：

```go
package main

import "fmt"

// 定义图书馆书籍的属性
type library struct{
	Title string
	Author string
	Subject string
	ID int
}

func main(){
	librarys:=library{"水浒传","王鹏程","武侠",00002020}
	fmt.Println(librarys)

}
```

数组使用如下：

```go
package main

import "fmt"

func main(){
	a :=[4]int{1,2,3,4}
	fmt.Println(a)
}
```



### golang中的数组与切片

​	数组是内置(build-in)类型,是一组同类型数据的集合，它是值类型，通过从0开始的下标索引访问元素值。在初始化后长度是固定的，无法修改其长度。当作为方法的入参传入时将复制一份数组而不是引用同一指针。数组的长度也是其类型的一部分，通过内置函数len(array)获取其长度。数组是值类型。

​	数组的长度不可改变，在特定场景中这样的集合就不太适用，Go中提供了一种灵活，功能强悍的内置类型**Slices**切片,与数组相比切片的长度是不固定的，可以追加元素，在追加时可能使切片的容量增大。切片中有两个概念：一是**len长度**，二是**cap容量**，长度是指已经被赋过值的最大下标+1，可通过内置函数len()获得。容量是指切片目前可容纳的最多元素个数，可通过内置函数cap()获得。切片是引用类型，因此在当传递切片时将引用同一指针，修改值将会影响其他的对象。切片可以通过数组来初始化，也可以通过内置函数make()初始化 .初始化时len=cap,在追加元素时如果容量cap不足时将按len的**2**倍扩容。

### golang中的访问权限

​	因为Go是以大小写来区分是公有还是私有，但都是针对包级别的，所以在包内所有的都能访问，而方法绑定本身只能绑定包内的类型，所以方法可以访问接收者所有成员。如果是包外调用某类型的方法，则需要看方法名是大写还是小写，大写能被包外访问，小写只能被包内访问。

### 在Go中怎样判断数据类型

```go
//  用反射：reflect.TypeOf(x)
package main

import (
        "fmt"
        "reflect"
)

func main() {
        var x float64 = 3.4
        fmt.Println("type:", reflect.TypeOf(x))
}
```

### Go Json处理

​	Golang 虽然自己就带了 JSON (encoding/json) 处理的库，但是不如像Python这种动态语言的语法灵活。encoding/json 最大的问题是不够灵活，需要预先定义很多的 struct 来进行编解码，这样对于处理结构不定的 JSON 文件非常不方便。

​	Go处理多层嵌套的json



### Go中的interface{}作用

**golang中空的interface即interface{}可以看作任意类型, 即C中的void *.**

```go
// BodyJson is the document as a serializable JSON interface.
func (s *IndexService) BodyJson(body interface{}) *IndexService {
	s.bodyJson = body
	return s
}

// BodyString is the document encoded as a string.
func (s *IndexService) BodyString(body string) *IndexService {
	s.bodyString = body
	return s
}
```



### go json解析Marshal和Unmarshal

### 解决go get golang.org/x/xxx被墙的问题

```bash
# 会失败
go get golang.org/x/tools/cmd/goimports
# 
```

使用如下办法1

```bash
go get -v -u github.com/yanjunhui/god

# 然后windows下添加环境变量，下面就有god.exe
D:\Program Files\GOPATH\bin
```

如下办法二：

```bash
# 直接下载好包到windows对应目录即可
D:\Program Files\GOPATH\src\golang.org\x\tools
```





