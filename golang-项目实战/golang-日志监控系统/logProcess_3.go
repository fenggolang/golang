package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	url2 "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)



type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(rc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader // 读取器
	write Writer //写入器
}

type ReadFromFile struct {
	path string //读取文件的路径
}

type WriteToInfluxDB struct {
	influxDBDsn string //influxdb数据源
}

type Message struct {
	TimeLocal                    time.Time
	BytesSent                    int
	Path, Method, Scheme, Status string
	UpstreamTime, RequestTime    float64
}

// 读取模块
func (r *ReadFromFile) Read(rc chan []byte) {
	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("打开文件出错:%s", err.Error()))
	}

	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2) // 文件指针移动到最后
	rd := bufio.NewReader(f)

	for {
		//逐行读取
		line, err := rd.ReadBytes('\n')

		if err == io.EOF {
			// 读取到文件末尾，等待文件生成
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("读取字节出错:%s", err.Error()))
		}
		rc <- line[:len(line)-1] // 去掉换行符
	}
}

// 解析模块
func (l *LogProcess) Process() {
	// 正则表达式
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)

	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v)) //匹配数据内容
		if len(ret) != 14 {
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}

		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 -0700", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings Split fail:", ret[6])
			continue
		}
		message.Method = reqSli[0]

		url, err := url2.Parse(reqSli[1])
		if err != nil {
			log.Println("url parse fail:", err.Error())
			continue
		}
		message.Path = url.Path

		// 协议：http
		message.Scheme = ret[5]
		message.Status = ret[7]

		upstreamTime, _ := strconv.ParseFloat(ret[12], 64)
		requestTime, _ := strconv.ParseFloat(ret[13], 64)
		message.UpstreamTime = upstreamTime
		message.RequestTime = requestTime

		l.wc <- message
	}
}

// 写入模块
func (w WriteToInfluxDB) Write(wc chan *Message) {
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
		wc:    make(chan *Message),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(30 * time.Second)
}
