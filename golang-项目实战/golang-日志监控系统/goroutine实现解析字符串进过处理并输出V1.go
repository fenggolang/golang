package main

import (
	"fmt"
	"strings"
	"time"
)

type LogProcess struct {
	rc          chan string // 读取模块到解析模块传递数据
	wc          chan string //解析模块到写入模块解析数据
	path        string      // 读取文件的路径
	influxDBDsn string      // influxdb data source
}

//读取模块,下面用*是基于性能上的一些考虑，表示引用，好处是1避免了拷贝和可以修改原来的值
func (l *LogProcess) ReadFromFile() {
	line := "message"
	// 把字符串写入rc管道
	l.rc <- line
}

// 解析模板
func (l *LogProcess) Process() {
	// 读取rc管道中的字符串并赋值给data变量
	data := <-l.rc
	// data变量解析并处理得到的字符
	l.wc <- strings.ToUpper(data)
}

// 写入模块
func (l *LogProcess) WriteToInfluxDB() {
	// 把处理后的管道中的数据输出
	fmt.Println(<-l.wc)
}

func main() {
	// 用&表示性能上的一些考虑
	lp := &LogProcess{
		rc:          make(chan string),
		wc:          make(chan string),
		path:        "/tmp/access.log",
		influxDBDsn: "username@password..",
	}
	go lp.ReadFromFile()
	//go (*lp).ReadFromFile() //因为lp是引用所以这才是标准写法，但是这样写不易读，所以golang可以省略掉*
	go lp.Process()
	go lp.WriteToInfluxDB()

	time.Sleep(1 * time.Second) //等待1秒等待goroutine执行完，否则可能无法输出预期的MESSAGE
}

// output:
// MESSAGE