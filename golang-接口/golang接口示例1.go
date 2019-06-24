package main

import "fmt"

// interface变量存储的是实现者的值
type I interface {
	Get() int
	Set(int)
}

type S struct {
	Age int
}

func (s S) Get() int {
	return s.Age
}

func (s *S) Set(age int) {
	s.Age = age
}

func f(i I) {
	i.Set(10)
	fmt.Println(i.Get())
}

func main() {
	s := S{}
	f(&s)

	var i I //声明i
	i = &s  //赋值s到i
	fmt.Println(i.Get())

	// 如何判断interface变量存储的是哪种类型
	if t, ok := i.(*S); ok {
		fmt.Println("s implements I", t)
	}

	// 如果需要区分多种类型，可以使用switch断言，更简单直接，这种断言方式只能在switch语句中使用。
	switch t := i.(type) {
	case *S:
		fmt.Println("i store *S", t)
	}
}
