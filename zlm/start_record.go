package zlm

import (
	"context"
)

// StartRecordReq 是 StartRecord 的参数
type StartRecordReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 类型
	Type RecordFileType `query:"type"`
	// 保存目录
	CustomizedPath string `query:"customized_path"`
	// mp4 录像时间大小，单位秒，0 则采用配置项
	MaxSecond string `query:"max_second"`
}

// StartRecordRes 是 StartRecord 的返回值
type StartRecordRes struct {
	CodeMsg
	// 成功与否
	Result bool `json:"result"`
}

const (
	StartRecordPath = apiPathPrefix + "/startRecord"
)

// StartRecord 调用 /index/api/startRecord ，开始录制
func StartRecord(ctx context.Context, ser Server, req *StartRecordReq, res *StartRecordRes) error {
	return Request(ctx, ser, StartRecordPath, req, res)
}
