一个最简单的Golang实现的http服务器

```go
package main

import (
	"fmt"
	"net/http"
)

// 所谓的http服务器，主要在于server端如何接受client端的request请求，并向client端返回response
func IndexHandle(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello Golang")
}
func main() {
	http.HandleFunc("/", IndexHandle)
	http.ListenAndServe(":9090", nil)
}
```

​	接收request的过程中，最重要的莫过于路由（`router`），即实现一个`Multiplexer`器。Go中既可以使用内置的mutilplexer --- `DefautServeMux`，也可以自定义。Multiplexer路由的目的就是为了找到处理器函数（`handler`），后者将对request进行处理，同时构建response。

简单总结就是下面这个流程：

```go
Client -> Requests -> Multiplexer(router) -> handler -> Response -> Client

因此，理解go中的http服务，最重要就是要理解Multiplexer和handler,Golang中的Multiplexer基于ServeMux结构，同时也实现了Handler接口。
type ServeMux struct {
	mu    sync.RWMutex
	m     map[string]muxEntry
	hosts bool // whether any patterns contain hostnames
}
```

- handler函数：具有func(w http.ResponseWriter, r *http.Request)签名的函数
- handler处理器(函数):经过HandlerFunc结构包装的handler函数，它实现了ServeHTTP接口方法的函数。调用handler处理器的ServeHTTP方法时，即调用handler函数本身。
- handler对象：实现了Handler接口ServeHTTP方法的结构

```go
handler处理器和handler对象的差别在于，一个是函数，另外一个是结构，它们都有实现了ServeHTTP方法。很多情况下它们的功能类似，下文就使用统称为handler。这算是Golang通过接口实现的类动态类型吧。

go的http服务都是基于handler进行处理。
```

---

#### 创建HTTP服务

```go
创建一个http服务，大致需要经历两个过程，首先需要注册路由，即提供url模式和handler函数的映射；其次就是实例化一个server对象，并开启对客户端的监听。
net/http包暴露的注册路由的api很简单，http.HandleFunc选取了DefaultServeMux作为multiplexer：

func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
    DefaultServeMux.HandleFunc(pattern, handler)
}
那么什么是DefaultServeMux呢？实际上，DefaultServeMux是ServeMux的一个实例。当然http包也提供了NewServeMux方法创建一个ServeMux实例，默认则创建一个DefaultServeMux：

注册好路由之后，启动web服务还需要开启服务器监听。http的ListenAndServer方法中可以看到创建了一个Server对象，并调用了Server对象的同名方法：

监听开启之后，一旦客户端请求到底，go就开启一个协程处理请求，主要逻辑都在serve方法之中。

serve方法比较长，其主要职能就是，创建一个上下文对象，然后调用Listener的Accept方法用来　获取连接数据并使用newConn方法创建连接对象。最后使用goroutein协程的方式处理连接请求。因为每一个连接都开起了一个协程，请求的上下文都不同，同时又保证了go的高并发。serve也是一个长长的方法：
```



 

 

 

 

 