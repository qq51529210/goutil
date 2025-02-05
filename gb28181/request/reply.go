package request

import (
	gsync "goutil/sync"
	"time"
)

var (
	replys = gsync.NewTimeoutContextPool()
)

// AddReply 添加
func AddReply(deviceID, sn string, data any, timeout time.Duration) *gsync.TimeoutContext {
	tx, _ := replys.New(deviceID+sn, data, timeout)
	return tx
}

// GetReply 获取
func GetReply(deviceID, sn string) *gsync.TimeoutContext {
	return replys.Get(deviceID + sn)
}

// XMLResult 用于接收 xml.response.Result 字段的值
type XMLResult struct {
	Result string
}
