#### golang数组

​	定义数组的常用方式：

- var arr1 [5]int

- arr2:=[3]int{1,3,5}

- arr3:=[...]int{2,4,6,8,10}

- var grid [4][5]bool

- [10]int和[20]int是不同的类型

- 调用func f(arr [10]int)会拷贝数组，go语言的参数传递只有一种传递方式就是**值传递**,在其他语言中传递一个数组作为函数参数都是引用传递，函数内改变数组相关值，会直接改变原始的值，而go就不会。

- 在go语言中一般不直接使用数组，go语言中使用更多的是数组切片。

  值得注意的是数组长度的数量是写在类型的前面，而[...]是让编译器帮我们自己数一下有几个数组元素。多维数组方括号也写在前面。

```go
package main

import "fmt"

func main(){
	// 定长未初始化的数组
	var arr1 [5]int
	// 定长并初始化的数组
	arr2:=[3]int{1,3,5}
	// 特殊数组：切片
	arr3:=[...]int{2,4,6,8,10}
	// 定义二维数组，下面表示2个长度为3的bool类型的数组
	var grid [2][3]bool

	fmt.Println(arr1,arr2,arr3)
	fmt.Println(grid)
}
// output:
[0 0 0 0 0] [1 3 5] [2 4 6 8 10]
[[false false false] [false false false]]
```

---

#### golang遍历数组的方式

- 普通for循环方式
- range关键字遍历数组
- 可通过_下划线省略变量
- 不仅range,任何地方都可以通过下划线省略变量
- 如果只要下标i,直接写：for i :=range numbers即可

```go
package main

import "fmt"

func main(){

	// 特殊数组：切片
	arr3:=[...]int{2,4,6,8,10}

	// 遍历数组方式一
	for i:=0;i<len(arr3);i++{
		fmt.Println(arr3[i])
	}
	// 遍历数组方式二
	for i:=range arr3{
		fmt.Println(arr3[i])
	}
	// 遍历数组方式二.2
	for i,v:=range arr3{
		fmt.Println("arr[",i,"]","=",v)
	}
	// 遍历数组方式二.3
	for _,v:=range arr3{
		fmt.Println(v)
	}
}
```

**求一个数组中最大值且最大值下标**

```go
package main

import "fmt"

func main(){
	numbers:=[5]int{1,2,3,4,5}
	maxi:=-1
	maxValue:=-1

	for i,v:=range numbers {
		if v>maxValue{
			maxi,maxValue=i,v
		}
	}
	fmt.Println("maxi是",maxi,"maxValue是：",maxValue)
}
```

**求一个数组的元素的和**

```go
package main

import "fmt"

func main(){
	numbers:=[5]int{1,2,3,4,5}
	sum:=0

	for _,v:=range numbers {
		sum+=v
	}
	fmt.Println("数组的和是：",sum)
}
```

---

#### golang中的range

- range意义明确，美观
- c++:没有类似的能力，只能通过for a< len(arr)方式
- java/Python:只能for each value,不能同时获取i,v

```go
package main

import "fmt"

// go中[3]int和[5]int是表示不同的类型

func printArray(arr [5]int){
	//只能在函数内可变，函数外不行,表明数组是值类型，想改变的话则指针传递即可
	arr[0]=100
	for i,v:= range arr{
		fmt.Println(i,v)
	}
}

func main(){
	var arr1  [5]int
	//arr2:=[3]int{1,3,5}
	arr3:=[...]int{2,4,6,8,10}

	//var grid [2][3]int

	printArray(arr1)
	//printArray(arr2) 报错，因为元素只有3个，不是5个
	printArray(arr3)
	//printArray(grid) 同arr2
}
```

```go
package main

import "fmt"

// go中[3]int和[5]int是表示不同的类型

func printArray(arr *[5]int){
  //因为是指针传递，所以会改变数组的值在原来的地方(不仅在函数内改变),当然go的指针传递之后修改数组的值，照样写arr[0]=xxx即可，不用前面加什么*(xxx)
	arr[0]=100
	for i,v:= range arr{
		fmt.Println(i,v)
	}
}
func main(){
	var arr1  [5]int
	//arr2:=[3]int{1,3,5}
	arr3:=[...]int{2,4,6,8,10}

	//var grid [2][3]int

	printArray(&arr1)
	fmt.Println("arr1[0]值已经改变",arr1[0])
	//printArray(arr2) 报错，因为元素只有3个，不是5个
	printArray(&arr3)
	fmt.Println("arr3[0]值已经改变",arr3[0])
	//printArray(grid) 同arr2
}
```

