package zlm

import (
	"context"
)

// StopRecordReq 是 StopRecord 的参数
type StopRecordReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 录像文件类型
	Type RecordFileType `query:"type"`
}

// StopRecordRes 是 StopRecord 的返回值
type StopRecordRes struct {
	CodeMsg
	// 成功与否
	Result bool `json:"result"`
}

const (
	StopRecordPath = apiPathPrefix + "/stopRecord"
)

// StopRecord 调用 /index/api/stopRecord ，停止录制
// 经过测试，-500 是找不到流，-1 是已经停止
func StopRecord(ctx context.Context, ser Server, req *StopRecordReq, res *StopRecordRes) error {
	return Request(ctx, ser, StopRecordPath, req, res)
}
