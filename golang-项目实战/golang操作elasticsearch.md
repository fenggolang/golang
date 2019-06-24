#### golang操作es

``` go
package main

/**
参考：https://olivere.github.io/elastic/
*/

import (
	"context"
	"fmt"

	elastic "gopkg.in/olivere/elastic.v5"
)

// olivere.elastic5.0对应 es 5.x版本
// olivere/elastic3.0对应 es 2.x版本
type Tweet struct {
	User    string
	Message string
}

func main() {
	// Starting with elastic.v5, you must pass a context to execute each service
	// 从elastic.v5开始，必须要通过一个上下文才能执行每一个服务
	ctx := context.Background()

	//client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://elasticsearch.csg.com/")) // 容器方式部署的es,暂时没有测试通过
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://172.17.12.79:9200/")) //宿主机上部署的es
	if err != nil {
		fmt.Println("连接elasticsearch集群失败", err)
		return
	}

	fmt.Println("成功连接上es集群")

	for i := 0; i < 20; i++ {
		tweet := Tweet{User: "wpc", Message: "wpc like golang promramming"}
		_, err = client.Index().
			Index("wpc-index").
			Type("tweet").
			Id(fmt.Sprintf("%d", i)).
			BodyJson(tweet).
			Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
			return
		}
	}

	fmt.Println("数据写入es成功")
}
```

