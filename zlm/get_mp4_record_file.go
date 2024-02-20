package zlm

import "context"

// GetMp4RecordFileReq 是 GetMp4RecordFile 的参数
type GetMp4RecordFileReq struct {
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
	// 录像文件保存自定义根目录，为空则采用配置文件设置
	CustomizedPath string `query:"customized_path"`
	// 流的录像日期，格式为2020-02-01
	// 如果不是完整的日期，那么是搜索录像文件夹列表
	// 否则搜索对应日期下的mp4文件列表
	Period string `query:"period"`
}

// getMp4RecordFileRes 是 GetMp4RecordFile 返回值
type getMp4RecordFileRes struct {
	apiError
	Data struct {
		Path     []string `json:"path"`
		RootPath string   `json:"rootPath"`
	} `json:"data"`
}

const (
	apiGetMp4RecordFile = "getMp4RecordFile"
)

// GetMp4RecordFile 调用 /index/api/getMp4RecordFile
func GetMp4RecordFile(ctx context.Context, req *GetMp4RecordFileReq) error {
	// 请求
	var res getMp4RecordFileRes
	if err := request(ctx, req.BaseURL, apiGetMp4RecordFile, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiGetMp4RecordFile
		return &res.apiError
	}
	return nil
}
