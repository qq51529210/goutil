package zlm

import (
	"context"
)

// SetServerConfigReq 是 SetServerConfig 参数
// M 是结构体
//
//	type config struct {
//		AliveInterval      string `query:"hook.alive_interval"`
//	}
type SetServerConfigReq[M any] struct {
	// 数据
	Data M
}

// SetServerConfigRes 是 SetServerConfig 的返回值
type SetServerConfigRes struct {
	CodeMsg
	Changed int `json:"changed"`
}

const (
	SetServerConfigPath = apiPathPrefix + "/setServerConfig"
)

// SetServerConfig 调用 /index/api/setServerConfig ，更新配置
func SetServerConfig[M Config](ctx context.Context, ser Server, req *SetServerConfigReq[M], res *SetServerConfigRes) error {
	return Request(ctx, ser, SetServerConfigPath, req.Data, res)
}
