package zlm

import (
	"context"
)

// RestartServerReq 是 RestartServer 参数
type RestartServerReq struct {
}

// RestartServerRes 是 RestartServer 返回值
type RestartServerRes struct {
	CodeMsg
}

const (
	RestartServerPath = apiPathPrefix + "/restartServer"
)

// RestartServer 调用 /index/api/restartServer ，重启服务器
func RestartServer(ctx context.Context, ser Server, req *RestartServerReq, res *RestartServerRes) error {
	return Request(ctx, ser, RestartServerPath, req, res)
}
