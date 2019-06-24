#### golang if语句

```go
package main

import (
"fmt"
"io/ioutil"
)

func main(){
	const filename ="abx.txt"
	contents,err:=ioutil.ReadFile(filename)
	if err!=nil{
		fmt.Println(err)
	}else {
		fmt.Printf("%s\n",contents)
	}

	// 下面这种写法说明：if的条件里可以赋值
	// if的条件里赋值的变量作用域就在这个if语句里
	if contents,err:=ioutil.ReadFile(filename);err==nil{
		fmt.Println(string(contents))
	}else {
		fmt.Println("不能打印文件内容",err)
	}
}
```

#### golang swtich语句

```go
package main

import "fmt"

// 注意：golang中的switch会自动添加break,除非使用fallthrough
func eval(a,b int,op string) int{
	var result int
	switch op{
	case "+":
		result=a+b
	case "-":
		result=a-b
	case "*":
		result=a*b
	case "/":
		result=a/b
	default:
		panic("不支持的运算符"+op)
	}
	return result
}
func main(){
	result:=eval(1,2,"+")
	fmt.Println("result:",result)
}

```

 注意，golang中的switch语句可以没有表达式

```go
package main

import "fmt"

func grade(score int) string{
	g:=""
	// switch后面可以没有表达式！！
	switch{
	case score <0 || score >100:
		panic(fmt.Sprintf("错误分数: %d",score))
	case score<60:
		g="F"
	case score <80:
		g="C"
	case score<90:
		g="B"
	case score<100:
		g="A"
	}
	return g
}
func main(){
	// output: F F C B A
	fmt.Println(
		grade(0),
		grade(59),
		grade(60),
		grade(82),
		grade(92),
		grade(100),
		//grade(102),
	)
}
```

---

#### golang for语句

​	golang中没有while,要用while其实用for就可以了

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// for的条件里不需要括号
// 十进制转二进制
func convertToBin(n int) string{
	result:=""
	if n ==0 {
		result="0"
	}

	// 省略初始条件，相当于while
	for ; n>0; n/=2{
		lsb:=n%2
		result = strconv.Itoa(lsb)+result
	}
	return result
}

func printFile(filename string) {
	file,err:=os.Open(filename)
	if err!=nil{
		panic(err)
	}

	scanner:=bufio.NewScanner(file)
	// 下面的for循环和其他语言的while循环作用一样
	for scanner.Scan(){
		//打印扫描出来的一行
		fmt.Println(scanner.Text())
	}
}

func forever(){
	// 死循环
	for {
		fmt.Println("abc")
	}
}
func main(){
	fmt.Println(
		convertToBin(5), // 101
		convertToBin(13),//1101
		convertToBin(3382887),
		convertToBin(0),
	)

	printFile("abx.txt")

	//forever()
}
```

总结：

​	for,if后面的条件没有括号，if条件里也可以定义变量

​	switch不需要break,switch后面可以不跟表达式，表达式在case后面跟也可以，也可以switch多个条件