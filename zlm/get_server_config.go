package zlm

import (
	"context"
	"errors"
)

var (
	// ErrConfig 错误的服务配置
	ErrConfig = errors.New("error config")
)

// GetServerConfigReq 是 GetServerConfig 参数
type GetServerConfigReq struct {
	// 服务标识，用于筛选配置
	ID string
}

// getServerConfigRes 是 GetServerConfig 的返回值
type getServerConfigRes struct {
	CodeMsg
	Data []map[string]string `json:"data"`
}

const (
	GetServerConfigPath = apiPathPrefix + "/getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig ，返回配置
func GetServerConfig(ctx context.Context, ser Server, req *GetServerConfigReq) (map[string]string, error) {
	// 请求
	var res getServerConfigRes
	if err := Request(ctx, ser, GetServerConfigPath, req, &res); err != nil {
		return nil, err
	}
	if res.Code != CodeOK {
		return nil, &res.CodeMsg
	}
	// 找到自己的配置
	for _, d := range res.Data {
		if d["general.mediaServerId"] == req.ID {
			return d, nil
		}
	}
	// 没有配置，流媒体服务有问题
	return nil, ErrConfig
}

// getServerConfigAndUnmarshalRes 是 GetServerConfigAndUnmarshal 的返回值
type getServerConfigAndUnmarshalRes[M ServerID] struct {
	CodeMsg
	Data []M `json:"data"`
}

// GetServerConfigAndUnmarshal 调用 /index/api/getServerConfig ，返回配置
func GetServerConfigAndUnmarshal[M ServerID](ctx context.Context, ser Server, req *GetServerConfigReq) (m M, err error) {
	// 请求
	var res getServerConfigAndUnmarshalRes[M]
	if err = Request(ctx, ser, GetServerConfigPath, req, &res); err != nil {
		return
	}
	if res.Code != CodeOK {
		err = &res.CodeMsg
		return
	}
	// 筛选
	for i := 0; i < len(res.Data); i++ {
		if res.Data[i].ID() == req.ID {
			m = res.Data[i]
			return
		}
	}
	return
}
