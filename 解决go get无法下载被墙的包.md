#### 方式一：

```bash
# 使用gopm(Go Package Manager)代替go下载,是go上的包管理工具，十分好用
# 1. 下载安装gopm
go get -u github.com/gpmgo/gopm
# 2. 使用gopm安装被墙的包
gopm get github.com/Shopify/sarama
```

---

#### 方式二：

```bash
#  golang 在 github 上建立了一个镜像库，如 https://github.com/golang/net 即是 https://golang.org/x/net 的镜像库.获取 golang.org/x/net 包（其他包类似），其实只需要以下步骤：
mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x
git clone https://github.com/golang/net.git
```

