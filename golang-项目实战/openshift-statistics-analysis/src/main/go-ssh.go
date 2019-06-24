package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
)

/**
模拟500用户拉去和推送镜像

用法：编译生成二进制文件放到linux上，启动的时候带-p参数指定端口，然后在windows上用postman发送命令连接到此服务上(ip:port)
*/

var (
	flagConf = flag.String("p", "7777", "指定端口")
)

func main() {
	flag.Parse()
	http.HandleFunc("/ssh", ssh)

	port := ":" + *flagConf
	fmt.Println("开启服务", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func ssh(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		d, _ := json.Marshal(err)
		w.Write(d)
		return
	}

	fmt.Println("执行：")
	fmt.Println(string(data))

	cmd := exec.Command("/bin/sh", "-c", string(data))

	// 获取输出对象，可以从该对象中读取输出结果
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	//保证关闭输出流
	defer stdout.Close()
	defer stderr.Close()

	//运行命令
	// 注意：Start执行不会等待命令完成就，Run会阻塞等待命令完成。
	if err := cmd.Start(); err != nil {
		d, _ := json.Marshal(err)
		w.Write(d)
		return
	}

	//读取输出结果
	errBytes, _ := ioutil.ReadAll(stderr)
	if errBytes != nil {
		w.Write(errBytes)
	}
	succBytes, _ := ioutil.ReadAll(stdout)
	if succBytes != nil {
		w.Write(succBytes)
	}
	return
}
