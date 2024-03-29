package zlm

import (
	"context"
)

// DelStreamPusherProxyReq 是 DelStreamPusherProxy 的参数
type DelStreamPusherProxyReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 流的唯一标识
	Key string `query:"key"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
}

// delStreamPusherProxyRes 是 DelStreamPusherProxy 返回值
type delStreamPusherProxyRes struct {
	apiError
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
	apiDelStreamPusherProxy = "delStreamPusherProxy"
)

// DelStreamPusherProxy 调用 /index/api/delStreamPusherProxy
func DelStreamPusherProxy(ctx context.Context, req *DelStreamPusherProxyReq) (bool, error) {
	// 请求
	var res delStreamPusherProxyRes
	if err := request(ctx, req.BaseURL, apiDelStreamPusherProxy, req, &res); err != nil {
		return false, err
	}
	if res.apiError.Code != codeTrue {
		// -500 是没有找到流，也算成功
		if res.apiError.Code != -500 {
			res.apiError.Path = apiDelStreamPusherProxy
			return false, &res.apiError
		}
	}
	return res.Data.Flag, nil
}
