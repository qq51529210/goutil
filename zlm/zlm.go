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

// httpCallRes 封装请求
func httpCallRes[ReqQuery, ResData any](ctx context.Context, ser *Server, path string, query *ReqQuery, res *ResData) error {
	return httpCall(ctx, ser.Secret, ser.APIBaseURL, path, query, func(response *http.Response) error {
		// 必须是 200
		if response.StatusCode != http.StatusOK {
			return gh.StatusError(response.StatusCode)
		}
		// 解析
		return json.NewDecoder(response.Body).Decode(res)
	})
}

// httpCall 封装请求
func httpCall[Query any](ctx context.Context, secret, apiBaseURL, path string, query *Query, onRes func(res *http.Response) error) error {
	// 参数
	q := make(url.Values)
	q.Set(querySecret, secret)
	q.Set(queryVHost, VHost)
	if query != nil {
		q = gh.Query(query, q)
	}
	// 请求
	old := time.Now()
	url := fmt.Sprintf("%s/index/api/%s?%s", apiBaseURL, path, q.Encode())
	err := gh.Request[any](ctx, http.DefaultClient, http.MethodGet, url, nil, nil, onRes)
	if err != nil {
		return err
	}
	Logger.Debugf("[%v] %s", time.Since(old), url)
	return nil
}
