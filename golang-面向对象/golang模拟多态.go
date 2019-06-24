package main

import "fmt"

// golang面向对象-模拟多态(通过interface关键字)
type Foo interface {
	qux()
}

type Bar struct{}
type Baz struct{}

func (b Bar) qux() {
	fmt.Println("Bar qux")
}
func (b Baz) qux() {
	fmt.Println("Baz qux")
}

func main() {
	var f Foo
	f = Bar{}
	fmt.Println("Bar", f)
	f = Baz{}
	fmt.Println("Baz", f)
}
