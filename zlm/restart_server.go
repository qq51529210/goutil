package zlm

import (
	"context"
)

// RestartServerReq 是 RestartServer 参数
type RestartServerReq struct {
}

const (
	RestartServerPath = apiPathPrefix + "/restartServer"
)

// RestartServer 调用 /index/api/restartServer ，重启服务器
func RestartServer(ctx context.Context, ser Server, req *RestartServerReq) error {
	var res CodeMsg
	if err := Request(ctx, ser, RestartServerPath, req, &res); err != nil {
		return err
	}
	if res.Code != CodeOK {
		return &res
	}
	return nil
}
