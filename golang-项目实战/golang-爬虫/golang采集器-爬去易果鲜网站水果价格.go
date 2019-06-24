package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

//golang 采集器

type MySpider struct {
	indexUrl string
}

func (this MySpider) readUrlBody() (string, error) {
	resp, err := http.Get(this.indexUrl)
	if err != nil {
		return "err", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "err", err
	}
	return string(body), err
}
func (this MySpider) catchCategoryUrl() []string {
	body, _ := this.readUrlBody()
	rcg := regexp.MustCompile(`class="catalogs-ad">(?sU:.*)<a href="(http://www.yiguo.com/products/.*?_channelhome.html)">`)
	urls := rcg.FindAllStringSubmatch(body, -1)
	//var cateUrl []string
	cateUrl := make([]string, len(urls))
	for i, u := range urls {
		cateUrl[i] = u[1]
	}
	return cateUrl
}
func (this MySpider) catchProductInfo() string {
	body, _ := this.readUrlBody()
	rcg := regexp.MustCompile(`<div class="p_info clearfix">(?sU:.*)<div class="p_name"><a href="http://www.yiguo.com/product/(?U:.*).html" target="_blank">(.*?)</a></div>(?sU:.*)<div class="p_price">(?sU:.*)<strong>(.*?)</strong>(?sU:.*)</div>(?sU:.*)</div>`)
	result := rcg.FindAllStringSubmatch(body, -1)
	for i := range result {
		line := result[i]
		fmt.Println(line[1], "<<======>>", line[2])
	}
	return ""
}

func (this MySpider) run() string {
	cateUrls := this.catchCategoryUrl()
	for _, u := range cateUrls {
		this.indexUrl = u
		this.catchProductInfo()
		break
	}
	return ""
}
func main() {
	//ms := MySpider{} // 也 ok
	ms := new(MySpider)
	ms.indexUrl = "http://www.yiguo.com"
	ms.run()
}
