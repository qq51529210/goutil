package zlm

import (
	"context"
	"encoding/json"
	"fmt"
	gh "goutil/http"
	"net/http"
	"net/url"
	"time"
)

// SetServerConfigReq 是 SetServerConfig 参数
type SetServerConfigReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 数据
	Data map[string]string
}

// setServerConfigRes 是 SetServerConfig 的返回值
type setServerConfigRes struct {
	apiError
	Changed int `json:"changed"`
}

const (
	apiSetServerConfig = "setServerConfig"
)

// SetServerConfig 调用 /index/api/setServerConfig ，返回配置
func SetServerConfig(ctx context.Context, req *SetServerConfigReq) error {
	q := make(url.Values)
	for k, v := range req.Data {
		q.Add(k, v)
	}
	q = gh.Query(req, q)
	url := fmt.Sprintf("%s?%s", requestURL[any](req.BaseURL, apiSetServerConfig, nil), q.Encode())
	// 请求
	var res setServerConfigRes
	old := time.Now()
	err := gh.JSONRequest(ctx, http.DefaultClient, http.MethodGet, url, nil, nil,
		func(res *http.Response) error {
			// 必须是 200
			if res.StatusCode != http.StatusOK {
				return gh.StatusError(res.StatusCode)
			}
			// 解析
			return json.NewDecoder(res.Body).Decode(&res)
		})
	if err != nil {
		return err
	}
	Logger.DebugfDepth(2, "[%v] %s", time.Since(old), url)
	//
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiSetServerConfig
		return &res.apiError
	}
	return nil
}

// setServerConfigAndUnmarshalRes 是 SetServerConfigAndUnmarshal 的返回值
type setServerConfigAndUnmarshalRes[M any] struct {
	apiError
	Data []M `json:"data"`
}

// SetServerConfigAndUnmarshal 调用 /index/api/setServerConfig ，返回配置
func SetServerConfigAndUnmarshal[M any](ctx context.Context, req *SetServerConfigReq) ([]M, error) {
	// 请求
	var res setServerConfigAndUnmarshalRes[M]
	if err := request(ctx, req.BaseURL, apiSetServerConfig, req, &res); err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiSetServerConfig
		return nil, &res.apiError
	}
	return res.Data, nil
}
