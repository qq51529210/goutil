package zlm

import (
	"context"
)

// OnStreamNoneReaderReq 表示 on_stream_none_reader 提交的数据
type OnStreamNoneReaderReq struct {
	// 虚拟主机
	VHost string `json:"vhost"`
	// 服务标识
	MediaServerID string `json:"mediaServerId"`
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamNoneReaderRes 表示 on_stream_none_reader 返回值
type OnStreamNoneReaderRes struct {
	Close bool `json:"close,omitempty"`
	Code  int  `json:"code"`
}

// OnStreamNoneReader 处理 zlm 的 on_stream_none_reader 回调
func OnStreamNoneReader(ctx context.Context, req *OnStreamNoneReaderReq, res *OnStreamNoneReaderRes) {
	res.Close = true
}
