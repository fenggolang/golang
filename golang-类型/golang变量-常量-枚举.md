#### golang内建变量类型

​	下面的rune是32位的，byte是8位的，这2个可以和整数混用的,byte是(u)int8整数的别名，rune是(u)int32整数的别名

- bool,string
- (u)int,(u)int8,(u)int16,(u)int32,(u)int64,uintptr
- byte rune
- float32，float64，complex64，complex128

#### golang强制类型转换

​	golang中没有隐式类型转换

```go
var a,b int = 3,4
var c int = math.Sqrt(a*a+b*b) //错误的写法
var c int = (int)math.Sqrt((float64(a*a+b*b))) //正确的写法，因为需要强制类型转换，golang没有隐式类型转换
```

```bash
package main

import (
	"fmt"
	"math"
)

//一直直角三角形，2边为3,4，求:斜边长度
func triangle(){
	var a,b int = 3,4
	var c int
	//c = math.Sqrt(a*a+b*b)
	c = (int)(math.Sqrt(float64(a*a+b*b)))
	fmt.Println(c)
}

func main(){
	//调用三角函数
	triangle()
}
```

#### golang常量

```go
package main

import (
	"math"
	"fmt"
)

//全局常量
const (
	aa="aa"
	bb="bb"
)

func consts(){
	//局部常量
	const filename = "abc.txt"
	// 定义下面2个常量a,b,没有定义其类型，当我们把这2个常量传递给math.Sqrt函数的
	// 时候，因为其参数是float64，编译器会自动把a,b识别为float64,而不用强制转换
	const a,b = 3,4
	var c int
	c = int(math.Sqrt(a*a+b*b))
	fmt.Println(filename,c)
}

func main(){
	consts()
	fmt.Println("全局常量",aa,bb)
}
```

#### golang特殊的常量

​	golang有一种特殊的常量：枚举，通过const块来实现枚举类型，下面的就叫枚举类型

```go
package main

import (
	"fmt"
)

//特殊的常量：枚举
func enums(){
	const(
		cpp=0
		java=1
		python=2
		golang=3
	)
	fmt.Println(cpp,java,python,golang)
}

func main(){
	enums()
}
```

​	golang中为这个枚举类型做了一个简化

```go
package main

import (
	"fmt"
)

//特殊的常量：枚举
func enums(){
	const(
		//iota默认是0，然后下面的以此递增
		cpp=iota
		_
		python
		golang
		javascript
	)
	fmt.Println(cpp,javascript,python,golang)

	//b,kb,mb,gb,tb,pb
	const(
		b = 1 << (10*iota)
		kb
		mb
		gb
		tb
		pb
	)
	fmt.Println("b:",b,",kb:",kb,",mb:",mb,",gb:",gb,",tb:",tb,",pb:",pb)
}

func main(){
	// 0 4 2 3
	// b: 1 ,kb: 1024 ,mb: 1048576 ,gb: 1073741824 ,tb: 1099511627776 ,pb: 1125899906842624
	enums()
}
```

