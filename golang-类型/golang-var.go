package main

import "fmt"

// 函数外定义变量
var aa = 3
var ss = "kk"

// 函式外面定义的变量不能使用:=,下面语法错误
//bb:=1

// 下面是简化写法
var (
	aaa = 3
	sss = "kk"
	bbb = true
)

func variableZeroValue() {
	var a int
	var s string

	fmt.Println(a, s)
}

func variableInitialValue() {
	var a, b int = 3, 4
	var s string = "abc"
	fmt.Println(a, b, s)
}

func variableTypeDeduction() {
	var a, b, c, s = 3, 4, true, "def"
	fmt.Println(a, b, c, s)
}

// 使用简化定义变量的写法:=
func variableShorter() {
	a, b, c, s := 3, 4, true, "def"
	b = 5
	fmt.Println(a, b, c, s)
	var (
		aa = 1
		bb = 1
	)
	fmt.Println(aa, bb)
}
func main() {
	fmt.Println("Hello world")
	variableZeroValue()
	variableInitialValue()
	variableTypeDeduction()
	variableShorter()
	fmt.Println(aaa, bbb, sss)
}

// 冒号等于:=只能在函数内部使用