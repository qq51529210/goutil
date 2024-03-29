package zlm

import (
	"context"
)

// StopRecordReq 是 StopRecord 的参数
type StopRecordReq struct {
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
	// 0为hls，1为mp4
	Type string `query:"type"`
}

// stopRecordRes 是 StopRecord 的返回值
type stopRecordRes struct {
	apiError
	// 成功与否
	Result bool `json:"result"`
}

const (
	apiStopRecord = "stopRecord"
)

// StopRecord 调用 /index/api/stopRecord
// 停止录制流
// 返回是否成功
func StopRecord(ctx context.Context, req *StopRecordReq) (bool, error) {
	var res stopRecordRes
	if err := request(ctx, req.BaseURL, apiStopRecord, req, &res); err != nil {
		return false, err
	}
	// 经过测试，-500 是找不到流，-1 是已经停止
	if res.apiError.Code != codeTrue &&
		res.apiError.Code != -500 &&
		res.apiError.Code != -1 {
		res.apiError.Path = apiStopRecord
		return false, &res.apiError
	}
	return res.Result, nil
}
