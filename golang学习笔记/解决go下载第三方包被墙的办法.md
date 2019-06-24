如果go get -u golang.org/x/text无法下载那么就变为下面方式下载

```bash
C:\Users\wpc>go get -u github.com/golang/text

C:\Users\wpc>

# 2. 然后把下载的文件copy到golang.org/x目录下
copy D:\Program Files\GOPATH\src\github.com\golang\text D:\Program Files\GOPATH\src\golang.org\x\text

# 3.然后执行安装
C:\Users\wpc>go install golang.org/x/text

C:\Users\wpc>
```

需要注意：

```bash
https://github.com/golang/text 和 https://golang.org/x/text 有的说是一样的有的说不一样
```

最好的解决办法是：

用god替换为god

```bash

```

