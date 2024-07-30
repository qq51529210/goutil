package zlm

import (
	"context"
)

// DelFFMPEGSourceReq 是 DelFFMPEGSource 的参数
type DelFFMPEGSourceReq struct {
	// addFFmpegSource 返回的 key
	Key string `query:"key"`
}

// DelFFMPEGSourceRes 是 DelFFMPEGSource 返回值
type DelFFMPEGSourceRes struct {
	CodeMsg
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

const (
	DelFFmpegSourcePath = apiPathPrefix + "/delFFmpegSource"
)

// DelFFmpegSource 调用 /index/api/delFFmpegSource ，返回是否成功，可以使用 close_streams 替代
// 经过测试 code=-500 应该是流不存在的意思
func DelFFmpegSource(ctx context.Context, ser Server, req *DelFFMPEGSourceReq, res *DelFFMPEGSourceRes) error {
	return Request(ctx, ser, DelFFmpegSourcePath, req, res)
}
