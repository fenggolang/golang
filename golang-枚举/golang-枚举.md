#### golang枚举

​	golang中没有真正意义上的枚举类型(java中的枚举类型是enum)，但是golang也可以实现枚举，通过golang的关键字iota和const的组合实现枚举。

**iota，特殊常量，可以认为是一个可以被编译器修改的常量。** 在每一个const关键字出现时，被重置为0,然后再下一个const出现之前，每出现一次iota,其所代表的数字会自动增加1，关键字iota定义常量组从0开始按行计数的自增枚举值。

---

#### 示例一

golang枚举实现：输出星期天对应0，星期一对应1，星期六对应6

```go
package main

import (
	"fmt"
)

const (
	Sunday = iota
	Monday // 通常省略后续行表达式
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func main() {
	fmt.Println(Sunday,Monday,Tuesday,Wednesday,Thursday,Friday,Saturday)
}
// output:
0 1 2 3 4 5 6
```

---

#### 示例二

KB，MB，GB，TB

```go
package main

import "fmt"

const (
	_ = iota //iota = 0
	KB int64 = 1 << (10 *iota) // 左移10位 ，2的10次方等于1024
	MB // iota = 1
	GB // 与KB表达式相同，但 iota = 2
	TB
)

func main(){
	fmt.Println(KB,MB,GB,TB)
}
// output
// 1024 1048576 1073741824 1099511627776
```

