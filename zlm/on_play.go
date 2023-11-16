package zlm

import (
	"context"
)

// OnPlayReq 表示 on_play 提交的数据
type OnPlayReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 推流的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 推流url参数
	Params string `json:"params"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnPlayRes 表示 on_play 的返回值
type OnPlayRes struct {
	apiError
}

// OnPlay 处理 zlm 的 on_play 回调
// 检查 sid 和 skey ，验证是否集群内存同步
// 否则检查播放 token
// 如果本服务有流，就直接播放了
// 如果没有流，等 on_stream_not_found 再处理同步问题
func OnPlay(ctx context.Context, req *OnPlayReq, res *OnPlayRes) {

}
