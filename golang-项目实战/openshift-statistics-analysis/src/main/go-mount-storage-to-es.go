package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"strconv"
	"strings"
	"time"

	"../../pkg/elasticsearch"
	"../../pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"
)

/**
挂载存储操作：模拟n用户挂载存储至1000个pod上
挂载存储的前提是pvc必须存在的吧，否则挂载存储的时候还需要创建pvc,有些时候更需要创建与之对应的pv
挂载存储是属于应用部署里面的，而创建pvc是属于存储场景测试中的
结果统计分析：统计挂载成功、失败数；统计pod运行成功的耗时
openshift project:mount-storage
*/

// 全局变量：挂载成功数，失败数全局变量,Pod状态为Running算成功，Pod状态为Failed算失败
// 挂载成功，失败总数应该是通过kibana或者grafana来出图的
//var totalSucc = 0
//var totalFail = 0

// 前提条件：等性能测试完了之后再执行watch操作

// 索引mapping
const watchMountStorageMapping = `
{
    "template": "mount-storage*",
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
        "mount-storage": {
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
                    "type": "long",
                    "index": "not_analyzed",
                    "store": true
                },
                 "podCreateDuration": {
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

	elasticsearch.CreateIndex(elasticClient, "mount-storage", watchMountStorageMapping)

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
		elasticsearch.CreateIndex(elasticClient, "mount-storage", watchMountStorageMapping)

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
		// 解析watch并重新组装成
		watchMountStorageAssemble := parseWatchMountStorage(string(msg))

		// 控制台输出组装之后的watch结果
		fmt.Println(watchMountStorageAssemble)

		// 数据入库
		_, err = elasticClient.Index().
			Index("mount-storage").
			Type("mount-storage").
			BodyString(watchMountStorageAssemble).
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
func parseWatchMountStorage(watchMountStorage string) string {
	// 如果Pod最终是Running状态则获取pod Running的时间；如果Pod最终状态是Failed则获取pod Failed的时间，Pod状态为Pending的情况已在上面直接过滤掉不处理
	var startRunningOrTerminated gjson.Result
	// pod有几种状态：Running、Failed、Successed,而Pending状态在上面已经过滤掉
	podStatus := gjson.Get(watchMountStorage, "object.status.phase")
	containerStatuses := gjson.Get(watchMountStorage, "object.status.containerStatuses")
	if strings.EqualFold(podStatus.String(), "Running") {
		if containerStatuses.Exists() {
			re := containerStatuses.Array()
			// Pod Running开始时间
			for _, v := range re {
				startedAt := v.Get("state.running.startedAt")
				if !startedAt.Exists() {
					continue
				}
				startRunningOrTerminated = startedAt
			}
		}
	}
	if strings.EqualFold(podStatus.String(), "Failed") {
		if containerStatuses.Exists() {
			re := containerStatuses.Array()
			// Pod Terminated结束时间
			for _, v := range re {
				terminatedAt := v.Get("state.terminated.finishedAt")
				if !terminatedAt.Exists() {
					continue
				}
				startRunningOrTerminated = terminatedAt
			}
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
	watchMountStorageResult := gjson.GetMany(watchMountStorage, "type", "object.kind", "object.metadata.name",
		"object.metadata.namespace", "object.metadata.creationTimestamp", "object.status.phase",
		"object.status.hostIP", "object.status.podIP", "object.status.startTime")

	return assembleWatchMountStorage(startRunningOrTerminated.String(), watchMountStorageResult)
}

/**
watch pod结果解析出的字段重新组装成Json格式
*/
func assembleWatchMountStorage(startRunningOrTerminated string, a []gjson.Result) string {

	toBeChangeStartRunningOrTerminated := startRunningOrTerminated
	toBeChangeStartTime := a[8].String()

	//转化所需模板
	timeLayout := time.RFC3339 // 2006-01-02T15:04:05Z07:00
	// 使用模板在对应时区转化为time.time类型,go里面转换时间搓是10位，也就是转为的是秒不是毫秒
	// 字符串类型解析为时间类型
	startRunningOrTerminatedParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartRunningOrTerminated, time.UTC)
	startTimeParse, _ := time.ParseInLocation(timeLayout, toBeChangeStartTime, time.UTC)

	// 时间转为时间搓
	startRunningOrTerminatedParseToTimestamp := startRunningOrTerminatedParse.Unix() //这个转换时间搓默认(无论是time.Local还是time.UTC)会加上8个小时,why?,也好，这样就自动加了8小时
	startTimeParseToTimestamp := startTimeParse.Unix()

	// 时间搓转回为时间
	// 设置时间搓，使用模板格式化为日期字符串,因为时间搓已经加了8小时，所以再转为字符串的话这个时间就和之前的多8小时了。
	startRunningOrTerminatedAddEightHours := time.Unix(startRunningOrTerminatedParseToTimestamp, 0).Format(timeLayout)
	startTimeAddEightHours := time.Unix(startTimeParseToTimestamp, 0).Format(timeLayout)
	systemCurrentTime := time.Unix(time.Now().Unix(), 0).Format(timeLayout)

	// 统计pod从创建到running需要的时间
	podCreateDuration := startRunningOrTerminatedParseToTimestamp - startTimeParseToTimestamp

	// 如果pod 状态是Running则pod总运行成功数加1，否则pod总运行失败数加1
	/*
		podIsRunningOrFailed := a[5].String()
		utils.HoursMintuesSeconds(podCreateDuration)

		if strings.EqualFold(podIsRunningOrFailed, "Running") {
			totalSucc++
		} else {
			totalFail++
		}
	*/

	watchMountStorageResult := `{"systemTimestamp":` + `"` + systemCurrentTime + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		//",\"watch_object_metadata_creationTimestamp\":" + `"` + a[4].String() + `"` +
		",\"status\":" + `"` + a[5].String() + `"` +
		",\"hostIP\":" + `"` + a[6].String() + `"` +
		",\"podIP\":" + `"` + a[7].String() + `"` +
		",\"startTime\":" + `"` + startTimeAddEightHours + `"` +
		",\"runningTimeOrTerminatedTime\":" + `"` + startRunningOrTerminatedAddEightHours + `"` +
		",\"podCreateDurationInt\":" + `"` + strconv.FormatInt(podCreateDuration, 10) + `"` +
		",\"podCreateDuration\":" + `"` + utils.HoursMintuesSeconds(podCreateDuration) + `"` +
		//",\"podMountSucc\":" + `"` + strconv.Itoa(totalSucc) + `"` +
		//",\"podMountFail\":" + `"` + strconv.Itoa(totalFail) + `"` +
		"}"

	return watchMountStorageResult
}
