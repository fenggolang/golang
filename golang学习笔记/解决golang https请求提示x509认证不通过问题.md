参考： [解决golang https请求提示x509: certificate signed by unknown authority](https://blog.bbzhh.com/index.php/archives/150.html)

代码部分

```go
	// 定义TLSClientConfig,跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

// 修改

client := &http.Client{}
// 为

client := &http.Client{Transport: tr}
```

