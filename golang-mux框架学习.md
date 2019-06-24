#### grilla/mux

```go
# gorilla/mux 框架包实现了一个请求路由器和调度器，用于将传入请求匹配到他们各自的处理程序。名称mux代表了<HTTP>请求多路复用器；与标准的http.ServeMux一样,mux.Router将传入的请求与注册路由列表进行匹配，并调用匹配URL或其他条件的路由处理程序。主要特点：
	1. 它实现了http.Handler接口，因此它与标准的http.ServeMux兼容；
	2. 请求可以根据URL host,path,path prefix,schemes,header和query values，HTTP 方法或使用用户自定义匹配器进行匹配。
	import gmux "github.com/gorilla/mux"
	gmux := gmux.NewRouter()
	/////////es 代理
	{
		esApi:=c.APIConfig.Routes.Eapi["m8Es"]
		esService :=es.EsService{Config: c,Perfix:esApi.Url}
		gmux.PathPrefix(esApi.Url).HandlerFunc(esService.Handle).Schemes(httpScheme)
	}
```

#### 使用mux框架步骤

##### 1. 安装

```bash
# -u标志指示可以使用网络来更新命名包和他们的依赖。 默认没有-u参数的情况下，get使用网络检出缺少包但不用它来查找现有包的更新。
go get -u github.com/gorilla/mux
```

##### 2. 示例程序

```go
package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"fmt"
)


func main() {
	// 注册一些url paths和handlers
	gmux:=mux.NewRouter()
	// 我们注册三个路由匹配URL paths到handlers处理程序,这相当于http.HandleFunc()的工
	// 作方式:如果传入的请求URL匹配其中一个路径，则相应的处理程序称为传递
	// (http.ResponseWriter, *http.Request)作为参数。也就是HomeHandler其实内部HomeHandler(w http.ResponseWriter, req *http.Request)
	gmux.HandleFunc("/products/{key}",ProductHandle)
	gmux.HandleFunc("/articles/{category}/",ArticlesCategoryHandle)
	gmux.HandleFunc("/articles/{category/{id:[0-9]+}",ArticleHandle)
	http.Handle("/",gmux)
}

func ProductHandle(w http.ResponseWriter, req *http.Request) {
	// ...
}
func ArticlesCategoryHandle(w http.ResponseWriter, req *http.Request) {
	vars:=mux.Vars(req)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,"Category: %v\n", vars["category"])
}
func ArticleHandle(w http.ResponseWriter, req *http.Request){
	// ...
}
```

##### 3.匹配路由

```go
package main

import (
	"github.com/gorilla/mux"
	"net/http"
)


func main() {
	gmux:=mux.NewRouter()
	// 只有域名domain为"www.abc.com.cn"时才匹配.
	gmux.Host("www.abc.com.cn")
	// 动态匹配子域名
	gmux.Host("{subdomain:[a-z]+}.abc.com.cn")
	// 还有如下几个匹配器可以添加:
	// 匹配路径path前缀
	gmux.PathPrefix("/products/")
	// 匹配HTTP 方法
	gmux.Methods("GET","POST")
	// 匹配URL schemes
	gmux.Schemes("http")
	// 匹配header values
	gmux.Headers("X-Requested-With","XMLHttpRequest")
	// 匹配query values
	gmux.Queries("key","value")
	// 用自定义的匹配方法
	gmux.MatcherFunc(func(req *http.Request,rm *mux.RouteMatch) bool{
		return true
	})
	// 常用的是：在一条路线上组合多个匹配器
	gmux.HandleFunc("/products",ProductsHandle).Host("www.abc.com.cn").Methods("GET").Schemes("http")

	// 一次又一次设置相同的匹配条件可能会很无聊，所以我们有一种方法可以将几条共享相同要求的路线分组。 我们称之为“子路由”。
	s:=gmux.Host("www.ccssogt.com.cn").Subrouter()
	s.HandleFunc("/products/",ProductsHandle)
	s.HandleFunc("/products/{key}",ProductHandle)
	s.HandleFunc("/articles/{category}/{id:[0-9]+}",ArticleHandle)
}
```

##### 4. 静态文件服务器

```go
package main

import (
	"flag"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/log"
	"net/http"
	"time"
)

// 静态文件
func main() {
	var dir string
	flag.StringVar(&dir, "dir", ".", "从中提供文件的目录。 默认为当前目录")
	flag.Parse()
	gmux := mux.NewRouter()
	// 这将提供htp://127.0.0.1:8000/static/<filename>下的文件
	gmux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler: gmux,
		Addr:    "127.0.0.1:8000",
		// 为创建的servers设置强制超时时间
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	// 访问：http://127.0.0.1:8000/static/,将会列出你当前项目下的所有文件和目录
	log.Fatalln(srv.ListenAndServe())
}
```

##### 5. 注册URLs

```go

```

