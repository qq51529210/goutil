package zlm

import (
	"context"
)

// CloseRTPServerReq 是 CloseRTPServer 的参数
type CloseRTPServerReq struct {
	apiCall
	// 调用closeRtpServer接口时提供的流ID
	Stream string `query:"stream_id"`
}

// closeRTPServerRes 是 CloseRTPServer 的返回值
type closeRTPServerRes struct {
	apiError
	// 是否找到记录并关闭
	Hit int `json:"hit"`
}

const (
	apiCloseRtpServer = "closeRtpServer"
)

// CloseRTPServer 调用 /index/api/closeRtpServer
// 关闭GB28181 RTP接收端口
// 返回成功的个数
func CloseRTPServer(ctx context.Context, req *CloseRTPServerReq) (int, error) {
	// 请求
	var res closeRTPServerRes
	err := request(ctx, &req.apiCall, apiCloseRtpServer, req, &res)
	if err != nil {
		return 0, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = req.apiCall.ID
		res.apiError.Path = apiCloseRtpServer
		return 0, &res.apiError
	}
	return res.Hit, nil
}
