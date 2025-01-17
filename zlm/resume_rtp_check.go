package zlm

import (
	"context"
)

// ResumeRtpCheckReq 是 ResumeRtpCheck 参数
type ResumeRtpCheckReq struct {
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream_id"`
}

const (
	ResumeRtpCheckPath = "resumeRtpCheck"
)

// ResumeRtpCheck 调用 /index/api/resumeRtpCheck ，恢复RTP检测
func ResumeRtpCheck(ctx context.Context, ser Server, req *ResumeRtpCheckReq) error {
	return Request(ctx, ser, ResumeRtpCheckPath, req, new(CodeMsg))
}
