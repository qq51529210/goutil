package zlm

import (
	"context"
	"encoding/json"
	"fmt"
	gh "goutil/http"
	"goutil/log"
	"net/http"
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
)

var (
	// Logger 用于打印
	Logger *log.Logger = log.DefaultLogger
)

func requestURL[Query any](baseURL, path string, query *Query) string {
	// 参数
	if query != nil {
		return fmt.Sprintf("%s/index/api/%s?%s", baseURL, path, gh.Query(query, nil).Encode())
	}
	return fmt.Sprintf("%s/index/api/%s", baseURL, path)
}

// request 封装请求
func request[Query, Response any](ctx context.Context, baseURL, path string, query *Query, data *Response) error {
	url := requestURL(baseURL, path, query)
	// 请求
	old := time.Now()
	err := gh.JSONRequest(ctx, http.DefaultClient, http.MethodGet, url, nil, nil,
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
	Logger.DebugfDepth(2, "[%v] %s", time.Since(old), url)
	return nil
}
