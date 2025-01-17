package zlm

import (
	"context"
)

// PauseRecordReq 是 PauseRecord 参数
type PauseRecordReq struct {
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 暂停/恢复
	Pause bool `query:"pause"`
}

const (
	PauseRecordPath = "pauseRecord"
)

// PauseRecord 调用 /index/api/pauseRecord ，暂停 mp4 文件推流
func PauseRecord(ctx context.Context, ser Server, req *PauseRecordReq) error {
	return Request(ctx, ser, PauseRecordPath, req, new(CodeMsg))
}
