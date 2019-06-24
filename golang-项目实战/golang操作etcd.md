#### golang连接etcd

```go
package main

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.18.2.13:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接etcd失败:", err)
		return
	}

	fmt.Println("连接成功")
	defer cli.Close()
}
```

---

#### golang 存取etcd

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.18.2.13:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}

	fmt.Println("连接成功")
	defer cli.Close()

	//设置1秒超时，访问etcd有超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	//操作etcd: 写入key-value
	_, err = cli.Put(ctx, "/wpc/interests", "golang")
	// 操作完毕，取消etcd
	cancel()
	if err != nil {
		fmt.Println("数据写入失败：", err)
		return
	}

	// 取值，设置超时为1秒
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	// 操作etcd:获取etcd的key对应的value的值
	resp, err := cli.Get(ctx, "/wpc/interests")
	// 操作完毕，取消etcd
	cancel()
	if err != nil {
		fmt.Println("获取数据失败：", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s：%s\n", ev.Key, ev.Value)
	}
}
```

---

#### golang etcd监听watch

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"172.18.2.13:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println("连接失败：", err)
		return
	}

	fmt.Println("连接成功")
	defer cli.Close()

	cli.Put(context.Background(), "/wpc/like", "golanging")
	for {
		// watch key监听节点
		rch := cli.Watch(context.Background(), "/wpc/like")
		for wresp := range rch {
			for _, ev := range wresp.Events {
				fmt.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}
		}
	}
}
```

