package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)



// golang实现日志监控系统v2
/**
执行此go代码：在项目当前路径下创建access.log文件并写入如下日志格式：
100.97.120.0 - - [04/Mar/2018:13:49:52 +0800] http "GET /foo?query=t HTTP/1.0" 200 612 "-" "KeepAliveClient" "-" 1.005 1.854
当程序运行起来执行，修改access.log日志文件，在文件末尾新增多行，控制台会检测到文件新内容生成并输出。
*/
type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader //读取器
	write Writer //写入器
}

type ReadFromFile struct {
	path string //读取文件的路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influxdb数据源
}

// 读取模块
func (r *ReadFromFile) Read(rc chan []byte) {
	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("打开文件出错:%s", err.Error()))
	}

	// 从文件末尾逐行读取文件内容
	f.Seek(0, 2) //文件指针移动到最后
	rd := bufio.NewReader(f)

	for {
		// 逐行读取
		line, err := rd.ReadBytes('\n')

		if err == io.EOF {
			// 读取到文件末尾，等待文件生成
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("读取字节出错:%s", err.Error()))
		}
		rc <- line[:len(line)-1] //去掉换行符
	}
}

// 解析模块：暂时简单处理解析
func (l *LogProcess) Process() {
	for v := range l.rc {
		l.wc <- strings.ToUpper(string(v))
	}
}

// 写入模块:暂未实现，只是输出
func (w WriteToInfluxDB) Write(wc chan string) {
	for v := range wc {
		fmt.Println(v)
	}
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}
	w := &WriteToInfluxDB{
		influxDBDsn: "username&password..",
	}

	lp := &LogProcess{
		rc:    make(chan []byte),
		wc:    make(chan string),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}
