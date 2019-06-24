package utils

import (
	"strconv"
)

// 输入多少秒转为几时几分几秒
func HoursMintuesSeconds(sec int64) string {
	var hours, mintues, seconds int64

	if sec > 3600 {
		sec = sec / 1000000000
	}

	hours = sec / 3600 //小时
	sec = sec % 3600   //剩下秒数
	mintues = sec / 60 //分钟
	seconds = sec % 60 //剩下秒数

	if strconv.FormatInt(hours, 10) == "0" {
		if strconv.FormatInt(mintues, 10) == "0" {
			return strconv.FormatInt(seconds, 10) + "秒"
		}

		return strconv.FormatInt(mintues, 10) + "分" + strconv.FormatInt(seconds, 10) + "秒"
	}
	return strconv.FormatInt(hours, 10) + "时" + strconv.FormatInt(mintues, 10) + "分" + strconv.FormatInt(seconds, 10) + "秒"
}
