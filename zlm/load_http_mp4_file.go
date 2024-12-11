package zlm

import "context"

// LoadHTTPMP4FileReq 是 LoadMP4File 的参数
type LoadHTTPMP4FileReq struct {
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// mp4 文件的路径
	Path string `query:"file_path"`
}

// LoadHTTPMP4FileRes 是 LoadMP4File 返回值
type LoadHTTPMP4FileRes struct {
	CodeMsg
	Data struct {
		Flag bool `json:"flag"`
	} `json:"data"`
}

const (
	LoadHTTPMP4FilePath = apiPathPrefix + "/loadHttpMP4File"
)

// LoadHTTPMP4File 调用 /index/api/loadHttpMP4File ，加载本地的 mp4 文件，主要用于推流
func LoadHTTPMP4File(ctx context.Context, ser Server, req *LoadHTTPMP4FileReq, res *LoadHTTPMP4FileRes) error {
	return Request(ctx, ser, LoadHTTPMP4FilePath, req, res)
}
