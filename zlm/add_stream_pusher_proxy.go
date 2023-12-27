package zlm

import (
	"context"
)

// AddStreamPusherProxyReq 是 AddStreamPusherProxy 参数
type AddStreamPusherProxyReq struct {
	apiCall
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 推流地址，需要与schema字段协议一致
	DstURL string `query:"dst_url"`
	// rtsp推流时，拉流方式，0：tcp，1：udp
	RTPType string `query:"rtp_type"`
	// 推流超时时间，单位秒，float类型
	TimeoutSec string `query:"timeout_sec"`
	// 推流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string `query:"retry_count"`
}

// addStreamPusherProxyRes 是 AddStreamPusherProxy 返回值
type addStreamPusherProxyRes struct {
	apiError
	Data AddStreamPusherProxyResData `json:"data"`
}

// AddStreamPusherProxyResData 是 addStreamPusherProxyRes 的 Data 字段
type AddStreamPusherProxyResData struct {
	// 流的唯一标识
	Key string
}

const (
	apiAddStreamPusherProxy = "addStreamPusherProxy"
)

// AddStreamPusherProxy 调用 /index/api/addStreamPusherProxy ，返回 key
func AddStreamPusherProxy(ctx context.Context, req *AddStreamPusherProxyReq) (string, error) {
	// 请求
	var res addStreamPusherProxyRes
	err := request(ctx, &req.apiCall, apiAddStreamPusherProxy, req, &res)
	if err != nil {
		return "", err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = req.apiCall.ID
		res.apiError.Path = apiAddStreamPusherProxy
		return "", &res.apiError
	}
	return res.Data.Key, nil
}
