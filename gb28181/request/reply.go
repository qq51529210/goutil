package request

import (
	gs "goutil/sync"
	"time"
)

var (
	replys gs.Map[string, *Reply]
)

func init() {
	replys.Init()
}

// Reply 用于有应答的请求
// 实现 context.Context 同步等待结果
type Reply struct {
	gs.Context
	// 池的 key
	key string
}

func (m *Reply) onFinish() {
	replys.Del(m.key)
}

// AddReply 添加
func AddReply(deviceID, sn string, data any, timeout time.Duration) *Reply {
	m := new(Reply)
	m.key = deviceID + sn
	m.Context.OnFinish = m.onFinish
	m.Context.Run(data, timeout)
	//
	replys.Set(m.key, m)
	//
	return m
}

// GetReply 获取
func GetReply(deviceID, sn string) *Reply {
	key := deviceID + sn
	return replys.Get(key)
}
