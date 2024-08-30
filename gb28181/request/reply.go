package request

import (
	gs "goutil/sync"
	"time"
)

var (
	replys gs.TimeoutContextPool
)

// Reply 用于有应答的请求，同步等待结果
type Reply gs.TimeoutContext

// AddReply 添加
func AddReply(deviceID, sn string, data any, timeout time.Duration) *gs.TimeoutContext {
	tx, _ := replys.New(deviceID+sn, data, timeout)
	return tx
}

// GetReply 获取
func GetReply(deviceID, sn string) *gs.TimeoutContext {
	return replys.Get(deviceID + sn)
}

// XMLResult 用于接收 xml.response.Result 字段的值
type XMLResult struct {
	Result string
}
