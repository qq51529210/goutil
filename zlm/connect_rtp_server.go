package zlm

import (
	"context"
)

// ConnectRTPServerReq 是 ConnectRTPServer 的参数
type ConnectRTPServerReq struct {
	// tcp主动模式时服务端地址
	DstURL string `query:"dst_url"`
	// tcp主动模式时服务端端口
	DstPort string `query:"dst_port"`
	// OpenRtpServer时绑定的流id
	Stream string `query:"stream_id"`
}

// connectRTPServerRes 是 OpenRTPServer 的返回值
type connectRTPServerRes struct {
	apiError
}

const (
	apiConnectRTPServer = "connectRtpServer"
)

// ConnectRTPServer 调用 /index/api/connectRtpServer
// 未找到文档说明，从 postman 上发现的
func (s *Server) ConnectRTPServer(ctx context.Context, req *ConnectRTPServerReq) error {
	// 请求
	var res connectRTPServerRes
	err := httpCallRes(ctx, s, apiConnectRTPServer, req, &res)
	if err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiConnectRTPServer
		return &res.apiError
	}
	return nil
}
