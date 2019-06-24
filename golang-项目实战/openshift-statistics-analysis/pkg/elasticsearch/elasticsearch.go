package elasticsearch

import (
	"fmt"

	 "gopkg.in/olivere/elastic.v3"
)

/**
创建索引
*/
func CreateIndex(elasticClient *elastic.Client, indexName string, indexMapping string) {
	existsIndex, err := elasticClient.IndexExists(indexName).Do()
	if err != nil {
		panic(err)
	}
	if !existsIndex {
		createIndex, err := elasticClient.CreateIndex(indexName).
			BodyString(indexMapping).
			Do()
		if err != nil {
			panic(err)
			return
		}
		if !createIndex.Acknowledged {
			// Not Acknowledged
			fmt.Println("索引还未创建，未确认")
		}
	}
}
