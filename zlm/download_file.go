package zlm

import "context"

// DownloadFileReq 是 DownloadFile 的参数
type DownloadFileReq struct {
	// 文件的绝对路径
	FilePath string `query:"file_path"`
}

// DownloadFileRes 是 DownloadFile 返回值
type DownloadFileRes struct {
	CodeMsg
}

const (
	DownloadFilePath = apiPathPrefix + "/downloadFile"
)

// DownloadFile 调用 /index/api/downloadFile ，下载文件，会触发 on_http_access 回调
func DownloadFile(ctx context.Context, ser Server, req *DownloadFileReq, res *DownloadFileRes) error {
	return Request(ctx, ser, DownloadFilePath, req, res)
}
