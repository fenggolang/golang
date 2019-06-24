package main

import (
	"fmt"
	"os"
)



type point struct {
	x,y int
}

func main() {
	// Go提供了几种打印格式，用来格式化一般的Go值，例如
	// 下面的%v打印了一个point结构体的对象的值,原样输出
	p:=point{1,2}

	// fmt.Printf:根据格式说明符打印格式并写入标准输出。返回写入的字节数和任何错误信息
	// func Printf(format string, a ...interface{}) (n int, err error)
	// Printf内部调用的是Fprintf(os.Stdout, format, a...)
	n,err:=fmt.Printf("%v\n",p) // 输出：{1 2}
	fmt.Println(n,err)

	//  fmt.Println:使用其操作数的默认格式格式化并写入标准输出。返回写入的字节数和任何错误信息
	// Println内部调用的是Fprintln(os.Stdout, a...)
	fmt.Println(p) // 输出并且换行：{1 2}
	fmt.Print(p) // 输出但不换行：{1 2}

	// 如果所格式化的值是一个结构体对象，那么`%+v`的格式化输出
	// 将包括结构体的成员名称和值
	fmt.Println()
	fmt.Printf("%+v\n",p) // 输出：{x:1 y:2}

	// `%#v`格式化输出将输出一个值的Go语言表示方式。
	fmt.Printf("%#v\n",p) // 输出：main.point{x:1, y:2}

	// 使用`%T`来输出一个值的数据类型
	fmt.Printf("%T\n",p) // 输出：main.point

	// 格式化布尔型变量
	fmt.Printf("%t\n",true) // 输出：true

	// 有很多的方式可以格式化整型，使用`%d`是一种
	//标准的以10(decimal)进制来输出整型的方式
	fmt.Printf("%d\n",10086) // 10086

	// 输出整型的二进制(binary)表示方式
	fmt.Printf("%b\n",10086) // 10011101100110

	// 打印出整型数值所对应的字符
	fmt.Printf("%c\n",10086) // ❦

	// 使用`%x`输出一个值的16进制表示方式
	fmt.Printf("%x\n",10086) // 2766

	// 浮点型数值也有几种格式化方法，最基本的一种是`%f`
	fmt.Printf("%f\n",12.3) // 12.300000

	// `%e`和`%E`使用科学计数法来输出整型
	fmt.Printf("%e\n",123400000.0) // 1.234000e+08
	fmt.Printf("%E\n",123400000.0) // 1.234000E+08

	// 使用`%s`输出基本的字符串
	fmt.Printf("%s\n","\"string\"") // "string"

	// 输出像Go源码中那样带双引号的字符串，需使用`%q`
	fmt.Printf("%q\n", "\"string\"") // "\"string\""
	fmt.Printf("%q\n", "string") // "string"

	// `%x`以16进制输出字符串，每个字符串的字节用2个字符输出
	fmt.Printf("%x\n","hex this")

	// 使用`%p`输出一个指针的值
	fmt.Printf("%p\n",&p) // 0xc042054080

	// 当输出数字的时候，经常需要去控制输出的宽度和精度。
	// 可以使用一个位于%后面的数字来控制输出的宽度，默认
	// 情况下输出是右对齐的，左边加上空格
	fmt.Printf("|%6d|%6d|\n",12,345) // |    12|   345|

	// 你也可以指定浮点数的输出宽度，同时你还可以指定浮点数
	// 的输出精度
	fmt.Printf("|%6.2f|%6.2f|\n", 1.2, 3.45)

	// 你也可以指定输出字符串的宽度来保证它们输出对齐。默认
	// 情况下，输出是右对齐的
	fmt.Printf("|%6s|%6s|\n", "foo", "b")

	// 为了使用左对齐你可以在宽度之前加上`-`号
	fmt.Printf("|%-6s|%-6s|\n", "foo", "b")

	// `Printf`函数的输出是输出到命令行`os.Stdout`的，你
	// 可以用`Sprintf`来将格式化后的字符串赋值给一个变量
	s:=fmt.Sprintf("%s","string") // string
	fmt.Println(s)
	// 你也可以使用`Fprintf`来将格式化后的值输出到`io.Writers`,也就是你自己指定输出到哪里，Printf默认输出到os.Stdout
	fmt.Fprintf(os.Stderr,"an %s\n", "error")
	fmt.Fprintf(os.Stdout, "an %s\n", "error")
}

/**
总结：
fmt.Printf指定格式化输出，返回输出到标准输出的字节数和error
fmt.Sprintf指定格式化输出返回字符串，这样可以赋值给一个变量而不是输出到标准输出
 */