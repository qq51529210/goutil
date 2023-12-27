package zlm

import (
	"context"
)

// RestartServerReq 是 RestartServer 参数
type RestartServerReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
}

// restartServerRes 是 RestartServer 的返回值
type restartServerRes struct {
	apiError
}

const (
	apiRestartServer = "restartServer"
)

// RestartServer 调用 /index/api/restartServer
// 重启服务器,只有Daemon方式才能重启，否则是直接关闭！
func RestartServer(ctx context.Context, req *RestartServerReq) error {
	var res restartServerRes
	err := request[any](ctx, req.BaseURL, apiRestartServer, nil, &res)
	if err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiRestartServer
		return &res.apiError
	}
	return nil
}
