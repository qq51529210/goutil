package zlm

import (
	"context"
	"errors"
)

var (
	// ErrConfig 错误的服务配置
	ErrConfig = errors.New("error config")
)

// getServerConfigRes 是 GetServerConfig 的返回值
type getServerConfigRes struct {
	apiError
	Data []map[string]string `json:"data"`
}

const (
	apiGetServerConfig = "getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig
func (s *Server) GetServerConfig(ctx context.Context) (map[string]string, error) {
	// 请求
	var res getServerConfigRes
	err := httpCallRes[any](ctx, s, apiGetServerConfig, nil, &res)
	if err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiGetServerConfig
		return nil, &res.apiError
	}
	// 找到自己的配置
	for _, d := range res.Data {
		if d["general.mediaServerId"] == s.ID {
			return d, nil
		}
	}
	// 没有配置，流媒体服务有问题
	return nil, ErrConfig
}

// getServerConfigRes2 是 GetServerConfig 的返回值
type getServerConfigRes2[M any] struct {
	apiError
	Data []M `json:"data"`
}

// GetServerConfig 调用 /index/api/getServerConfig
// 它这个返回确实是一个数组，不知道为什么
func GetServerConfig[M any](ctx context.Context, s *Server) ([]M, error) {
	// 请求
	var res getServerConfigRes2[M]
	err := httpCallRes[any](ctx, s, apiGetServerConfig, nil, &res)
	if err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiGetServerConfig
		return nil, &res.apiError
	}
	return res.Data, nil
}
