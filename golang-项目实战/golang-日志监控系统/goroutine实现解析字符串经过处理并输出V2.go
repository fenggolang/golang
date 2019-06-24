package main

import (
	"fmt"
	"strings"
	"time"
)

// 定义接口提升扩展性
type Reader interface {
	Read(rc chan string)
}
type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc    chan string // 读取模块到解析模块传递数据
	wc    chan string //解析模块到写入模块解析数据
	read  Reader
	write Writer
}

type ReadFromFile struct {
	path string // 读取文件的路径

}
type WriteToInfluxDB struct {
	influxDBDsn string // influxdb data source
}

func (r *ReadFromFile) Read(rc chan string) {
	line := "message"
	// 把字符串写入rc管道 rc <- line
	rc <- line
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	// 把处理后的管道中的数据输出
	fmt.Println(<-wc)
}

// 解析模板
func (l *LogProcess) Process() {
	// 读取rc管道中的字符串并赋值给data变量
	data := <-l.rc
	// data变量解析并处理得到的字符
	l.wc <- strings.ToUpper(data)
}

func main() {
	// 实例化读取器
	r := &ReadFromFile{
		path: "/tmp/access.log",
	}
	// 实例化写入器
	w := &WriteToInfluxDB{
		influxDBDsn: "username&password..",
	}
	// 用&表示性能上的一些考虑
	lp := &LogProcess{
		rc: make(chan string),
		wc: make(chan string),
		// 把读取和写入模块注入到Log结构体中
		read:  r,
		write: w,
	}
	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(1 * time.Second) //等待1秒等待goroutine执行完，否则可能无法输出预期的MESSAGE
}
