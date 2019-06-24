package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// go实现http文件上传小功能
const tp1 = `<html>
<head>
<title>上传文件</title>
</head>
<body>
  <form enctype="multipart/form-data" action="/upload" method="post">
    <input type="file" name="uploadfile" />
    <input type="hidden" name="token" value="{...{.}...}" />
    <input type="submit" value="upload" />
  </form>
</body>
</html>
`

// 写一个上传页面的html
func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tp1))
}

// 上传文件函数
func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	f, err := os.OpenFile(handler.Filename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
	fmt.Fprintln(w, "上传文件成功！")
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/upload", upload)
	http.ListenAndServe(":1789", nil)
}
