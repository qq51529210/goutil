package util

import (
	"gbs/validate"
	"time"
)

// GBTime 返回国标格式的时间字符串
func GBTime() string {
	return time.Now().Format(validate.GBTimeForamt)
}

// GBTimestamp 解析并返回时间戳
func GBTimestamp(t string) int64 {
	_t, err := time.ParseInLocation(validate.GBTimeForamt, t, time.Local)
	if err != nil {
		return 0
	}
	return _t.Unix()
}
