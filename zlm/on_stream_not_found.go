package zlm

import (
	"context"
)

// OnStreamNotFoundReq 表示 on_stream_not_found 提交的数据
type OnStreamNotFoundReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 播放的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 播放url参数
	Params string `json:"params"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamNotFound 处理 zlm 的 on_stream_not_found 回调
func OnStreamNotFound(ctx context.Context, req *OnStreamNotFoundReq, res *CodeMsg) {
}
