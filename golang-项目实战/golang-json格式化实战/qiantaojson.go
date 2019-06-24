package main

import (
	"encoding/json"
	"fmt"
)

// 缩进对于嵌套结构也同样使用，如下例中，为Person结构体添加Job嵌套结构之后
// json输出的可读性依然不错。
type Job struct {
	Location string
	Type     string
}
type Person struct {
	Name string
	ID   uint64
	Age  uint8
	Job  *Job
}

func main() {
	p := &Person{
		Name: "wangpengcheng",
		ID:   888888,
		Age:  25,
		Job: &Job{
			Location: "四川省成都市高新区",
			Type:     "golang programmer",
		},
	}

	b, _ := json.MarshalIndent(p, "", "\t")
	fmt.Println(string(b))
}

// output
/**
{
	"Name": "wangpengcheng",
	"ID": 888888,
	"Age": 25,
	"Job": {
		"Location": "四川省成都市高新区",
		"Type": "golang programmer"
	}
}
*/
