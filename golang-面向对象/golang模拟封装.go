package main

import "fmt"

// Golang中的面向对象-模拟封装
type Foo struct {
	baz string
}

func (f *Foo) echo() {
	fmt.Println(f.baz)
}

func main() {
	f := Foo{"hello"}
	f.echo()
}
