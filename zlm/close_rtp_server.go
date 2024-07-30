package zlm

import (
	"context"
)

// CloseRTPServerReq 是 CloseRTPServer 的参数
type CloseRTPServerReq struct {
	// 调用 openRtpServer 接口时提供的流ID
	Stream string `query:"stream_id"`
}

// CloseRTPServerRes 是 CloseRTPServer 的返回值
type CloseRTPServerRes struct {
	CodeMsg
	// 是否找到记录并关闭
	Hit int `json:"hit"`
}

const (
	CloseRtpServerPath = apiPathPrefix + "/closeRtpServer"
)

// CloseRTPServer 调用 /index/api/closeRtpServer ，关闭GB28181 RTP接收端口，返回成功的个数
func CloseRTPServer(ctx context.Context, ser Server, req *CloseRTPServerReq, res *CloseRTPServerRes) error {
	return Request(ctx, ser, CloseRtpServerPath, req, res)
}
