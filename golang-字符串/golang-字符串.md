#### golang 字符串

- rune相当于go的char

```bash
package main

import "fmt"

func main(){
	s := "Yes我爱慕课网!"
	fmt.Println("字符串s长度：",len(s)) //19。中文字符，每个字符占用3字节，UTF-8编码，可变长编码
	fmt.Printf("%s\n", []byte(s)) // Yes我爱慕课网!

	fmt.Printf("%X\n", []byte(s)) // 打印出16进制
}
```

#### 字符串常见操作

- Fields,Split Join
- Contains, Index
- ToLower,ToUpper
- Trim,TrimRight,TrimLeft

---

#### 判断字符串是否是某字符串开头

```go
package main

import (
	"fmt"
	"strings"
)

// 判断是不是以某个字符串开头
func main(){
	str := "hello world"
	res0 := strings.HasPrefix(str,"http://")
	res1:=strings.HasPrefix(str,"hello")

	fmt.Printf("res0 is %v\n",res0) // res0 is false
	fmt.Printf("res1 is %v\n",res1) // res1 is true
}
```

#### 判断字符串是否是某字符串结尾

```go
package main

import (
	"fmt"
	"strings"
)

// 判断是不是以某个字符串结尾
func main(){
	str := "hello world"
	res0 := strings.HasSuffix(str,"http://")
	res1:=strings.HasSuffix(str,"world")

	fmt.Printf("res0 is %v\n",res0) // res0 is false
	fmt.Printf("res1 is %v\n",res1) // res1 is true
}
```

