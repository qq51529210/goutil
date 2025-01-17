package zlm

import (
	"context"
)

// PauseRtpCheckReq 是 PauseRtpCheck 参数
type PauseRtpCheckReq struct {
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream_id"`
}

const (
	PauseRtpCheckPath = "pauseRtpCheck"
)

// PauseRtpCheck 调用 /index/api/pauseRtpCheck ，暂停RTP检测
func PauseRtpCheck(ctx context.Context, ser Server, req *PauseRtpCheckReq) error {
	return Request(ctx, ser, PauseRtpCheckPath, req, new(CodeMsg))
}
