package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)


/**
 监控模块的实现
1. 总处理日志行数
2. 系统吞吐量
3. read channel; 长度
4. write channel 长度
5. 运行总时间
6. 错误数
*/
type Reader interface {
	Read(rc chan []byte)
}
type Writer interface {
	Write(wc chan *Message)
}

type LogProcess struct {
	rc    chan []byte
	wc    chan *Message
	read  Reader // 读取器
	write Writer // 写入器
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

const (
	TypeHandleLine = 0
	TypeErrNum     = 1
)

var TypeMonitorChan = make(chan int, 200)

// 系统状态监控：定义结构体存储监控数据,然后下面通过HTTP接口的形式将数据暴露出去
type SystemInfo struct {
	HandleLine   int     `json:"handleLine"`   // 总处理日志行数
	Tps          float64 `json:"tps"`          // 系统吞吐量
	ReadChanLen  int     `json:"readChanLen"`  // read channel 长度
	WriteChanLen int     `json:"writeChanLen"` // write channel长度
	RunTime      string  `json:"runTime"`      // 运行总时间
	ErrNum       int     `json:"errNum"`       // 错误数
}

// 定义一个Monitor作为监控模块的封装
type Monitor struct {
	startTime time.Time //监控开始时间
	data      SystemInfo
	tpsSli    []int
}

func (m *Monitor) start(lp *LogProcess) {
	// 定义一个协程：判断是ErrNum加1还是HandleLine加1，避免同时加1
	go func() {
		for n := range TypeMonitorChan {
			switch n {
			case TypeErrNum:
				m.data.ErrNum += 1
			case TypeHandleLine:
				m.data.HandleLine += 1
			}
		}
	}()

	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-ticker.C
			m.tpsSli = append(m.tpsSli, m.data.HandleLine)
			if len(m.tpsSli) > 2 {
				m.tpsSli = m.tpsSli[1:]
			}
		}
	}()

	// 请求路径为：/monitor
	http.HandleFunc("/monitor", func(writer http.ResponseWriter, request *http.Request) {
		m.data.RunTime = time.Now().Sub(m.startTime).String()
		m.data.ReadChanLen = len(lp.rc)
		m.data.WriteChanLen = len(lp.wc)

		if len(m.tpsSli) >= 2 {
			m.data.Tps = float64(m.tpsSli[1]-m.tpsSli[0]) / 5
		}

		ret, _ := json.MarshalIndent(m.data, "", "\t")

		// 将内容输出到http writer里面
		io.WriteString(writer, string(ret))

	})
	// 监听9193端口,ListenAndServe方法是阻塞的方法，所以正常情况下程序不会运行中断
	http.ListenAndServe(":9193", nil)
}

// 读取模块
func (r *ReadFromFile) Read(rc chan []byte) {
	// 打开文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("打开文件出错：%s", err.Error()))
	}

	// 从文件末尾开始逐行读取文件内容
	f.Seek(0, 2) //文件指针移动到最后
	rd := bufio.NewReader(f)

	for {
		//逐行读取
		line, err := rd.ReadBytes('\n')

		if err == io.EOF {
			// 读取到文件末尾，等待文件生成
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			panic(fmt.Sprintf("ReadBytes出错:%s", err.Error()))
		}
		// 处理行数写入chan
		TypeMonitorChan <- TypeHandleLine
		rc <- line[:len(line)-1] //去掉换行符
	}
}

// 写入模块
func (w WriteToInfluxDB) Write(wc chan *Message) {
	// 分割influxdb数据源，解析出每一部分
	infSli := strings.Split(w.influxDBDsn, "@")

	// 创建一个HTTPClient
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
		// 创建一个point监控点和添加到批处理中
		// Tags: Path,Method,Scheme,Status
		tags := map[string]string{"Path": v.Path, "Method": v.Method, "Scheme": v.Scheme, "Status": v.Status}
		// Fields: UpstreamTime,RequestTime,BytesSent
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

		// 批量写入
		if err := c.Write(bp); err != nil {
			log.Fatal(err)
		}

		// 关闭客户端资源
		if err := c.Close(); err != nil {
			log.Fatal(err)
		}

		log.Println("数据写入influxdb成功!")
	}
}

// 解析模块
func (l *LogProcess) Process() {
	//100.97.120.0 - - [08/Jan/2016:10:40:18 +0800] http "GET /foo?query=t HTTP/1.0" 200 612 "-" "KeepAliveClient" "-" 1.005 1.854
	//100.97.120.0 - - [04/Mar/2018:13:49:52 +0000] http "GET /foo?query=t HTTP/1.0" 200 612 "-" "KeepAliveClient" "-" 1.005 1.854

	//([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)
	r := regexp.MustCompile(`([\d\.]+)\s+([^ \[]+)\s+([^ \[]+)\s+\[([^\]]+)\]\s+([a-z]+)\s+\"([^"]+)\"\s(\d{3})\s+(\d+)\s+\"([^"]+)\"\s+\"(.*?)\"([\d\.-]+)\"\s+([\d\.-]+)\s+([\d\.-]+)`)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	for v := range l.rc {
		ret := r.FindStringSubmatch(string(v)) //匹配数据内容
		if len(ret) != 14 {
			TypeMonitorChan <- TypeErrNum
			fmt.Println(ret[:])
			log.Println("FindStringSubmatch fail:", string(v))
			continue
		}

		message := &Message{}
		//log.Println(ret[4])
		t, err := time.ParseInLocation("02/Jan/2006:15:04:05 -0700", ret[4], loc)
		if err != nil {
			TypeMonitorChan <- TypeErrNum
			log.Println("ParseInLocation fail:", err.Error(), ret[4])
			continue
		}
		message.TimeLocal = t

		byteSent, _ := strconv.Atoi(ret[8])
		message.BytesSent = byteSent

		// GET /foo?query=t HTTP/1.0
		reqSli := strings.Split(ret[6], " ")
		if len(reqSli) != 3 {
			TypeMonitorChan <- TypeErrNum
			log.Println("strings Split fail:", ret[6])
			continue
		}
		message.Method = reqSli[0]

		url, err := url2.Parse(reqSli[1])
		if err != nil {
			TypeMonitorChan <- TypeErrNum
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

func main() {
	// https://github.com/influxdata/influxdb/tree/master/client

	var path, influxDsn string
	flag.StringVar(&path, "path", "./access.log", "填写需要读取的文件路径")
	flag.StringVar(&influxDsn, "influxDsn", "http://172.18.2.200:8086@wpc@wpc@wpc_test@s", "influxdb数据源")
	flag.Parse()

	r := &ReadFromFile{
		//path: "./access.log",
		path: path,
	}

	w := &WriteToInfluxDB{
		//influxDBDsn: "http://172.18.2.200:8086@wpc@wpc@wpc_test@s",
		influxDBDsn: influxDsn,
	}

	lp := &LogProcess{
		// 一般情况下，解析模块比读取模块要慢，容易造成阻塞状态，所以
		// 我们限制下rc和wc的chan大小为200，这样相当于这2个chan都有了缓存
		rc:    make(chan []byte, 200),
		wc:    make(chan *Message, 200),
		read:  r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	for i := 0; i < 2; i++ {
		// 解析处理模块协程一般比上面的Read读取模块要慢,我们开2个goroutine去执行
		go lp.Process()
	}
	for i := 0; i < 4; i++ {
		// Write模块最慢,我们开4个goroutine去执行
		go lp.write.Write(lp.wc)
	}

	m := &Monitor{
		startTime: time.Now(),
		data:      SystemInfo{},
	}
	m.start(lp)

	// 监控信息查询方式：
	// curl 127.0.0.1:9193/monitor
	// windows下编译成linux 二进制文件
	/**
	set GOARCH=amd64
	set GOOS=linux
	go build -o logProcess logProcess.go
	*/
	// [root@iot ~]# ./logProcess -path=/root/access.log -influxDsn=http://172.18.2.200:8086@wpc@wpc@wpc_test@s
}
