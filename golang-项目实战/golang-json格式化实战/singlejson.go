package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string
	ID   uint64
	Age  uint8
}

func main() {
	p := &Person{
		Name: "Peter",
		ID:   1024,
		Age:  30,
	}

	b, _ := json.MarshalIndent(p, "", "\t")
	/*
		{
			"Name": "Peter",
			"ID": 1024,
			"Age": 30
		}
	*/

	//b, _ := json.Marshal(p) // {"Name":"Peter","ID":1024,"Age":30}
	fmt.Println(string(b))
}
