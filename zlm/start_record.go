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

// startRecordRes 是 StartRecord 的返回值
type startRecordRes struct {
	CodeMsg
	// 成功与否
	Result bool `json:"result"`
}

const (
	StartRecordPath = apiPathPrefix + "/startRecord"
)

// StartRecord 调用 /index/api/startRecord ，开始录制，返回是否成功
func StartRecord(ctx context.Context, ser Server, req *StartRecordReq) (bool, error) {
	// 请求
	var res startRecordRes
	if err := Request(ctx, ser, StartRecordPath, req, &res); err != nil {
		return false, err
	}
	if res.Code != CodeOK {
		return false, &res.CodeMsg
	}
	return res.Result, nil
}
