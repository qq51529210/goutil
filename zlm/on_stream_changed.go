package zlm

import (
	"context"
)

// OnStreamChangedReq 表示 on_stream_changed 提交的数据，保存注册和注销的所有字段
type OnStreamChangedReq struct {
	// 服务器id，通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流注册或注销
	Regist bool `json:"regist"`
	// 流的媒体信息
	MediaListData
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamChanged 处理 zlm 的 on_stream_changed 回调
func OnStreamChanged(ctx context.Context, req *OnStreamChangedReq, res *CodeMsg) {
}
