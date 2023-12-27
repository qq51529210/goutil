package zlm

import (
	"context"
)

// RestartServerReq 是 RestartServer 参数
type RestartServerReq struct {
	apiCall
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
	err := request[any](ctx, &req.apiCall, apiRestartServer, nil, &res)
	if err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = req.apiCall.ID
		res.apiError.Path = apiRestartServer
		return &res.apiError
	}
	return nil
}
