package zlm

import (
	"context"
)

// OnStreamNoneReaderReq 表示 on_stream_none_reader 提交的数据
type OnStreamNoneReaderReq struct {
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
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamNoneReaderRes 表示 on_stream_none_reader 返回值
type OnStreamNoneReaderRes struct {
	Close *bool `json:"close"`
	Code  int   `json:"code"`
}

// OnStreamNoneReader 处理 zlm 的 on_stream_none_reader 回调
func OnStreamNoneReader(ctx context.Context, req *OnStreamNoneReaderReq, res *OnStreamNoneReaderRes) {
	res.Close = new(bool)
	*res.Close = true
}
