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
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 服务标识，用于筛选配置
	ID string
}

// getServerConfigRes 是 GetServerConfig 的返回值
type getServerConfigRes struct {
	apiError
	Data []map[string]string `json:"data"`
}

const (
	apiGetServerConfig = "getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig ，返回配置
func GetServerConfig(ctx context.Context, req *GetServerConfigReq) (map[string]string, error) {
	// 请求
	var res getServerConfigRes
	if err := request(ctx, req.BaseURL, apiGetServerConfig, req, &res); err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiGetServerConfig
		return nil, &res.apiError
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
type getServerConfigAndUnmarshalRes[M any] struct {
	apiError
	Data []M `json:"data"`
}

// GetServerConfigAndUnmarshal 调用 /index/api/getServerConfig ，返回配置
func GetServerConfigAndUnmarshal[M any](ctx context.Context, req *GetServerConfigReq) ([]M, error) {
	// 请求
	var res getServerConfigAndUnmarshalRes[M]
	if err := request(ctx, req.BaseURL, apiGetServerConfig, req, &res); err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiGetServerConfig
		return nil, &res.apiError
	}
	return res.Data, nil
}
