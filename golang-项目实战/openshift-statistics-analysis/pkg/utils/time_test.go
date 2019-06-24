package utils

import (
	"fmt"
	"testing"
)

//go test是go语言自带的测试工具，其中包含的是两类，单元测试和性能测试

func Test_HoursMintuesSeconds(t *testing.T) {
	myTime := HoursMintuesSeconds(35)
	fmt.Println("myTime:", myTime)
}
