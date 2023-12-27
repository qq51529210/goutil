package zlm

import (
	"context"
	"encoding/json"
	"fmt"
	gh "goutil/http"
	"goutil/log"
	"net/http"
	"net/url"
	"time"
)

// 参数常量
const (
	True  = "1"
	False = "0"
)

// 流的协议
const (
	RTMP = "rtmp"
	RTSP = "rtsp"
	HLS  = "hls"
	TS   = "ts"
	FMP4 = "fmp4"
)

const (
	// VHost 默认的 vhost
	VHost = "__defaultVhost__"
	// 查询参数的名称
	queryVHost  = "vhost"
	querySecret = "secret"
)

var (
	// Logger 用于打印
	Logger *log.Logger = log.DefaultLogger
)

type apiCall struct {
	// 标识
	ID string
	// http://localhost:8080
	BaseURL string
	// 配置的密钥
	Secret string
}

func (m *apiCall) url(path string, query any) string {
	// 参数
	q := make(url.Values)
	q.Set(querySecret, m.Secret)
	q.Set(queryVHost, VHost)
	if query != nil {
		q = gh.Query(query, q)
	}
	return fmt.Sprintf("%s/index/api/%s?%s", m.BaseURL, path, q.Encode())
}

// request 封装请求
func request[Query, Response any](ctx context.Context, call *apiCall, path string, query *Query, data *Response) error {
	url := call.url(path, query)
	// 请求
	old := time.Now()
	err := gh.Request[any](ctx, http.DefaultClient, http.MethodGet, url, nil, nil,
		func(res *http.Response) error {
			// 必须是 200
			if res.StatusCode != http.StatusOK {
				return gh.StatusError(res.StatusCode)
			}
			// 解析
			return json.NewDecoder(res.Body).Decode(data)
		})
	if err != nil {
		return err
	}
	Logger.Debugf("[%v] %s", time.Since(old), url)
	return nil
}
