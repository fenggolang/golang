#### panic作用

- 停止当前函数执行;
- 一直向上返回，执行每一层的defer
- 如果没有遇见recover,程序退出
- panic不要随便用，不要程序导出都是panic

```go
package main

import "fmt"

func tryDefer() {
	for i := 0; i < 3; i++ {
		defer fmt.Println(i)
		if i == 3 {
			panic("printed too many")
		}
	}
}
func main() {
	tryDefer()
}
```

---

#### recover作用

- 仅在defer调用中使用
- 获取panic的值
- 如果无法处理，可重新panic

```go
package main

import (
	"fmt"
)

func tryRecover() {
	defer func() {
		r := recover()
		if err, ok := r.(error); ok {
			fmt.Println("Error occurred:", err)
		} else {
			panic(r)
		}
	}()

	// 第一种
	//panic(errors.New("this is an error"))

	// 第二种
	panic(123)

	// 第三种
	//b := 0
	//a := 5 / b
	//fmt.Println(a)
}
func main() {
	tryRecover()
}

```

#### panic和error

```bash
# 什么使用用error,什么时候用panic,注意：panic尽量不要用

意料之中的：使用error,如：文件打不开

意料之外的：使用panic,如：数组越界

defer + panic + recover三者加起来进行错误处理
```





