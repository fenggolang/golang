#### golang封装

- go语言仅支持封装，不支持继承和多态
- go语言没有class,只有struct,对struct进行封装实现面向对象编程
- go语言没有构造函数的说法
- 只有使用指针才可以改变结构内容
- nil指针也可以调用方法！

#### golang接受者

- 值接收者：go语言特有
- 指针接收者：要改变内容必须使用指针接受者，结构过大也考虑使用指针接受者，如果有指针接受者，最好都是指针接受者

####  golang包和封装

- 名称一般使用驼峰式
- 首字母大写：public
- 首字母小写：private
- 包：每个目录一个包，但是和其他语言不一样的是，每一个包名称不一定和目录名称一致，比如目录名称是example,但是go代码最上面的package main定义一个main包，一个目录下只能有一个包，也就是比如说一个代码目录是example,里面有许多文件，文件最上面的package 包名称必须总共只有一个。
- 为结构定义的方法必须放在同一个包中，可以是不同文件

#### golang扩展已有类型

​	如何扩充系统类型或者别人的类型：第一种是定义别名，第二种是使用组合

#### golang GOPATH

​	找包的路径，golang.org包是被墙了的，我们可以通过执行来安装gopm

```bash
# 下载安装gopm
go get -v github.com/gpmgo/gopm

# 然后用gopm来安装golang.org的包
gopm get -g -v -u golang.org/x/tools/cmd/goimports
```



