package zlm

import (
	"context"
)

// OnPlayReq 表示 on_play 提交的数据
type OnPlayReq struct {
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

// OnPlay 处理 zlm 的 on_play 回调
// 检查 sid 和 skey ，验证是否集群内存同步
// 否则检查播放 token
// 如果本服务有流，就直接播放了
// 如果没有流，等 on_stream_not_found 再处理同步问题
func OnPlay(ctx context.Context, req *OnPlayReq, res *CodeMsg) {
}
