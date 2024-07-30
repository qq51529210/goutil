package zlm

import "context"

// GetMP4RecordFileReq 是 GetMP4RecordFile 的参数
type GetMP4RecordFileReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 录像文件保存自定义根目录，为空则采用配置文件设置
	CustomizedPath string `query:"customized_path"`
	// 流的录像日期，格式为 2020-02-01
	// 如果不是完整的日期，那么是搜索录像文件夹列表
	// 否则搜索对应日期下的 mp4 文件列表
	Period string `query:"period"`
}

// GetMp4RecordFileRes 是 GetMP4RecordFile 返回值
type GetMp4RecordFileRes struct {
	CodeMsg
	Data struct {
		Path     []string `json:"path"`
		RootPath string   `json:"rootPath"`
	} `json:"data"`
}

const (
	GetMP4RecordFilePath = apiPathPrefix + "/getMp4RecordFile"
)

// GetMP4RecordFile 调用 /index/api/getMp4RecordFile 查询指定日期的录像文件
func GetMP4RecordFile(ctx context.Context, ser Server, req *GetMP4RecordFileReq, res *GetMp4RecordFileRes) error {
	return Request(ctx, ser, GetMP4RecordFilePath, req, res)
}
