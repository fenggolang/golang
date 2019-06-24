package main

import (
	"fmt"
	"strings"
	"time"
)



type LogProcess struct {
	rc          chan string
	wc          chan string
	path        string // 读取文件的路径
	influxDBDsn string // influxdb数据源
}

// 读取模块:从文件中读取，其实我们这里简单一点，直接赋值读取字符串而已
func (l *LogProcess) ReadFromFile() {
	line := "message"
	l.rc <- line
}

// 解析模块: 把管道里面的字符串全部解析处理为大写
func (l *LogProcess) Process() {
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

// 写入模块:写入influxdb,其实我们这里处理简单一点，直接输出控制台
func (l *LogProcess) WriteToInfluxDB() {
	fmt.Println(<-l.wc)
}

func main() {
	lp := &LogProcess{
		rc:          make(chan string),
		wc:          make(chan string),
		path:        "/tmp/access.log",
		influxDBDsn: "username&password..",
	}

	go lp.ReadFromFile()
	go lp.Process()
	go lp.WriteToInfluxDB()

	time.Sleep(1 * time.Second)

}

// OUTPUT: MESSAGE