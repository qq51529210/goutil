package gb28181

import (
	"time"
)

const (
	// TimeForamt 国标时间格式
	TimeForamt = "2006-01-02T15:04:05"
)

// Time 返回国标格式的时间字符串
func Time() string {
	return time.Now().Format(TimeForamt)
}

// IsTime 验证 t 是否国标时间格式
func IsTime(t string) bool {
	_, err := time.Parse(TimeForamt, t)
	return err == nil
}

// Timestamp 解析标时间格式的 t 并返回时间戳
func Timestamp(t string) int64 {
	_t, err := time.ParseInLocation(TimeForamt, t, time.Local)
	if err != nil {
		return 0
	}
	return _t.Unix()
}

// TimeFromTimestamp 返回时间戳 ts 的国标格式时间字符串
func TimeFromTimestamp(ts int64) string {
	return time.Unix(ts, 0).Format(TimeForamt)
}

// CheckID 验证 id 是否国标编号
func CheckID(id string) bool {
	if len(id) != 20 {
		return false
	}
	return IsNumber(id)
}

// IsNumber 验证 id 是否全数字
func IsNumber(id string) bool {
	for _, n := range id {
		if n < '0' || n > '9' {
			return false
		}
	}
	return true
}
