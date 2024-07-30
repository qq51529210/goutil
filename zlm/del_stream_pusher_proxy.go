package zlm

import (
	"context"
)

// DelStreamPusherProxyReq 是 DelStreamPusherProxy 的参数
type DelStreamPusherProxyReq struct {
	// addStreamPusherProxy 返回的 key
	Key string `query:"key"`
}

// DelStreamPusherProxyRes 是 DelStreamPusherProxy 返回值
type DelStreamPusherProxyRes struct {
	CodeMsg
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

const (
	DelStreamPusherProxyPath = apiPathPrefix + "/delStreamPusherProxy"
)

// DelStreamPusherProxy 调用 /index/api/delStreamPusherProxy 停止推流，可以使用 close_streams 替代
// 经过测试 code=-500 应该是流不存在的意思
func DelStreamPusherProxy(ctx context.Context, ser Server, req *DelStreamPusherProxyReq, res *DelStreamPusherProxyRes) error {
	return Request(ctx, ser, DelStreamPusherProxyPath, req, res)
}
