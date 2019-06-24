package main



//构建部署应用操作之统计构建环节耗时。
//openshift project:build-deploy-app

import (
	"fmt"
	"net/http"
	"strconv"

	"crypto/tls"

	"time"

	"openshift-statistics-analysis/pkg/const"

	"openshift-statistics-analysis/pkg/elasticsearch"
	"openshift-statistics-analysis/pkg/utils"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"
)

// 索引mapping
const buildDeployAppMapping = `
{
    "template": "build-deploy-app*",
    "settings": {
	    "number_of_replicas": 0,
        "number_of_shards": 5,
		"index": {
			"store": {
				"compress": {
					"stored": true,
					"tv": true
				}
			}
		}
    },
    "mappings": {
        "build-deploy-app": {
            "_source": {
                "enabled": true
            },
			"_ttl":{
				"enabled": true,
				"default": "30d"
			},
			"_all" : {
				"enabled" : false
			},
            "properties": {
                "type": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                "kind": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "name": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "namespace": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "status": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "hostIP": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "podIP": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "durationInt": {
                    "type": "long",
                    "index": "not_analyzed",
                    "store": true
                },
                 "duration": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                }
            }
        }
    }
}
`

func main() {
	var token *string
	token, err := utils.GetToken(_const.OpenshiftUrl, _const.Username, _const.Password)
	base64url := "base64url.bearer.authorization.k8s.io." + utils.TokenToBase64Encode(*token) + ",undefined"

	elasticClient, err := elastic.NewClient(elastic.SetSniff(true), elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}

	url := "wss://cluster.prod.guizhou:8443/oapi/v1/namespaces/pt-ec/builds?watch=true"
	//url := "wss://openshift-master.m8.ccs:443/oapi/v1/namespaces/pt-ec/builds?watch=true"
	newReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	newReq.Header.Set("Sec-WebSocket-Protocol", base64url)
	newReq.Header.Set("Origin", "https://"+newReq.URL.Host)
	newReq.Header.Set("Access-Control-Allow-Origin", "*")

	d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sendConn, _, err := d.Dial(newReq.URL.String(), newReq.Header)
	if err != nil {
		fmt.Println(err)
		return
	}

	elasticsearch.CreateIndex(elasticClient, "build-deploy-app", buildDeployAppMapping)

	// 保持websocket回话连接，保证websocket在watch pod的信息没有变化的情况下不会自动断开
	sendConn.SetReadDeadline(time.Time{})
	i := 0
	// 坑爹的中间节点可能会认为一份连接在一段时间内没有数据发送就等于失效，它们会自作主张的切断这些连接
	// 解决方案，WebSocket 的设计者们也早已想过。就是让服务器和客户端能够发送 Ping/Pong
	// Ping和Pong则没有明确定定义，一般用于心跳消息，而且Pong一般是讲Ping的消息原封不动的发送回去
	go func() {
		for {
			//fmt.Printf("交互数据保持连接次数：%v \n", i)
			sendConn.WriteMessage(1, []byte("hello"))
			//sendConn.PingHandler() // 发送ping包
			time.Sleep(time.Second * 5)
			i += 1
		}

	}()
LABEL:
	for {
		_, msg, err := sendConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//
		objectStatusDuration := gjson.Get(string(msg), "object.status.duration")
		if !objectStatusDuration.Exists() {
			goto LABEL
		}
		fmt.Println(string(msg))
		// 解析watch返回的结果并重新组装
		watchBuildAssemble := parseWatchBuild(string(msg))

		// 输出重新组装之后的watch builds结果
		fmt.Println(watchBuildAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("build-deploy-app").
			Type("build-deploy-app").
			BodyString(watchBuildAssemble).
			Do()
		if err != nil {
			panic(err)
			return
		}
	}
	fmt.Println(newReq.Header)
}

/**
解析watch builds返回结果
*/
func parseWatchBuild(watchBuild string) string {
	/**
	type  watch类型
	object.kind  资源类型
	object.metadata.name 资源名称
	object.metadata.namespace  资源所属namespace
	object.metadata.completionTimestamp 资源完成时间
	object.status.duration 资源创建持续时间
	object.status.phase 资源所属阶段，资源状态
	object.status.startTimestamp 资源创建时间
	*/
	watchBuildResult := gjson.GetMany(watchBuild, "type", "object.kind", "object.metadata.name",
		"object.metadata.namespace", "object.status.completionTimestamp", "object.status.duration",
		"object.status.phase", "object.status.startTimestamp")
	return assembleWatchBuild(watchBuildResult)
}

/**
重新组装watch builds结果
*/
func assembleWatchBuild(a []gjson.Result) string {
	timeLayout := time.RFC3339
	// 需要处理的时间字段
	toBeChangeCompletionTimestamp := a[4].String()
	toBeChangeStartTimestamp := a[7].String()

	// 字符串类型解析为时间类型
	completionTimestampParse, _ := time.ParseInLocation(timeLayout, toBeChangeCompletionTimestamp, time.UTC)
	startTimestampParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartTimestamp, time.UTC)

	// 时间转为时间搓，自动加8小时
	completionTimestampParseToTimestamp := completionTimestampParse.Unix()
	startTimestampParseToTimestamp := startTimestampParse.Unix()

	// 时间搓转回时间，得到增加8小时候的时间
	completionTimestampAddEightHours := time.Unix(completionTimestampParseToTimestamp, 0).Format(timeLayout)
	startTimestampADDEightHours := time.Unix(startTimestampParseToTimestamp, 0).Format(timeLayout)

	// 持续时间字段转为时分秒格式的字符串
	duration := utils.HoursMintuesSeconds(a[5].Int())
	//fmt.Println("a[5]:", a[5])
	//fmt.Println("a[5].Int():", a[5].Int())
	//time.Sleep(2 * time.Second)
	watchBuildResult := `{"systemTimestamp":` + `"` + time.Unix(time.Now().Unix(), 0).Format(timeLayout) + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		",\"startTimestamp\":" + `"` + startTimestampADDEightHours + `"` +
		",\"completionTimestamp\":" + `"` + completionTimestampAddEightHours + `"` +
		",\"durationInt\":" + `"` + strconv.FormatInt(a[5].Int()/1000000000, 10) + `"` +
		",\"duration\":" + `"` + duration + `"` +
		",\"status\":" + `"` + a[6].String() + `"` +
		"}"

	return watchBuildResult
}
