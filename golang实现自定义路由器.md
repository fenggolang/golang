```go
package main

import (
	"fmt"
	"net/http"
)

/**
Go支持外部实现路由器，ListenAndServe的第二个参数就是配置外部路由器(如果第二个参数为nil表示使用内置默认路由器)，
它是一个Handler接口。即外部路由器实现Handler接口。
*/
/*// Handler接口
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}*/

type MyMux struct {
}

func sayHelloName(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintln(w, "Hello myroute")
}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		sayHelloName(w, req)
		return
	}
	http.NotFound(w, req) //未路由成功
	return
}

func main() {
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)
}
```

