package zlm

import (
	"context"
)

// OnStreamNotFoundReq 表示 on_stream_not_found 提交的数据
type OnStreamNotFoundReq struct {
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
	// url 查询字符串
	Params string `json:"params"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamNotFound 处理 zlm 的 on_stream_not_found 回调
func OnStreamNotFound(ctx context.Context, req *OnStreamNotFoundReq, res *CodeMsg) {
}
