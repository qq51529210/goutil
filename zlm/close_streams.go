package zlm

import (
	"context"
)

// CloseStreamsReq 是 CloseStreams 参数
type CloseStreamsReq struct {
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 是否强制关闭
	Force Boolean `query:"force"`
}

// CloseStreamsRes 是 CloseStreams 的返回值
type CloseStreamsRes struct {
	CodeMsg
	// 命中的流个数
	CountHit int `json:"count_hit"`
	// 被关闭的流个数
	CountClosed int `json:"count_closed"`
}

const (
	CloseStreamsPath = apiPathPrefix + "/close_streams"
)

// CloseStreams 调用 /index/api/close_streams ，关闭流
// 经过测试 code=-500 应该是流不存在的意思
func CloseStreams(ctx context.Context, ser Server, req *CloseStreamsReq, res *CloseStreamsRes) error {
	return Request(ctx, ser, CloseStreamsPath, req, res)
}
