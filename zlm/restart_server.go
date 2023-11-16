package zlm

import (
	"context"
)

// restartServerRes 是 RestartServer 的返回值
type restartServerRes struct {
	apiError
}

const (
	apiRestartServer = "restartServer"
)

// RestartServer 调用 /index/api/restartServer
// 重启服务器,只有Daemon方式才能重启，否则是直接关闭！
func (s *Server) RestartServer(ctx context.Context) error {
	var res restartServerRes
	err := httpCallRes[any](ctx, s, apiRestartServer, nil, &res)
	if err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiRestartServer
		return &res.apiError
	}
	return nil
}
