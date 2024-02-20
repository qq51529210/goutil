package zlm

import "context"

// DownloadFileReq 是 DownloadFile 的参数
type DownloadFileReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// mp4 文件的绝对路径
	FilePath string `query:"file_path"`
}

// downloadFileRes 是 DownloadFile 返回值
type downloadFileRes struct {
	apiError
}

const (
	apiDownloadFile = "downloadFile"
)

// DownloadFile 调用 /index/api/downloadFile
func DownloadFile(ctx context.Context, req *DownloadFileReq) error {
	// 请求
	var res downloadFileRes
	if err := request(ctx, req.BaseURL, apiDownloadFile, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiDownloadFile
		return &res.apiError
	}
	return nil
}
