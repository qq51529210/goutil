package zlm

import (
	"context"
)

// StopSendRTPReq 是 StopSendRTP 参数
type StopSendRTPReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 根据ssrc关停某路rtp推流，不传时关闭所有推流
	SSRC string `query:"ssrc"`
}

// stopSendRTPRes 是 StopSendRTP 返回值
type stopSendRTPRes struct {
	apiError
}

const (
	apiStopSendRTP = "stopSendRtp"
)

// StopSendRTP 调用 /index/api/stopSendRtp
// 停止 GB28181 rtp 推流。
func StopSendRTP(ctx context.Context, req *StopSendRTPReq) error {
	var res stopSendRTPRes
	if err := request(ctx, req.BaseURL, apiStopSendRTP, req, &res); err != nil {
		return err
	}
	// 经过测试，-500 是找不到流，-1 是已经停止
	if res.apiError.Code != codeTrue &&
		res.apiError.Code != -500 &&
		res.apiError.Code != -1 {
		res.apiError.Path = apiStopSendRTP
		return &res.apiError
	}
	return nil
}
