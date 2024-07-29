package zlm

import (
	"context"
)

// DelFFMPEGSourceReq 是 DelFFMPEGSource 的参数
type DelFFMPEGSourceReq struct {
	// addFFmpegSource 返回的 key
	Key string `query:"key"`
}

// delFFMPEGSourceRes 是 DelFFMPEGSource 返回值
type delFFMPEGSourceRes struct {
	CodeMsg
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

const (
	DelFFmpegSourcePath = apiPathPrefix + "/delFFmpegSource"
)

// DelFFmpegSource 调用 /index/api/delFFmpegSource ，返回是否成功，可以使用 close_streams 替代
func DelFFmpegSource(ctx context.Context, ser Server, req *DelFFMPEGSourceReq) (bool, error) {
	// 请求
	var res delFFMPEGSourceRes
	if err := Request(ctx, ser, DelFFmpegSourcePath, req, &res); err != nil {
		return false, err
	}
	// 经过测试，-500 应该是不存在的意思
	// 不存在也当它成功了
	if res.Code != CodeOK && res.Code != -500 {
		return false, &res.CodeMsg
	}
	return res.Data.Flag, nil
}
