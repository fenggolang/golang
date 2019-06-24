package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)



var (
	//自定义数值形态的测量数据:瞬时值
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "CPU 目前的温度",
	})
	// 计数形态的测量数据，并带有自定义标签：累计值
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "硬盘发生错误的次数",
		},
		[]string{"device"},
	)
)

// 初始化函数
func init() {
	//测量数据必须注册才会暴露给外界知道：
	prometheus.MustRegister(cpuTemp)
	prometheus.MustRegister(hdFailures)
}

func main() {
	// 配置测量数据的数值
	cpuTemp.Set(65.3)
	hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// 我们会用prometheus所提供的预设处理函数在"/metrics"路径监控着。
	// 这会暴露我们的数据内容,所以prometheus就能够获取到这些数据
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// 浏览器页面访问：127.0.0.1:8080/metrics
// ref:https://yami.io/golang-prometheus/