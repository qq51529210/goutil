package zlm

import (
	"context"
)

// DelStreamPusherProxyReq 是 DelStreamPusherProxy 的参数
type DelStreamPusherProxyReq struct {
	// addStreamPusherProxy 返回的 key
	Key string `query:"key"`
}

// delStreamPusherProxyRes 是 DelStreamPusherProxy 返回值
type delStreamPusherProxyRes struct {
	CodeMsg
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

// DelStreamPusherProxyResData 是 addStreamPusherProxyRes 的 Data 字段
type DelStreamPusherProxyResData struct {
	// 唯一标识
	Key string
}

const (
	DelStreamPusherProxyPath = apiPathPrefix + "/delStreamPusherProxy"
)

// DelStreamPusherProxy 调用 /index/api/delStreamPusherProxy 停止推流，可以使用 close_streams 替代
func DelStreamPusherProxy(ctx context.Context, ser Server, req *DelStreamPusherProxyReq) (bool, error) {
	// 请求
	var res delStreamPusherProxyRes
	if err := Request(ctx, ser, DelStreamPusherProxyPath, req, &res); err != nil {
		return false, err
	}
	// 经过测试，-500 应该是不存在的意思
	// 不存在也当它成功了
	if res.Code != CodeOK && res.Code != -500 {
		return false, &res.CodeMsg
	}
	//
	return res.Data.Flag, nil
}
