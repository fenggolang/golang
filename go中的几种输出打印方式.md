- fmt.Printf
- fmt.Println
- fmt.Print


```go
// fmt.Printf指定格式化格式并返回写入到标准输出，返回写入的字节数和error
// fmt.Println使用默认数值类型的格式化写入到标准输出，返回写入的字节数和error
// fmt.Print同上，只是不换行

type point struct {
	x,y int
}

p:=point{1,2}
fmt.Printf("%v\n",p) // 输出：{1 2}
```

- fmt.Sprinf
- fmt.Sprintln
- fmt.Sprint


```go
// fmt.Sprintf指定格式化格式返回格式化后的字符串
// fmt.Sprintln使用默认格式化格式返回格式化后的字符串
// fmt.Sprint同上,只是不换行

// `Printf`函数的输出是输出到命令行`os.Stdout`的，你
// 可以用`Sprintf`来将格式化后的字符串赋值给一个变量
s:=fmt.Sprintf("%s","string") // string
fmt.Println(s)
```

- fmt.Fprintf
- fmt.Fprintln
- fmt.Fprint


```go
// fmt.Fprintf 自己指定输出到哪里，比如os.Stdout,os.Stderr,然后指定格式化格式，然后返回写入的字节数和error
// fmt.Fprintln 自己指定输出到哪里，比如os.Stdout,os.Stderr,使用默认格式化格式，然后返回写入的字节数和error
// fmt.Fprint 同上，只是不换行
fmt.Fprintf(os.Stderr,"an %s\n", "error")
```

- fmt.Errorf


```go

```

