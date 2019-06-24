Go语言有四种特殊的类型通常是看作引用传递

- 切片(Slice)
- 接口(Interface)
- Map
- Channel

```
var arr1 = [...]int{1, 5, 2, 3, 7}是数组定义，不是切片
```

​	数组转切片方式

```bash
sli := arr[:]
```

