```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// 获取request.URL的项目Field信息
func golangRequestURLFieldInfo(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println(req.Form)
	fmt.Println("req.URL.Path", req.URL.Path)         // 路径（相对路径可能省略前导斜杠）
	fmt.Println("req.URL.Host", req.URL.Host)         // host or host:port
	fmt.Println("req.URL.RawQuery", req.URL.RawQuery) //原始查询，编码的查询值，没有'？'
	fmt.Println("req.URL.Scheme", req.URL.Scheme)
	fmt.Println("req.URL.User", req.URL.User)             // 结构体，用户和密码信息
	fmt.Println("req.URL.Fragment", req.URL.Fragment)     //片段的引用，没有'＃'
	fmt.Println("req.URL.ForceQuery", req.URL.ForceQuery) //追加（强制）查询（'？'）即使RawQuery为空
	fmt.Println("req.URL.Opaque", req.URL.Opaque)         // 编码的不透明的数据
	fmt.Println("req.URL.RawPath", req.URL.RawPath)       // 编码路径提示（请参阅EscapedPath方法）
	fmt.Println(req.Form["url_long"])
	for k, v := range req.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	w.Write([]byte("hello world,hello golang http server!"))
}

func main() {
	http.HandleFunc("/", golangRequestURLFieldInfo)
	http.HandlerFunc
	log.Fatal(http.ListenAndServe(":9090", nil))
}

```

