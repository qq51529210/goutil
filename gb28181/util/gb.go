package util

import (
	"time"
)

const (
	// GBTimeForamt 国标时间格式
	GBTimeForamt = "2006-01-02T15:04:05"
)

// GBTime 返回国标格式的时间字符串
func GBTime() string {
	return time.Now().Format(GBTimeForamt)
}

// IsGBTime 验证国标时间
func IsGBTime(timeStr string) bool {
	_, err := time.Parse(GBTimeForamt, timeStr)
	return err == nil
}

// IsGBID 验证国标编号
func IsGBID(id string) bool {
	if len(id) != 20 {
		return false
	}
	for _, n := range id {
		if n < '0' || n > '9' {
			return false
		}
	}
	return true
}

// GBTimestamp 解析并返回时间戳
func GBTimestamp(t string) int64 {
	_t, err := time.ParseInLocation(GBTimeForamt, t, time.Local)
	if err != nil {
		return 0
	}
	return _t.Unix()
}
