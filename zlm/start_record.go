package zlm

import (
	"context"
)

// StartRecordReq 是 StartRecord 的参数
type StartRecordReq struct {
	apiCall
	// 筛选虚拟主机
	VHost string `query:"vhost"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 0为hls，1为mp4
	Type string `query:"type"`
	// 录像保存目录
	CustomizedPath string `query:"customized_path"`
	// mp4录像切片时间大小,单位秒，置0则采用配置项
	MaxSecond string `query:"max_second"`
}

// startRecordRes 是 StartRecord 的返回值
type startRecordRes struct {
	apiError
	// 成功与否
	Result bool `json:"result"`
}

const (
	apiStartRecord = "startRecord"
)

// StartRecord 调用 /index/api/startRecord
// 开始录制hls或MP4
// 返回是否成功
func StartRecord(ctx context.Context, req *StartRecordReq) (bool, error) {
	// 请求
	var res startRecordRes
	err := request(ctx, &req.apiCall, apiStartRecord, req, &res)
	if err != nil {
		return false, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = req.apiCall.ID
		res.apiError.Path = apiStartRecord
		return false, &res.apiError
	}
	//
	return res.Result, nil
}
