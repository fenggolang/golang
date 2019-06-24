1、首先调用Http.HandleFunc，按如下顺序执行：

1. 调用了DefaultServerMux的HandleFunc。
2. 调用了DefaultServerMux的Handle。
3. 往DefaultServerMux的map[string] muxEntry中增加对应的handler和路由规则。

2、调用http.ListenAndServe(":9090",nil)，按如下顺序执行：

1. 实例化Server。
2. 调用Server的ListenAndServe()。
3. 调用net.Listen("tcp",addr)监听端口。
4. 启动一个for循环，在循环体中Accept请求。
5. 对每个请求实例化一个Conn，并且开启一个goroutine为这个请求进行服务go c.serve()。
6. 读取每个请求的内容w,err:=c.readRequest()。
7. 判断handler是否为空，如果没有设置handler，handler默认设置为DefaultServeMux。
8. 调用handler的ServeHttp。
9. 根据request选择handler，并且进入到这个handler的ServeHTTP,
   mux.handler(r).ServeHTTP(w,r)
10. 选择handler

- 判断是否有路由能满足这个request（循环遍历ServeMux的muxEntry）。
- 如果有路由满足，调用这个路由handler的ServeHttp。
- 如果没有路由满足，调用NotFoundHandler的ServeHttp。

Go支持外部实现路由器，ListenAndServe的第二个参数就是配置外部路由器(如果是nil则使用内置默认路由器)，它是一个Handler接口。即外部路由器实现Hanlder接口。

 