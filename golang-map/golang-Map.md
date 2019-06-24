#### golang Map

map定义：map[K]V,map[K1]map[K2]V

map定义方式：

- 直接赋值
- make创建

```go
package main
func main(){
	m:=map[string]string{
		"name":"ccmouse",
		"course":"golang",
		"site":"imooc",
		"quality":"notbad",
	}
```

```go
package main
import "fmt"
func main(){
	m2:=make(map[string]int)
	fmt.Println(m2)
}
```

---

#### map操作总结

- 创建：make(map[string]int)
- 获取元素：m[key]
- key不存在时，获取Value类型的初始值，也就是不会报错
- 用value,ok:=m[key]来判断是否存在key
- 用delete删除一个key
- 使用range 遍历key,或者遍历key.value对
- map是一个无序的字典，不保证遍历顺序，如需顺序，需手动对key排序。
- 使用len获得元素个数
- map使用hash表，必须可以比较相等
- 除了slice,map,function的内建类型都可以作为key
- struct类型不包含上述字段，也可以作为key

---

#### golang map实战

​	寻找最长不含有重复字符的字串

abcabcbb --> abc

bbbbb --> b

pwwkew --> wke

```go
package main

import (
	"fmt"
)

func lengthOfNonRepeatingSubStr(s string) int {
	lastOccurred := make(map[rune]int)
	start := 0
	maxLength := 0

	for i, ch := range []rune(s) {
		if lastI, ok := lastOccurred[ch]; ok && lastI >= start {
			start = lastI + 1
		}
		if i-start+1 > maxLength {
			maxLength = i - start + 1
		}
		lastOccurred[ch] = i
	}

	return maxLength
}

func main() {
	fmt.Println(
		lengthOfNonRepeatingSubStr("abcabcbb"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("bbbbb"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("pwwkew"))
	fmt.Println(
		lengthOfNonRepeatingSubStr(""))
	fmt.Println(
		lengthOfNonRepeatingSubStr("b"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("abcdef"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("这里是慕课网"))
	fmt.Println(
		lengthOfNonRepeatingSubStr("一二三二一"))
	fmt.Println(
		lengthOfNonRepeatingSubStr(
			"黑化肥挥发发灰会花飞灰化肥挥发发黑会飞花"))
}

```

