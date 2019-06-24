package main


/**
删除pod操作：统计删除pod的耗时；新pod运行成功耗时
openshift project:delete-pod
*/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"strings"

	"openshift-statistics-analysis/pkg/utils"

	"openshift-statistics-analysis/pkg/elasticsearch"

	"time"

	"strconv"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"
)

// 索引mapping
const watchDeletePodMapping = `
{
    "template": "delete-pod*",
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
        "delete-pod": {
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
                 "podCreateDurationInt": {
                    "type": "integer",
                    "index": "not_analyzed",
                    "store": true
                },
                 "podCreateDuration": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "podDeleteDurationInt": {
                    "type": "integer",
                    "index": "not_analyzed",
                    "store": true
                },
                 "podDeleteDuration": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "deletionGracePeriodSeconds": {
                    "type": "integer",
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
	//url := "ws://openshift-master.m8.ccs:443/api/v1/namespaces/logging/pods?watch=true"
	//url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/pods?watch=true&region=beijing"
	url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/pods?watch=true"
	newReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	newReq.Header.Set("Sec-WebSocket-Protocol", base64url)
	newReq.Header.Set("Origin", "https://"+newReq.URL.Host)
	newReq.Header.Set("Access-Control-Allow-Origin", "*")
	//d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, HandshakeTimeout: 30 * time.Hour}

	d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sendConn, _, err := d.Dial(newReq.URL.String(), newReq.Header)
	if err != nil {
		fmt.Printf("createConnErr:%v \n", err)
		return
	}
	//sendConn.SetReadDeadline(time.Now().Add(time.Duration(10) * time.Hour))
	//sendConn.SetReadDeadline(time.Now().Add(1 * time.Hour))
	//if err = sendConn.SetReadDeadline(time.Now().Add(time.Hour * 3)); err != nil {
	//	panic(err)
	//}
	//go func() {
	//	time.Sleep(2 * time.Hour)
	//	fmt.Println("休眠了2小时秒才输出")
	//}()
	//startTime := time.Now()
	//fmt.Println("开始：", startTime.UTC())
	//outTime := startTime.Add(time.Second * 120)
	//fmt.Println("断开时间：", outTime.UTC())
	sendConn.SetReadDeadline(time.Time{})
	elasticsearch.CreateIndex(elasticClient, "delete-pod", watchDeletePodMapping)
	//utils.KeepAlive(sendConn, 10*time.Hour)
	//var conn net.Conn
	i := 0
	go func() {
		for {
			//fmt.Printf("交互数据保持连接次数：%v \n", i)
			sendConn.WriteMessage(1, []byte("hello"))
			time.Sleep(time.Second * 5)
			i += 1
		}

	}()
LABEL:

	//fmt.Println("循环")

	for {
		elasticsearch.CreateIndex(elasticClient, "delete-pod", watchDeletePodMapping)

		//utils.KeepAlive(sendConn, 3*time.Second)
		//utils.HandleConnection(conn, 1111)
		//if err = sendConn.SetReadDeadline(time.Now().Add(time.Hour * 3)); err != nil {
		//	panic(err)
		//}

		_, msg, err := sendConn.ReadMessage()

		if err != nil {
			fmt.Printf("readMessageErr:%v \n", err)
			break
		}

		//如果一直是Pending状态则直接跳过，重新watch新的结果进行处理分析
		podIsPending := gjson.Get(string(msg), "object.status.phase")
		if strings.EqualFold(podIsPending.String(), "Pending") {
			continue
		}

		// 解析watch并重新组装成
		watchDeletePodAssemble := parseWatchDeletePod(string(msg))
		podDeleteType := gjson.Get(watchDeletePodAssemble, "type")
		podDeleteDuration := gjson.Get(watchDeletePodAssemble, "podDeleteDuration")
		podCreateDuration := gjson.Get(watchDeletePodAssemble, "podCreateDuration")
		deleteDuration, _ := strconv.ParseInt(podDeleteDuration.String(), 10, 32)
		if deleteDuration < -1 {
			continue
		}

		// 如果watch type等于MODIFIED且podDeleteDuration=-1则过滤掉
		if strings.EqualFold(podDeleteType.String(), "MODIFIED") {
			if strings.EqualFold(podCreateDuration.String(), "-1") {
				goto LABEL
			}
		}

		// 控制台输出组装之后的watch pod
		fmt.Println("time:%v info:%v \n", time.Now().UTC(), watchDeletePodAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("delete-pod").
			Type("delete-pod").
			BodyString(watchDeletePodAssemble).
			Do() //执行数据入索引库操作
		if err != nil {
			panic(err)
			return
		}
	}
	fmt.Println("结束", time.Now().UTC())
	//fmt.Println(newReq.Header)
}

/**
解析watch pod返回结果
*/
func parseWatchDeletePod(watchDeletePod string) string {
	// 定义变量用于存储pod开始running的时间点
	var startRunning gjson.Result

	containerStatuses := gjson.Get(watchDeletePod, "object.status.containerStatuses")
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
	object.metadata.deletionGracePeriodSeconds 优雅删除时间
	object.metadata.deletionTimestamp 接收到删除命令的时间
	object.status.phase 资源所属阶段，资源状态
	object.status.hostIP 资源所在node节点IP
	object.status.podIP 资源如果是pod则是分配给pod的IP
	object.status.startTime 资源创建时间
	*/
	watchPodResult := gjson.GetMany(watchDeletePod, "type", "object.kind", "object.metadata.name", "object.metadata.namespace",
		"object.metadata.creationTimestamp", "object.metadata.deletionGracePeriodSeconds", "object.metadata.deletionTimestamp",
		"object.status.phase", "object.status.hostIP", "object.status.podIP", "object.status.startTime")

	return assembleWatchDeletePod(startRunning.String(), watchPodResult)
}

/**
watch pod结果解析出的字段重新组装成Json格式
*/
func assembleWatchDeletePod(startRunning string, a []gjson.Result) string {

	toBeChangeStartRunning := startRunning
	toBeChangeDeletionTimestamp := a[6].String()
	toBeChangeStartTime := a[10].String()

	//转化所需模板
	timeLayout := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	// 使用模板在对应时区转化为time.time类型,go里面转换时间搓是10位，也就是转为的是秒不是毫秒
	// 字符串类型解析为时间类型
	startRunningParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartRunning, time.UTC)
	startTimeParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartTime, time.UTC)
	deletionTimestampParse, _ := time.ParseInLocation(timeLayout, toBeChangeDeletionTimestamp, time.UTC)

	// 时间转为时间搓
	startRunningParseToTimestamp := startRunningParse.Unix() //这个转换时间搓默认(无论是time.Local还是time.UTC)会加上8个小时,why?,也好，这样就自动加了8小时
	startTimeParseToTimestamp := startTimeParse.Unix()
	deletionTimestampParseToTimestamp := deletionTimestampParse.Unix()

	// 时间搓转回为时间
	// 设置时间搓，使用模板格式化为日期字符串,因为时间搓已经加了8小时，所以再转为字符串的话这个时间就和之前的多8小时了。
	startRunningAddEightHours := time.Unix(startRunningParseToTimestamp, 0).Format(timeLayout)
	startTimeAddEightHours := time.Unix(startTimeParseToTimestamp, 0).Format(timeLayout)
	deletionTimestampAddEightHours := time.Unix(deletionTimestampParseToTimestamp, 0).Format(timeLayout)
	systemCurrentTime := time.Unix(time.Now().Unix(), 0).Format(timeLayout)

	// 统计pod从创建到running需要的时间
	podCreateDuration := startRunningParseToTimestamp - startTimeParseToTimestamp
	// 统计pod从接收删除命令到真正DELETED需要的时间
	podDeleteDuration := time.Now().Unix() - deletionTimestampParseToTimestamp
	if a[6].String() == "" {
		deletionTimestampAddEightHours = "2006-01-02T15:04:05+08:00"
		podDeleteDuration = -1
	}
	if startRunning == "" {
		startRunningAddEightHours = "2006-01-02T15:04:05+08:00"
		podCreateDuration = -1
	}

	createDuration := utils.HoursMintuesSeconds(podCreateDuration)
	deleteDuration := utils.HoursMintuesSeconds(podDeleteDuration)
	watchDeletePodResult := `{"systemTimestamp":` + `"` + systemCurrentTime + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		//",\"watch_object_metadata_creationTimestamp\":" + `"` + a[4].String() + `"` +
		//",\"watch_object_metadata_deletionGracePeriodSeconds\":" + `"` + a[5].String() + `"` +
		",\"hostIP\":" + `"` + a[8].String() + `"` +
		",\"podIP\":" + `"` + a[9].String() + `"` +
		",\"status\":" + `"` + a[7].String() + `"` +
		",\"deletionTimestamp\":" + `"` + deletionTimestampAddEightHours + `"` +
		",\"startTime\":" + `"` + startTimeAddEightHours + `"` +
		",\"runningTime\":" + `"` + startRunningAddEightHours + `"` +
		",\"podCreateDurationInt\":" + `"` + strconv.FormatInt(podCreateDuration, 10) + `"` +
		",\"podCreateDuration\":" + `"` + createDuration + `"` +
		",\"podDeleteDurationInt\":" + `"` + strconv.FormatInt(podDeleteDuration, 10) + `"` +
		",\"podDeleteDuration\":" + `"` + deleteDuration + `"` +
		"}"

	return watchDeletePodResult
}
