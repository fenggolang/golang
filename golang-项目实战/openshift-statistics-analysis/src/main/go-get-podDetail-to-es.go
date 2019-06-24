package main



/**
需求：获取指定项目下的所有pod的详细信息并入库(文本或者写入es然后通过kibana查询es数据导出成excel表格)
pod信息字段：
namespace pod_name pod_id container_name dc_name
*/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"time"

	"../../pkg/elasticsearch"
	"../../pkg/utils"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"

	"strings"
)

// 索引mapping
const watchPodDetailMapping = `
{
    "template": "pod-detail*",
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
        "pod-detail": {
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
                "namespace": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                "pod_name": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "pod_id": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "container_name": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "dc_name": {
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

	//url := "wss://openshift-master.m8.ccs:443/api/v1/namespaces/logging/pods?watch=true"
	url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/pods?watch=true&region=beijing"
	newReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	newReq.Header.Set("Sec-WebSocket-Protocol", base64url)
	//newReq.Header.Set("Sec-WebSocket-Protocol", "base64url.bearer.authorization.k8s.io.Slh4OFVtM216NlVxWnpIaVZZREkzOVBPeDBvTXYwdnl0eVcza1lGZFZIbw, undefined")
	newReq.Header.Set("Origin", "https://"+newReq.URL.Host)
	newReq.Header.Set("Access-Control-Allow-Origin", "*")

	d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sendConn, _, err := d.Dial(newReq.URL.String(), newReq.Header)
	if err != nil {
		fmt.Println(err)
		return
	}

	elasticsearch.CreateIndex(elasticClient, "pod-detail", watchPodDetailMapping)

	// 保持websocket回话连接，保证websocket在watch pod的信息没有变化的情况下不会自动断开
	sendConn.SetReadDeadline(time.Time{})
	i := 0
	go func() {
		for {
			//fmt.Printf("交互数据保持连接次数：%v \n", i)
			sendConn.WriteMessage(1, []byte("hello"))
			time.Sleep(time.Second * 5)
			i += 1
		}

	}()

	for {
		elasticsearch.CreateIndex(elasticClient, "pod-detail", watchPodDetailMapping)

		_, msg, err := sendConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		//如果一直是Pending状态则直接跳过，重新watch新的结果进行处理分析
		podIsPending := gjson.Get(string(msg), "object.status.phase")
		if strings.EqualFold(podIsPending.String(), "Pending") {
			continue
		}

		// 解析watch并重新组装
		watchPodDetailAssemble := parseWatchPodDetail(string(msg))

		//如果是Succeeded状态则直接跳过，Succeeded状态一般是build pod完成之后遗留没删除的。重新watch新的结果进行处理分析
		podIsSucceeded := gjson.Get(string(msg), "object.status.phase")
		if strings.EqualFold(podIsSucceeded.String(), "Succeeded") {
			continue
		}

		// 如果部署名为空一般是daemonset,过滤掉
		podIsHasDeployName := gjson.Get(string(msg), "object.metadata.labels.deploymentconfig")
		if strings.EqualFold(podIsHasDeployName.String(), "") {
			continue
		}
		// 控制台输出组装之后的watch pod
		fmt.Println(watchPodDetailAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("pod-detail").
			Type("pod-detail").
			BodyString(watchPodDetailAssemble).
			Do() //执行数据入索引库操作
		if err != nil {
			panic(err)
			return
		}
	}
	fmt.Println(newReq.Header)
}

/**
解析watch pod返回结果
*/
func parseWatchPodDetail(watchPodDetail string) string {
	var containerName gjson.Result

	containers := gjson.Get(watchPodDetail, "object.spec.containers")
	if containers.Exists() {
		re := containers.Array()
		for _, v := range re {
			tempContainerName := v.Get("name")
			if !tempContainerName.Exists() {
				continue
			}
			containerName = tempContainerName
		}
	}

	/**
	type  watch类型
	object.kind  资源类型
	object.metadata.namespace  资源所属namespace
	object.metadata.name pod_name
	object.metadata.uid pod_id
	object.metadata.labels.deploymentconfig 部署名称
	object.spec.containers.name 容器名称
	*/
	watchPodDetailResult := gjson.GetMany(watchPodDetail, "type", "object.kind", "object.metadata.namespace",
		"object.metadata.name", "object.metadata.uid", "object.spec.containers.name", "object.metadata.labels.deploymentconfig", "object.status.startTime")

	return assembleWatchPodDetail(containerName.String(), watchPodDetailResult)
}

/**
watch pod结果解析出的字段重新组装成Json格式
*/
func assembleWatchPodDetail(containerName string, a []gjson.Result) string {

	toBeChangeTimestamp := a[7].String()
	timeLayout := time.RFC3339
	timestampParse, _ := time.ParseInLocation(timeLayout, toBeChangeTimestamp, time.UTC)
	timestampParseToTimestamp := timestampParse.Unix()
	timestampAddEightHours := time.Unix(timestampParseToTimestamp, 0).Format(timeLayout)
	//fmt.Println(timestampAddEightHours)

	//timestampAddEightHours, _ := time.Parse(time.RFC3339, a[7].String())
	watchPodDetailResult := `{"type":` + `"` + a[0].String() + `"` +
		//",\"timestamp\":" + `"` + a[7].String() + `"` +
		",\"timestamp\":" + `"` + timestampAddEightHours + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"namespace\":" + `"` + a[2].String() + `"` +
		",\"pod_name\":" + `"` + a[3].String() + `"` +
		",\"pod_id\":" + `"` + a[4].String() + `"` +
		",\"container_name\":" + `"` + containerName + `"` +
		",\"dc_name\":" + `"` + a[6].String() + `"` +
		"}"
	return watchPodDetailResult
}
