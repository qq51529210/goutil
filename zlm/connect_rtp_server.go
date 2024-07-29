package zlm

import (
	"context"
)

// ConnectRTPServerReq 是 ConnectRTPServer 的参数
type ConnectRTPServerReq struct {
	// 服务端地址
	DstIP string `query:"dst_url"`
	// 服务端端口
	DstPort string `query:"dst_port"`
	// OpenRtpServer 时绑定的流标识
	Stream string `query:"stream_id"`
}

const (
	ConnectRTPServerPath = apiPathPrefix + "/connectRtpServer"
)

// ConnectRTPServer 调用 /index/api/connectRtpServer ，用于主动模式拉流
func ConnectRTPServer(ctx context.Context, ser Server, req *ConnectRTPServerReq) error {
	// 请求
	var res CodeMsg
	if err := Request(ctx, ser, ConnectRTPServerPath, req, &res); err != nil {
		return err
	}
	if res.Code != CodeOK {
		return &res
	}
	return nil
}
