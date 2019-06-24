#### golang切片

```go
package main

import "fmt"

//函数参数是数组名称,数组长度必须是8，因为arr原始长度也是8.否则go语言识别为只要数组长度不相同的话都认为是不同的类型
func udateArrayByArrayName(arr [8]int){
	arr[0]=100
}

// 函数参数值指针,数组长度必须是8，因为arr原始长度也是8.否则go语言识别为只要数组长度不相同的话都认为是不同的类型
func updateArrByPointer(arr *[8]int){
	arr[0] = 101
}

// 数组中括号中不加长度就是一个slice,切片.
func updateSlice(s []int) {
	s[0] = 102
}

func main(){
	// 定义一个切片
	arr:=[...]int{0,1,2,3,4,5,6,7}
	// 取出2到6，前开后闭
	s1:=arr[2:6]
	s2:=arr[2:]
	s3:=arr[:6]
	s4:=arr[:]
	fmt.Println("arr[2:6]=",s1) // output:arr[2:6]= [2 3 4 5]
	fmt.Println("arr[2:]=",s2) // arr[2:]= [2 3 4 5 6 7]
	fmt.Println("arr[:6]=",s3) // arr[:6]= [0 1 2 3 4 5]
	fmt.Println("arr[:]=",s4) // arr[:]= [0 1 2 3 4 5 6 7]

	// 我们知道在go语言中，只有值传递，包括数组也是(其他语言中数组时引用传递)，如果想要修改原始的值，则需要通过指针方式作为函数参数
	fmt.Println("原来的arr数组的值是：",arr)

	// 测试传递函数参数为普通的数组名称
	udateArrayByArrayName(arr)
	fmt.Println("传递的函数参数为数组名称之后数组的值是：",arr) // 没变

	// 测试传递函数参数为指针
	updateArrByPointer(&arr)
	fmt.Println("传递的函数参数为数组指针之后数组的值是：",arr) // 指针传递会修改原始数组的值

	// 测试传递函数参数为切片
	updateSlice(arr[:])
	fmt.Println("传递的函数参数为切片之后数组的值是：",arr)// 切片传递会修改原始数组的值

	// 切片再切片
	s:=arr[:5]
	fmt.Println("再切片s:",s) //[102 1 2 3 4]
	s=s[2:]
	fmt.Println("再再切片ss:",s) //[2 3 4]
}

// 切片总结：
/**
Slice切片本身没有数据，是对底层array数组的一个view,
,一个slice切片可以继续slice切片，俗称重新slice,这些重切片都是view同一个数组切片
 */
```

#### golang切片扩展

```go
package main

import "fmt"

func main(){
	arr:=[...]int{0,1,2,3,4,5,6,7}
	// s1取arr切片下标为2,3,4,5的元素
	s1:=arr[2:6]
	// s2取arr切片下标为3,4的元素
	s2:=s1[3:5] //取s1[3],s1[4].s1[4]根本就不在s1数组中
	fmt.Println("s1:=",s1) //[2,3,4,5]
	fmt.Println("s2:=",s2) // [5,6] //即使切片长度溢出，只要不溢出最底层arr切片的长度都回取出来

	// slice切片里面装了三个变量：ptr指针，len长度，cap容量
	// slice切片可以向后扩展，但是不可以向前扩展，也就是向后的隐藏元素可以看到，向前s1隐藏的元素看不到。
	// s[i]不可以超越len(s),向后扩展不可以超越底层数组cap(s)

	// 格式化输出
	fmt.Printf("s1=%v, len(s1)=%d, cap(s1)=%d\n",s1,len(s1),cap(s1)) // s1=[2 3 4 5], len(s1)=4, cap(s1)=6
	fmt.Printf("s2=%v, len(s2)=%d, cap(s2)=%d\n",s2,len(s2),cap(s2)) // s2=[5 6], len(s2)=2, cap(s2)=3

	//fmt.Println(s1[3:7]) 这个不行，因为超过arr的总容量了。
}
```

---

#### 向切片中添加元素

```go
package main

import "fmt"

func main(){
	arr:=[...]int{0,1,2,3,4,5,6,7}
	s1:=arr[2:6]
	s2:=s1[3:5]
	s3:=append(s2,10)
	s4:=append(s3,11)
	s5:=append(s4,12)

	fmt.Println("s1=",s1) // [2,3,4,5]
	fmt.Println("s2=",s2) // [5,6]
	fmt.Println("s3=",s3) // [5,10]
	fmt.Println("s4=",s4) // [5,10,11]
	fmt.Println("s5=",s5) // [5,10,11,12]
	fmt.Println("arr=",arr) // [{0,1,2,3,4,5,6,10]
}

// 结论：向slice添加元素的时候，如果超越底层切片容量cap,系统会重新分配更大的底层数组。
//      由于值传递的关系，必须接受app的返回值。
// go语言没有null,他用nil表示没有
```

---

#### 切片常见操作

- 创建切片
- 向切片中添加元素
- 拷贝切片
- 删除切片中元素

```go
package main

import "fmt"

// 打印切片长度和容量
/**
len=0, cap=0
len=1, cap=1
len=2, cap=2
len=3, cap=4
len=4, cap=4
len=5, cap=8
len=6, cap=8
len=7, cap=8
len=8, cap=8
len=9, cap=16
 */
func printSlice(s []int){
   fmt.Printf("s=%v,len=%d, cap=%d\n",s,len(s),cap(s))
}
func main() {
   var s []int // Zero value for slice is nil

   for i:=0;i<10;i++{
      printSlice(s)
      s = append(s,2*i+1)
   }
   fmt.Println(s) // [1 3 5 7 9 11 13 15 17 19]

   s1:=[]int{2,4,6,8}
   printSlice(s1) //len=4, cap=4

   s2:=make([]int,16)
   printSlice(s2) // len=16, cap=16

   s3:=make([]int,10,32)
   printSlice(s3) // len=10, cap=32

   // 拷贝切片把s1的元素拷贝到s2切片中
   fmt.Println("Coping slice")
   copy(s2,s1)
   printSlice(s2) // len=16, cap=16

   //删除切片中的指定一个元素,并且后面的元素前移一位
   fmt.Println("Deleting elements from slice")
   s2 = append(s2[:3],s2[4:]...)
   // 删除切片中的一个元素不会缩小切片容量大写，但是其len长度会缩小1
   printSlice(s2)

   // 删除元素头
   fmt.Println("Poping from front")
   front:=s2[0]
   s2=s2[1:]
   // 删除元素尾
   fmt.Println("Popping from back")
   tail:=s2[len(s2)-1]
   s2=s2[:len(s2)-1]
   // 删除的头和尾的元素是：
   fmt.Println(front,tail)
   // 删除头和尾元素纸盒s2的值是：
   printSlice(s2)
}


```

​	总结，创建Slice切片的几种方式

- var s []int 
- s2:=make([]int,16)
- s3:=make([]int,10,32)
- 把其他切片赋值得到新切片
- 切片再自身切片


slice类似于python中的切片。可以理解为一个指向数组的指针。而且这个指针还有两个域：长度，容量。长度为当前元素个数，容量为最大元素个数，类似于STL中vector的size和capacity。正常情况下slice和数组共享一段存储空间，当我们对slice操作使得slice超出数组大小时，则slice会重新分配空间。也可以用make为slice分配空间。

```go
package main

import "fmt"



func main() {

	a := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	slice1 := a[5:10]
	fmt.Println(slice1) //[6,7,8,9,10]

	a[5] = 10
	fmt.Println(slice1) //[10,7,8,9,10]

	slice1 = append(slice1, 100) //超出了原始数组的容量，会重新建一个数组，不会再指向原来的数组了
	fmt.Println(slice1)          // [10 7 8 9 10 100]

	a[6] = 11           //
	fmt.Println(slice1) // [10 7 8 9 10 100]

	// map
	mymap := make(map[int]int)
	mymap[10] = 10 * 10
	mymap[5] = 5 * 5
	mymap[15] = 15 * 15

	for k, v := range mymap {
		fmt.Println(k, v)
	}
}
```

