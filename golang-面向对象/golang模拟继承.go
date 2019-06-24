package main

import "fmt"

// golang面向对象-模拟继承
type Foo struct {
	baz string
}
type Bar struct {
	Foo
}

func (f *Foo) echo(){
	fmt.Println(f.baz)
}

func main(){
	b:=Bar{Foo{"hehh"}}
	b.echo()
}