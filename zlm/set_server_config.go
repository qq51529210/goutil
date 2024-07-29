package zlm

import (
	"context"
	"encoding/json"
	ghttp "goutil/http"
	"net/http"
)

// SetServerConfigReq 是 SetServerConfig 参数
type SetServerConfigReq struct {
	// 数据
	Data map[string]string
}

// setServerConfigRes 是 SetServerConfig 的返回值
type setServerConfigRes struct {
	CodeMsg
	Changed int `json:"changed"`
}

const (
	SetServerConfigPath = apiPathPrefix + "/setServerConfig"
)

// SetServerConfig 调用 /index/api/setServerConfig ，更新配置
func SetServerConfig(ctx context.Context, ser Server, req *SetServerConfigReq) error {
	// 参数
	q := initRequestQuery(ser)
	for k, v := range req.Data {
		q.Add(k, v)
	}
	// 请求
	var res setServerConfigRes
	if err := ghttp.JSONRequest(ctx, http.DefaultClient, http.MethodGet,
		ser.BaseURL()+SetServerConfigPath, q, nil,
		func(res *http.Response) error {
			// 必须是 200
			if res.StatusCode != http.StatusOK {
				return ghttp.StatusError(res.StatusCode)
			}
			// 解析
			return json.NewDecoder(res.Body).Decode(&res)
		}); err != nil {
		return err
	}
	if res.Code != CodeOK {
		return &res.CodeMsg
	}
	return nil
}

// SetServerConfig2 调用 /index/api/setServerConfig ，更新配置
func SetServerConfig2[M any](ctx context.Context, ser Server, m M) error {
	// 请求
	var res setServerConfigRes
	if err := Request(ctx, ser, SetServerConfigPath, &m, &res); err != nil {
		return err
	}
	if res.Code != CodeOK {
		return &res.CodeMsg
	}
	return nil
}
