package zlm

import "context"

// DeleteRecordDirectoryReq 是 DeleteRecordDirectory 的参数
type DeleteRecordDirectoryReq struct {
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
	// 流的录像日期，格式为2020-01-01
	// 如果不是完整的日期，那么会删除失败
	Period string `query:"period"`
}

// deleteRecordDirectoryRes 是 DeleteRecordDirectory 返回值
type deleteRecordDirectoryRes struct {
	apiError
	Path string `json:"path"`
}

const (
	apiDeleteRecordDirectory = "deleteRecordDirectory"
)

// DeleteRecordDirectory 调用 /index/api/deleteRecordDirectory
func DeleteRecordDirectory(ctx context.Context, req *DeleteRecordDirectoryReq) error {
	// 请求
	var res deleteRecordDirectoryRes
	if err := request(ctx, req.BaseURL, apiDeleteRecordDirectory, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiDeleteRecordDirectory
		return &res.apiError
	}
	return nil
}
