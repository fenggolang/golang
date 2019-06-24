package main

import (
	"io"
	"net/http"
	"os"
)

// go语言使用http下载文件
var (
	url = "http://172.18.1.249:9081/m8cloud/m8-v4.0/demo/bigdata/WordCount.jar"
)

func main() {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	// 指定下载文件的位置,文件夹必须存在
	f, err := os.Create("E:\\tmp\\WordCount.jar")
	if err != nil {
		panic(err)
	}
	io.Copy(f, res.Body) // 把res响应的内容体Body复制到f文件中
}
