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

// AddStreamPusherProxyRes 是 AddStreamPusherProxy 返回值
type AddStreamPusherProxyRes struct {
	CodeMsg
	Data struct {
		// 流的唯一标识
		Key string
	} `json:"data"`
}

const (
	AddStreamPusherProxyPath = apiPathPrefix + "/addStreamPusherProxy"
)

// AddStreamPusherProxy 调用 /index/api/addStreamPusherProxy ，把本服务器的直播流（rtsp/rtmp）推送到其他服务器
func AddStreamPusherProxy(ctx context.Context, ser Server, req *AddStreamPusherProxyReq, res *AddStreamPusherProxyRes) error {
	return Request(ctx, ser, AddStreamPusherProxyPath, req, res)
}
