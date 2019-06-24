package main

import (
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"strconv"
	"strings"

	"crypto/tls"

	"time"

	"../../pkg/elasticsearch"
	"../../pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"
)



//构建部署应用操作之统计部署环节耗时。
//openshift project:build-deploy-app

// 索引mapping
const deployMapping = `
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
	//curl -H "Authorization: Bearer MdkqUn2AFrcBeYxKN6TKoEErcEHdC0YiRdlx1A5kHWE" https://openshift-master:8443/oapi/v1/namespaces/build-deploy-app/pods?watch=true
	url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/pods?watch=true&beijing=true"
	//url := "wss://openshift-master.m8.ccs:443/api/v1/namespaces/pt-ec/pods?watch=true&beijing=true"
	newReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	newReq.Header.Set("Sec-WebSocket-Protocol", base64url)
	//newReq.Header.Set("Sec-WebSocket-Protocol", "base64url.bearer.authorization.k8s.io.Slh4OFVtM216NlVxWnpIaVZZREkzOVBPeDBvTXYwdnl0eVcza1lGZFZIbw,undefined")
	newReq.Header.Set("Origin", "https://"+newReq.URL.Host)
	newReq.Header.Set("Access-Control-Allow-Origin", "*")

	d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sendConn, _, err := d.Dial(newReq.URL.String(), newReq.Header)
	if err != nil {
		fmt.Println(err)
		return
	}

	elasticsearch.CreateIndex(elasticClient, "build-deploy-app", deployMapping)

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
		//
		elasticsearch.CreateIndex(elasticClient, "build-deploy-app", deployMapping)

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

		//如果是Succeeded状态则直接跳过，Succeeded状态一般是build pod完成之后遗留没删除的。重新watch新的结果进行处理分析
		podIsSucceeded := gjson.Get(string(msg), "object.status.phase")
		if strings.EqualFold(podIsSucceeded.String(), "Succeeded") {
			continue
		}

		// 如果type是DELETED则过滤掉,这种类型是新pod起来之后需要干掉的原始的pod
		watchType := gjson.Get(string(msg), "type")
		if strings.EqualFold(watchType.String(), "DELETED") {
			continue
		}

		// 如果watch pod中存在object.metadata.deletionTimestamp表示是旧的pod正在被删除，过滤掉，我们只需要通过新的pod
		podIsDelete := gjson.Get(string(msg), "object.metadata.deletionTimestamp")
		if podIsDelete.Exists() {
			continue
		}

		// 解析watch返回的结果并重新组装
		watchDeployAssemble := parseWatchDeploy(string(msg))

		duration := gjson.Get(watchDeployAssemble, "duration")
		if strings.ContainsAny(duration.String(), "-") {
			continue
		}

		// 输出重新组装之后的watch deploy结果
		fmt.Println(watchDeployAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("build-deploy-app").
			Type("build-deploy-app").
			BodyString(watchDeployAssemble).
			Do()
		if err != nil {
			panic(err)
			return
		}
	}
	fmt.Println(newReq.Header)
}

/**
解析watch deploy pod返回结果
*/
func parseWatchDeploy(watchDeploy string) string {
	// 定义变量用于存储pod开始running的时间点
	var startRunning gjson.Result

	containerStatuses := gjson.Get(watchDeploy, "object.status.containerStatuses")
	if containerStatuses.Exists() {
		re := containerStatuses.Array()
		for _, v := range re {
			startedAt := v.Get("state.running.startedAt")
			if !startedAt.Exists() {
				continue
			}
			startRunning = startedAt
		}
	}

	/**
	type  watch类型
	object.kind  资源类型
	object.metadata.name 资源名称
	object.metadata.namespace  资源所属namespace
	object.metadata.creationTimestamp 资源创建时间
	object.status.hostIP 资源所在node节点IP
	object.status.podIP 资源如果是pod则是分配给pod的IP
	object.status.phase 资源所属阶段，资源状态
	object.status.startTime 资源创建时间
	*/
	watchDeployResult := gjson.GetMany(watchDeploy, "type", "object.kind", "object.metadata.name",
		"object.metadata.namespace", "object.metadata.creationTimestamp", "object.status.hostIP",
		"object.status.podIP", "object.status.phase", "object.status.startTime")
	return assembleWatchDeploy(startRunning.String(), watchDeployResult)
}

/**
重新组装watch deploy结果
*/
func assembleWatchDeploy(startRunning string, a []gjson.Result) string {

	toBeChangeStartRunning := startRunning
	toBeChangeStartTime := a[8].String()

	//转化所需模板
	timeLayout := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	// 使用模板在对应时区转化为time.time类型,go里面转换时间搓是10位，也就是转为的是秒不是毫秒
	// 字符串类型解析为时间类型
	startRunningParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartRunning, time.UTC)
	startTimeParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartTime, time.UTC)

	// 时间转为时间搓
	startRunningParseToTimestamp := startRunningParse.Unix() //这个转换时间搓默认(无论是time.Local还是time.UTC)会加上8个小时,why?,也好，这样就自动加了8小时
	startTimeParseToTimestamp := startTimeParse.Unix()

	// 时间搓转回为时间
	// 设置时间搓，使用模板格式化为日期字符串,因为时间搓已经加了8小时，所以再转为字符串的话这个时间就和之前的多8小时了。
	startRunningAddEightHours := time.Unix(startRunningParseToTimestamp, 0).Format(timeLayout)
	startTimeAddEightHours := time.Unix(startTimeParseToTimestamp, 0).Format(timeLayout)
	systemCurrentTime := time.Unix(time.Now().Unix(), 0).Format(timeLayout)

	// 统计pod从创建到running需要的时间
	podCreateDuration := startRunningParseToTimestamp - startTimeParseToTimestamp

	if startRunning == "" {
		startRunningAddEightHours = "2006-01-02T15:04:05+08:00"
		podCreateDuration = -1
	}

	// 持续时间字段转为时分秒格式的字符串
	duration := utils.HoursMintuesSeconds(podCreateDuration)

	watchDeployResult := `{"systemTimestamp":` + `"` + systemCurrentTime + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		//",\"watch_object_metadata_creationTimestamp\":" + `"` + a[4].String() + `"` +
		//",\"hostIP\":" + `"` + a[5].String() + `"` +
		//",\"podIP\":" + `"` + a[6].String() + `"` +
		",\"startTime\":" + `"` + startTimeAddEightHours + `"` +
		",\"runningTime\":" + `"` + startRunningAddEightHours + `"` +
		",\"durationInt\":" + `"` + strconv.FormatInt(podCreateDuration, 10) + `"` +
		",\"duration\":" + `"` + duration + `"` +
		",\"status\":" + `"` + a[7].String() + `"` +
		"}"
	return watchDeployResult
}
