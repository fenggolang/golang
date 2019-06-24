#### golang函数与方法

 - golang函数可以返回多个值;返回值写在最后面
 - golang函数返回多个值时可以起名字，仅用于非常简单的函数，对于调用者来说没区别，随便用什么变量名称接收返回值;
 - golang主要还是函数式编程语言,是golang中的一等公民。函数的返回值，函数体内都可以是函数，函数的参数也可以是函数
 - golang中没有函数重载，操作符重载，默认参数等花哨的语法，它有一个可变参数列表的的函数语法sumArgs(values...int) int
 - 在c/c++,java等语言中，函数和方法没有明显的区别(可以理解为只是同一个东西的两个名字而已)，但是在golang中是完全两个不同的东西。官方解释：方法是包含了接受者的函数，接收者可以是自己定义的一个类型，这个类型可以是struct,interface,甚至我们可以重定义基本数据类型。可以把接受者当作一个class，而这些方法就是类的成员函数
 - golang中有2个特殊的函数：main函数和init函数；main函数作为一个程序的入口，只能有一个。init函数在每个package包是可选的，可有可无，甚至可以有多个(建议一个package中只能有一个init函数)，init函数自动调用不用我们手动调用。

操作示例

```go
package main

import (
	"fmt"
	"math"
	"reflect"
	"runtime"
)

// 返回一个值的函数
func eval(a,b int,op string) (int,error) {
	switch op{
	case "+":
		return a+b,nil
	case "-":
		return a-b,nil
	case "*":
		return a*b,nil
	case "/":
		//return a/b
		q,_:=div(a,b)
		//只能返回一个值
		return q,nil
	default:
		return 0,fmt.Errorf("不支持运算: %s", op)
		//panic("不支持的算数运算："+op)
	}
}

// 返回2个值的函数
// 13/3=4...1
//func div(a,b int) (int,int) {
// 给返回值取名字，q是商，r是余数,建议这种直接return
func div(a,b int) (q,r int) {
	return a / b, a%b
}

// 函数返回值取名之后直接,不建议分开，建议上面一种方式
func div2(a,b int) (q, r int){
	q = a/b
	r =  a%b
	return
}

// 参数也是函数
func apply(op func(int,int) int,a,b int) int{
	// 反射
	p:=reflect.ValueOf(op).Pointer()
	opName:=runtime.FuncForPC(p).Name()
	fmt.Printf("Calling function %s with args \n" + "(%d, %d)", opName,a,b)
	return op(a,b)
}

// 重写系统函数
func pow(a,b int) int {
	return int(math.Pow(float64(a),float64(b)))
}

// 可变参数长度函数
func sum(numbers ...int) int{
	s:=0
	for i:=range numbers{
		s+=numbers[i]
	}
	return s
}
func main(){
	fmt.Println(eval(1,2,"+"))

	//fmt.Println("13除3的结果",div(13,3)) //错误的写法，div返回2个值就不能在前面加字符串了。
	fmt.Println(div(13,3))
	// 如果函数div返回值取了名字，则调用此函数的时候，编辑器一般会自动把函数返回值赋值给其名称，自动生成
	q,r:=div(13,3)
	fmt.Println(q,r)

	fmt.Println(div2(13,3))


	if result,err:=eval(3,4,"x");err!=nil{
		fmt.Println("出错了",err)
	}else{
		fmt.Println("成功了",result)
	}

	fmt.Println(apply(pow,3,4))

	// 直接使用匿名函数
	fmt.Println(apply(
		func(a int,b int) int{
			return int(math.Pow(float64(a),float64(b)))
		},3,4))

	//
	fmt.Println(sum(1,2,3,4,5))
}
```

---

#### 函数式编程 vs 函数指针

- 函数是一等公民：参数、变量、返回值都可以是函数；
- 高阶函数
- 函数-->闭包

```go
package main

import "fmt"

//
func adder() func(int) int {
	sum := 0
	return func(v int) int {
		sum += v
		return sum
	}
}
func main() {
	a := adder()
	for i := 0; i < 10; i++ {
		fmt.Printf("0 + 1 + ... + %d = %d\n", i, a(i))
	}
}

// output
0 + 1 + ... + 0 = 0
0 + 1 + ... + 1 = 1
0 + 1 + ... + 2 = 3
0 + 1 + ... + 3 = 6
0 + 1 + ... + 4 = 10
0 + 1 + ... + 5 = 15
0 + 1 + ... + 6 = 21
0 + 1 + ... + 7 = 28
0 + 1 + ... + 8 = 36
0 + 1 + ... + 9 = 45
```

#### go闭包应用:斐波拉契数列

```go
package main

import "fmt"

func fibonacci() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func main() {
	f := fibonacci()

	fmt.Println(f()) //1
	fmt.Println(f()) //1
	fmt.Println(f()) //2
	fmt.Println(f()) //3
	fmt.Println(f()) //5
	fmt.Println(f()) //8
	fmt.Println(f()) //13
	fmt.Println(f()) //21
}

```



​	
