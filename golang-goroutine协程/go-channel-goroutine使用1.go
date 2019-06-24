package main

import "fmt"

func sendData(ch chan<- string) {
	ch <- "go"
	ch <- "java"
	ch <- "c"
	ch <- "c++"
	ch <- "python"
	close(ch)
}

func getData(ch <-chan string, chClose chan bool) {
	for {
		str, ok := <-ch
		if !ok {
			fmt.Println("chan is close.")
			break
		}
		fmt.Println(str)
	}
	chClose <- true
}
func main() {
	ch := make(chan string, 10)
	chClose := make(chan bool, 1)
	go sendData(ch)
	go getData(ch, chClose)
	<-chClose
	close(chClose)
}

/**
go
java
c
c++
python
 */