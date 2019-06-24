```go
package main

import (
   "flag"
   "log"
   "net/http"

   "github.com/prometheus/client_golang/prometheus/promhttp"
)

// 一个最小的示例：如何用golang操作prometheus
// 本地访问localhost:8080，将会返回prometheus自身监控指标

var addr = flag.String("listen-address", ":8080", "这个地址去监听http请求")

func main() {
   flag.Parse()
   http.Handle("/metrics", promhttp.Handler())
   log.Fatal(http.ListenAndServe(*addr, nil))
}

// url:localhost:8080
```