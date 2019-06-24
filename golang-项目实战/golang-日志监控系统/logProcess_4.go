package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	url2 "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)



type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader //读取器
	write Writer //写入器
}

type ReadFromFile struct {
	path string //读取文件的路径
}

type WriteToInfluxDB struct {
	influxDBDsn string // influxdb数据源
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
		panic(fmt.Sprintf("打开文件错误:%s", err.Error()))
	}
	//f.Close()
	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2) //文件指针移动到最后
	rd := bufio.NewReader(f)

	for {
		// 逐行读取
		line, err := rd.ReadBytes('\n')

		if err == io.EOF {
			//读取到文件末尾，等待文件生成
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes 出错：%s", err.Error()))
		}
		rc <- line[:len(line)-1] //去掉换行符
	}
}

// 写入模块
func (w WriteToInfluxDB) Write(wc chan *Message) {
	infSli := strings.Split(w.influxDBDsn, "@")

	// 创建一个influxdb HTTP客户端
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     infSli[0],
		Username: infSli[1],
		Password: infSli[2],
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// 创建一个新的point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  infSli[3],
		Precision: infSli[4],
	})
	if err != nil {
		log.Fatal(err)
	}

	for v := range wc {
		// 创建一个点并添加到批处理中。
		// Tags: Path,Method, Scheme, Status
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status}
		// Fields: UpstreamTime, RequestTime, BytesSent
		fields := map[string]interface{}{
			"UpstreamTime": v.UpstreamTime,
			"RequestTime":  v.RequestTime,
			"BytesSent":    v.BytesSent,
		}

		pt, err := client.NewPoint("nginx_log", tags, fields, v.TimeLocal)
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)

		// 批量写
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// 关闭客户端资源
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("写入数据成功！")
	}
}

// 解析模块
func (l *LogProcess) Process() {
	// 定义正则表达式
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v)) //匹配数据内容
		if len(ret) != 14 {
			log.Println("FindStringSubmatch失败：", string(v))
			continue
		}

		message := &Message{}
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 -0700", ret[4], loc)
		if err != nil {
			log.Println("ParseInLocation 失败：", err.Error(), ret[4])
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			log.Println("strings Split failL:", ret[6])
			continue
		}
		message.Method = reqSli[0]

		url, err := url2.Parse(reqSli[1])
		if err != nil {
			log.Println("url 解析失败：", err.Error())
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

func main() {
	// https://github.com/influxdata/influxdb/tree/master/client

	// influxdb创建数据库:CREATE DATABASE "wpc_test"
	// influxdb创建用户：CREATE USER "wpc" WITH PASSWORD 'wpc'
	var path, influxDsn string
	flag.StringVar(&path, "path", "./access.log", "读取文件路径")
	flag.StringVar(&influxDsn, "influxDsn", "http://172.18.2.200:8086@wpc@wpc@wpc_test@s", "influxdb数据源")
	flag.Parse()

	r := &ReadFromFile{
		//path:"./access.log"
		path: path,
	}
	w := &WriteToInfluxDB{
		influxDBDsn: influxDsn,
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

	// 去influxdb中查询插入的数据：
	// use wpc_test;select * from nginx_log

	time.Sleep(30 * time.Second)
}
