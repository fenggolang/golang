package main


/**
扩容操作：统计请求发送成功到调度到资源时间、调度到资源到pod运行成功的时间
openshift project:scale-pod
pod可能通过dc控制也可能是通过rc控制，若通过rc控制则下面改为rc
oc scale --replicas=2 dc hawkular-alert-extend
*/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"strconv"

	"../../pkg/elasticsearch"
	"../../pkg/utils"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"

	"strings"
	"time"
)

// 索引mapping
const watchScalePodMapping = `
{
    "template": "scale-pod*",
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
        "scale-pod": {
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
                    "type": "integer",
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

	//curl -H "Authorization: Bearer R734aFCxWg_G8FXcX1WNFVAZ3whctrZpbiifBmdIeFw" https://openshift-master:8443/api/v1/namespaces/build-deploy-app/builds?watch=true
	//url := "wss://openshift-master:8443/api/v1/namespaces/scale-pod/pods?watch=true"
	url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/pods?watch=true&region=beijing"
	newReq, err := http.NewRequest("GET", url, nil)
	newReq.Close = true
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

	elasticsearch.CreateIndex(elasticClient, "scale-pod", watchScalePodMapping)

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
		elasticsearch.CreateIndex(elasticClient, "scale-pod", watchScalePodMapping)

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
		watchScalePodAssemble, scaleFlag := parseWatchScalePod(string(msg))
		if scaleFlag == 1 {
			continue
		}
		// 如果podCreateDuration为-1表示可以直接过滤掉，此时pod是Pending状态
		podCreateDuration := gjson.Get(watchScalePodAssemble, "podCreateDuration")
		if strings.EqualFold(string(podCreateDuration.Str), "-1") {
			continue
		}

		//如果是Succeeded状态则直接跳过，Succeeded状态一般是build pod完成之后遗留没删除的。重新watch新的结果进行处理分析
		podIsSucceeded := gjson.Get(string(msg), "object.status.phase")
		if strings.EqualFold(podIsSucceeded.String(), "Succeeded") {
			continue
		}

		// 控制台输出组装之后的watch pod
		fmt.Println(watchScalePodAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("scale-pod").
			Type("scale-pod").
			BodyString(watchScalePodAssemble).
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
func parseWatchScalePod(watchScalePod string) (string, int) {
	// 定义变量用于存储pod开始running的时间点
	var startRunning gjson.Result

	containerStatuses := gjson.Get(watchScalePod, "object.status.containerStatuses")
	if containerStatuses.Exists() {
		re := containerStatuses.Array()
		for _, v := range re {
			startedAt := v.Get("state.running.startedAt")
			// 如果pod yaml内容没有startedAt则表示没有running,出在waiting状态，因为此时正在ContainerCreating,
			// 但是有的可能一直都不能running该如何处理
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
	object.status.phase 资源所属阶段，资源状态
	object.status.hostIP 资源所在node节点IP
	object.status.podIP 资源如果是pod则是分配给pod的IP
	object.status.startTime 资源创建时间
	*/
	watchScalePodResult := gjson.GetMany(watchScalePod, "type", "object.kind", "object.metadata.name",
		"object.metadata.namespace", "object.metadata.creationTimestamp", "object.status.phase",
		"object.status.hostIP", "object.status.podIP", "object.status.startTime", "object.metadata.deletionTimestamp")

	return assembleWatchScalePod(startRunning.String(), watchScalePodResult)
}

/**
watch pod结果解析出的字段重新组装成Json格式
*/
func assembleWatchScalePod(startRunning string, a []gjson.Result) (string, int) {
	// 处理扩容之后，如果缩容的话，会重复watch到之前watch到的数据，需要通过自定义标记过滤掉
	scaleFlag := 0
	// 如果object.metadata.deletionTimestamp有值表示是缩容操作，会重复watch数据到es，所以标记为1，过滤掉
	if !strings.EqualFold(a[9].String(), "") {
		scaleFlag = 1
	}

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
	duration := utils.HoursMintuesSeconds(podCreateDuration)
	if startRunning == "" {
		startRunningAddEightHours = "2006-01-02T15:04:05+08:00"
		podCreateDuration = -1
	}

	watchScalePodResult := `{"systemTimestamp":` + `"` + systemCurrentTime + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		//",\"watch_object_metadata_creationTimestamp\":" + `"` + a[4].String() + `"` +
		",\"status\":" + `"` + a[5].String() + `"` +
		",\"hostIP\":" + `"` + a[6].String() + `"` +
		",\"podIP\":" + `"` + a[7].String() + `"` +
		",\"startTime\":" + `"` + startTimeAddEightHours + `"` +
		",\"startedAt\":" + `"` + startRunningAddEightHours + `"` +
		",\"durationInt\":" + `"` + strconv.FormatInt(podCreateDuration, 10) + `"` +
		",\"duration\":" + `"` + duration + `"` +
		"}"

	return watchScalePodResult, scaleFlag
}
