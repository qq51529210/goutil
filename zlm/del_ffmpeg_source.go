package zlm

import (
	"context"
)

// DelFFMPEGSourceReq 是 DelFFMPEGSource 的参数
type DelFFMPEGSourceReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 流的唯一标识
	Key string `query:"key"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
}

// delFFMPEGSourceRes 是 DelFFMPEGSource 返回值
type delFFMPEGSourceRes struct {
	apiError
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

// DelFFMPEGSourceResData 是 addFFMPEGSourceRes 的 Data 字段
type DelFFMPEGSourceResData struct {
	// 唯一标识
	Key string
}

const (
	apiDelFFmpegSource = "delFFmpegSource"
)

// DelFFmpegSource 调用 /index/api/delFFmpegSource
func DelFFmpegSource(ctx context.Context, req *DelFFMPEGSourceReq) (bool, error) {
	// 请求
	var res delFFMPEGSourceRes
	if err := request(ctx, req.BaseURL, apiDelFFmpegSource, req, &res); err != nil {
		return false, err
	}
	if res.apiError.Code != codeTrue {
		// -500 是没有找到流，也算成功
		if res.apiError.Code != -500 {
			res.apiError.Path = apiDelFFmpegSource
			return false, &res.apiError
		}
	}
	return res.Data.Flag, nil
}
