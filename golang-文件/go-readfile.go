package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// 读取文件内容并打印输出
func main() {
	var filePath = "E:\\golang websocket开发.md"
	//contents, err := ioutil.ReadFile(filePath)
	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("文件读取出错")
		return
	}
	// 因为contents是[]byte类型，直接转换成string类型后会多一行空行
	// 需要使用strings.Replace转换去掉
	result := strings.Replace(string(contents), "\n", "", 1)
	fmt.Println(result)
}
