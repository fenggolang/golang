#### golang interface定义

​	在Go中，**接口是一组方法签名**。 当一个类型为接口中的所有方法提供定义时，它被称为实现该接口。它与面向对象编程非常相似。接口指定类型应具有的方法，类型决定如何实现这些方法。

#### 创建和实现interface

​	在Golang中只要实现了接口定义的方法，就是实现了该interface

```go
package main

import "fmt"

// 定义接口:
type VowelsFinder interface {
	// 查看26个字母中的元音字母
	FindVowels() []rune
}

type MyString string

// 实现接口：
func (ms MyString) FindVowels() []rune {
	var vowels []rune
	for _, rune := range ms {
		if rune == 'a' || rune == 'e' || rune == 'i' || rune == 'o' || rune == 'u' {
			vowels = append(vowels, rune)
		}
	}
	return vowels
}
func main() {
	name := MyString("Sam Anderson") //类型转换
	var v VowelsFinder               //定义一个接口类型的变量
	v = name
	fmt.Printf("Vowels are %c", v.FindVowels())
}

// output:
Vowels are [a e o]
Process finished with exit code 0
```

---

#### 接口的实际用途

​	如果我们将上面的输出代码：

```go
fmt.Printf("Vowels are %c", v.FindVowels())
```

​	替换为：

```go
fmt.Printf("Vowels are %c", name.FindVowels())
```

​	程序同样的输出，而没有使用我们定义的接口。（v变量删除定义）。

下面我们通过案例来解释：

​	假设某公司有两类员工，一类普通员工和一类高级员工，但是基本薪资是相同的，高级员工多拿奖金。计算公司为员工的总开支。

```go
package main

import "fmt"

//薪资计算器接口
type SalaryCalculator interface {
	CalculateSalary() int
}

//普通挖掘机员工
type Contract struct {
	empId    int
	basicpay int
}

//有蓝翔技校证的员工
type Permanent struct {
	empId    int
	basicpay int
	jiangjin int //奖金
}

func (p Permanent) CalculateSalary() int {
	return p.basicpay + p.jiangjin
}
func (c Contract) CalculateSalary() int {
	return c.basicpay
}

//总开支
func totalExpense(s []SalaryCalculator) {
	expense := 0
	for _, v := range s {
		expense = expense + v.CalculateSalary()
	}
	fmt.Printf("总开支：$%d", expense)
}

func main() {
	pemp1 := Permanent{1, 3000, 10000}
	pemp2 := Permanent{2, 3000, 20000}
	cemp1 := Contract{3, 3000}
	employees := []SalaryCalculator{pemp1, pemp2, cemp1}
	totalExpense(employees)
}

// output:
总开支：$39000
```

---

#### 接口的内部表现

​	一个接口可以被人围殴是由一个元组(类型，值)在内部表示的。type是接口的基础具体类型，value是具体类型的值

```go
package main

import "fmt"

type Test interface {
	Tester()
}
type MyFloat float64

func (m MyFloat) Tester() {
	fmt.Println(m)
}

func describe(t Test) {
	fmt.Printf("Interface 类型 %T,值：%v\n", t, t)
}

func main() {
	var t Test
	f := MyFloat(89.7)
	t = f
	describe(t)
	//describe(f)
	t.Tester()
	//f.Tester()
}
// output:
Interface 类型 main.MyFloat,值：89.7
89.7
```

---

#### 空接口

​	具有0个方法的接口称为空接口。它表示为interface{}。由于空接口有0个方法，所有类型都实现了空接口。

```go
package main

import "fmt"

func describe(i interface{}) {
	fmt.Printf("Type = %T,value = %v\n", i, i)
}

func main() {
	//任何类型的变量传入都可以
	s := "Hello World,Hello, Golang"
	i := 55
	strt := struct {
		name string
	}{
		"风王集团",
	}

	describe(s)
	describe(i)
	describe(strt)
}
// output:
Type = string,value = Hello World,Hello, Golang
Type = int,value = 55
Type = struct { name string },value = {风王集团}
```

---

#### 类型断言

​	类型断言用于提取接口的基础值，语法：i.(T)

```go
package main

import "fmt"

func assert(i interface{}) {
	s := i.(int)
	fmt.Println(s)
}

func main() {
	var s interface{} = 55
	assert(s)
}
// output
55
```

程序打印的是int值， 但是如果我们给s 变量赋值的是string类型，程序就会panic。所以可以将以上程序改写为：

```go
package main

import "fmt"

func assert(i interface{}) {
	v, ok := i.(int)
	fmt.Println(v, ok)
}

func main() {
	var s interface{} = 55
	assert(s)
	var i interface{} = "Steven Paul"
	assert(i)
}
// output:
55 true
0 false
```

---

#### 类型判断

类型判断的语法类似于类型断言。在类型断言的语法i.(type)中，类型type应该由类型转换的**关键字**type替换。

```go
package main

import "fmt"

func findType(i interface{}) {
	switch i.(type) {
	case string:
		fmt.Printf("String: %s\n", i.(string))
	case int:
		fmt.Printf("Int: %d\n", i.(int))
	default:
		fmt.Printf("Unknown type\n")
	}
}

func main() {
	findType("Chain")
	findType(77)
	findType(89.98)
}

// output:
String: Chain
Int: 77
Unknown type
```

​	还可以将类型与接口进行比较。如果我们有一个类型并且该类型实现了一个接口，那么可以将它与它实现的接口进行比较。

```go
package main

import "fmt"

type Describer interface {
	Describe()
}
type Person struct {
	name string
	age  int
}

func (p Person) Describe() {
	fmt.Printf("%s is %d years old", p.name, p.age)
}

func findType(i interface{}) {
	switch v := i.(type) {
	case Describer:
		v.Describe()
	default:
		fmt.Printf("unkonwn type\n")
	}
}

func main() {
	findType("Chain")
	p := Person{
		name: "China Wang",
		age:  25,
	}
	findType(p)
}
// output:
unkonwn type
China Wang is 25 years old
```

---

go语言的面向对象只支持封装，不支持继承多态，go语言的多态是通过接口来实现的，go语言中的接口比传统的面向对象的语言的接口要灵活很多。严格说，go属于结构化系统。go的接口实现是隐式的，只要实现接口中的方法即可表示实现了该接口。

#### 示例一

```go
package main

import (
	"fmt"
)

// 接口
type Retriever interface {
	Get(url string) string
}

// 结构体
type RetrieverStruct struct {
	Content string
}

// 实现接口的方法
func (r RetrieverStruct ) Get(url string) string {
	return r.Content
}

// 函数
func download(r Retriever) string {
	return r.Get("www.abc.com.cn")
}

func main() {
	var r Retriever
	r = RetrieverStruct{"this is a my house"}
	fmt.Println(download(r))
}

// output
this is a my house
```



