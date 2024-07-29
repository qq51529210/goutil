package zlm

import (
	"context"
)

// AddStreamPusherProxyReq 是 AddStreamPusherProxy 参数
type AddStreamPusherProxyReq struct {
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 推流地址，需要与 schema 字段协议一致
	DstURL string `query:"dst_url"`
	// rtsp 拉流方式
	RTPType RTSPRTPType `query:"rtp_type"`
	// 推流超时时间，单位秒，float类型
	Timeout string `query:"timeout_sec"`
	// 推流重试次数，不传则无限重试
	RetryCount string `query:"retry_count"`
}

// addStreamPusherProxyRes 是 AddStreamPusherProxy 返回值
type addStreamPusherProxyRes struct {
	CodeMsg
	Data AddStreamPusherProxyResData `json:"data"`
}

// AddStreamPusherProxyResData 是 addStreamPusherProxyRes 的 Data 字段
type AddStreamPusherProxyResData struct {
	// 流的唯一标识
	Key string
}

const (
	AddStreamPusherProxyPath = apiPathPrefix + "/addStreamPusherProxy"
)

// AddStreamPusherProxy 调用 /index/api/addStreamPusherProxy ，把本服务器的直播流（rtsp/rtmp）推送到其他服务器，返回 key
func AddStreamPusherProxy(ctx context.Context, ser Server, req *AddStreamPusherProxyReq) (string, error) {
	// 请求
	var res addStreamPusherProxyRes
	if err := Request(ctx, ser, AddStreamPusherProxyPath, req, &res); err != nil {
		return "", err
	}
	if res.Code != CodeOK {
		return "", &res.CodeMsg
	}
	return res.Data.Key, nil
}
