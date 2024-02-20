package zlm

import "context"

// LoadMP4FileReq 是 LoadMP4File 的参数
type LoadMP4FileReq struct {
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
	// mp4 文件的绝对路径
	FilePath string `query:"file_path"`
}

// loadMP4FileRes 是 LoadMP4File 返回值
type loadMP4FileRes struct {
	apiError
}

const (
	apiLoadMP4File = "loadMP4File"
)

// LoadMP4File 调用 /index/api/loadMP4File
func LoadMP4File(ctx context.Context, req *LoadMP4FileReq) error {
	// 请求
	var res loadMP4FileRes
	if err := request(ctx, req.BaseURL, apiLoadMP4File, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiLoadMP4File
		return &res.apiError
	}
	return nil
}
