HTTP服务器的使用模式：

Handle处理函数:只要函数的签名为func (w http.ResponseWriter, req *http.Request),均可作为处理函数，即它可以被转换为http.HandlerFunc函数类型。