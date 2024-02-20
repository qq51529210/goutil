package zlm

import (
	"context"
)

// CloseStreamsReq 是 CloseStreams 参数
type CloseStreamsReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 是否强制关闭(有人在观看是否还关闭)，0/1
	Force string `query:"force"`
}

// closeStreamsRes 封装 CloseStreamsRes
type closeStreamsRes struct {
	apiError
	CloseStreamsRes
}

// CloseStreamsRes 是 closeStreams 返回值
type CloseStreamsRes struct {
	// 筛选命中的流个数
	CountHit int `json:"count_hit"`
	// 被关闭的流个数，可能小于count_hit
	CountClosed int `json:"count_closed"`
}

const (
	apiCloseStreams = "close_streams"
)

// CloseStreams 调用 /index/api/close_streams
// 关闭流(目前所有类型的流都支持关闭)
func CloseStreams(ctx context.Context, req *CloseStreamsReq) (*CloseStreamsRes, error) {
	// 请求
	req.Force = True
	var res closeStreamsRes
	if err := request(ctx, req.BaseURL, apiCloseStreams, req, &res); err != nil {
		return nil, err
	}
	// 经过测试，-500 应该是不存在的意思
	if res.apiError.Code != codeTrue && res.Code != -500 {
		res.apiError.Path = apiCloseStreams
		return nil, &res.apiError
	}
	//
	return &res.CloseStreamsRes, nil
}
