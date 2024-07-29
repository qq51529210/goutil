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

// closeStreamsRes 封装 CloseStreamsRes
type closeStreamsRes struct {
	CodeMsg
	CloseStreamsRes
}

// CloseStreamsRes 是 closeStreams 返回值
type CloseStreamsRes struct {
	// 命中的流个数
	CountHit int `json:"count_hit"`
	// 被关闭的流个数
	CountClosed int `json:"count_closed"`
}

const (
	CloseStreamsPath = apiPathPrefix + "/close_streams"
)

// CloseStreams 调用 /index/api/close_streams ，关闭流
func CloseStreams(ctx context.Context, ser Server, req *CloseStreamsReq) (*CloseStreamsRes, error) {
	// 请求
	var res closeStreamsRes
	if err := Request(ctx, ser, CloseStreamsPath, req, &res); err != nil {
		return nil, err
	}
	// 经过测试，-500 应该是不存在的意思
	// 不存在也当它成功了
	if res.Code != CodeOK && res.Code != -500 {
		return nil, &res.CodeMsg
	}
	//
	return &res.CloseStreamsRes, nil
}
