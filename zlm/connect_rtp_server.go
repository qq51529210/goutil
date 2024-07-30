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

// ConnectRTPServerRes 是 ConnectRTPServer 的返回值
type ConnectRTPServerRes struct {
	CodeMsg
}

const (
	ConnectRTPServerPath = apiPathPrefix + "/connectRtpServer"
)

// ConnectRTPServer 调用 /index/api/connectRtpServer ，用于主动模式拉流
func ConnectRTPServer(ctx context.Context, ser Server, req *ConnectRTPServerReq, res *ConnectRTPServerRes) error {
	return Request(ctx, ser, ConnectRTPServerPath, req, res)
}
