package main



import (
	"crypto/tls"
	"fmt"
	"net/http"
	"openshift-statistics-analysis/pkg/const"
	"strconv"
	"strings"
	"time"

	"openshift-statistics-analysis/pkg/elasticsearch"

	"../../pkg/utils"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"gopkg.in/olivere/elastic.v3"
)

/**
创建PVC操作：统计pvc从创建到分配成功的时间
openshift project:create-pvc
*/

// 索引mapping
const watchCreatePvcMapping = `
{
    "template": "create-pvc*",
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
        "create-pvc": {
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
                 "storageClassName": {
                    "type": "string",
                    "index": "not_analyzed",
                    "store": true
                },
                 "pvName": {
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
                 "createPvcDurationInt": {
                    "type": "integer",
                    "index": "not_analyzed",
                    "store": true
                },
              	"createPvcDuration": {
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
	url := "wss://cluster.prod.guizhou:8443/api/v1/namespaces/pt-ec/persistentvolumeclaims?watch=true&region=beijing"
	//url := "wss://openshift-master:8443/api/v1/namespaces/create-pvc/persistentvolumeclaims?watch=true"
	newReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	newReq.Header.Set("Sec-WebSocket-Protocol", base64url)
	//newReq.Header.Set("Sec-WebSocket-Protocol", "base64url.bearer.authorization.k8s.io.MXIwOWN2RWZVWElMR3Q2ZGs3ZE5EbllxSEpqcnZTeTFnS1Vud0R4Nmt6VQ, undefined")
	newReq.Header.Set("Origin", "https://"+newReq.URL.Host)
	newReq.Header.Set("Access-Control-Allow-Origin", "*")

	d := websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	sendConn, _, err := d.Dial(newReq.URL.String(), newReq.Header)
	if err != nil {
		fmt.Println(err)
		return
	}

	elasticsearch.CreateIndex(elasticClient, "create-pvc", watchCreatePvcMapping)

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
		elasticsearch.CreateIndex(elasticClient, "create-pvc", watchCreatePvcMapping)

		_, msg, err := sendConn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		//解析watch并重新组装
		watchCreatePvcAssemble, flag := parseWatchCreatePvc(string(msg))
		if flag == 1 {
			continue
		}
		statusPhase := gjson.Get(watchCreatePvcAssemble, "status")
		if strings.EqualFold(string(statusPhase.Str), "Pending") {
			continue
		}
		watchType := gjson.Get(watchCreatePvcAssemble, "type")
		if strings.EqualFold(string(watchType.Str), "DELETED") {
			continue
		}

		// 控制台输出组装之后的watch pod
		//fmt.Println(watchCreatePvcAssemble)
		//fmt.Println(string(msg))

		// 数据入库
		_, err = elasticClient.Index().
			Index("create-pvc").
			Type("create-pvc").
			BodyString(watchCreatePvcAssemble).
			Do() //执行数据入索引库操作
		if err != nil {
			panic(err)
			return
		}
	}
	fmt.Println(newReq.Header)
}

/**
解析watch create pvc返回结果
*/
func parseWatchCreatePvc(watchCreatePvc string) (string, int) {
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
	watchPodResult := gjson.GetMany(watchCreatePvc, "type", "object.kind", "object.metadata.name",
		"object.metadata.namespace", "object.metadata.creationTimestamp", "object.status.phase",
		"object.status.capacity.storage", "object.metadata.deletionTimestamp", "object.spec.storageClassName", "object.spec.volumeName")

	return assembleWatchCreatePvc(watchPodResult)
}

/**
watch create pvc结果解析出的字段重新组装成Json格式
*/
func assembleWatchCreatePvc(a []gjson.Result) (string, int) {
	flag := 0
	// 如果object.metadata.deletionTimestamp有值表示是删除操作，会重复watch数据到es，所以标记为1，过滤掉
	if !strings.EqualFold(a[7].String(), "") {
		flag = 1
	}

	toBeChangeCreationTimestamp := a[4].String()

	//转化所需模板
	timeLayout := time.RFC3339 // 2006-01-02T15:04:05Z07:00

	// 使用模板在对应时区转化为time.time类型,go里面转换时间搓是10位，也就是转为的是秒不是毫秒
	// 字符串类型解析为时间类型
	creationTimestampParse, _ := time.ParseInLocation(timeLayout, toBeChangeCreationTimestamp, time.UTC)

	// 时间转为时间搓
	creationTimestampParseToTimestamp := creationTimestampParse.Unix() //这个转换时间搓默认(无论是time.Local还是time.UTC)会加上8个小时,why?,也好，这样就自动加了8小时

	// 时间搓转回时间
	// 设置时间搓，使用模板格式化为日期字符串,因为时间搓已经加了8小时，所以再转为字符串的话这个时间就和之前的多8小时了。
	creationTimestampAddEightHours := time.Unix(creationTimestampParseToTimestamp, 0).Format(timeLayout)
	systemCurrentTime := time.Unix(time.Now().Unix(), 0).Format(timeLayout)
	createPvcDuration := time.Now().Unix() - creationTimestampParseToTimestamp
	duration := utils.HoursMintuesSeconds(createPvcDuration)
	watchCreatePvcResult := `{"systemTimestamp":` + `"` + systemCurrentTime + `"` +
		",\"type\":" + `"` + a[0].String() + `"` +
		",\"kind\":" + `"` + a[1].String() + `"` +
		",\"name\":" + `"` + a[2].String() + `"` +
		",\"namespace\":" + `"` + a[3].String() + `"` +
		",\"storageClassName\":" + `"` + a[8].String() + `"` +
		",\"pvName\":" + `"` + a[9].String() + `"` +
		",\"creationTimestamp\":" + `"` + creationTimestampAddEightHours + `"` +
		",\"status\":" + `"` + a[5].String() + `"` +
		",\"capacity\":" + `"` + a[6].String() + `"` +
		",\"createPvcDurationInt\":" + `"` + strconv.FormatInt(createPvcDuration, 10) + `"` +
		",\"createPvcDuration\":" + `"` + duration + `"` +
		"}"

	fmt.Println("系统时间：", systemCurrentTime, ",pvc创建到分配成功的时间：", creationTimestampAddEightHours)
	return watchCreatePvcResult, flag
}
