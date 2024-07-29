package zlm

import "context"

// LoadMP4FileReq 是 LoadMP4File 的参数
type LoadMP4FileReq struct {
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// mp4 文件的路径
	Path string `query:"file_path"`
}

const (
	LoadMP4FilePath = apiPathPrefix + "/loadMP4File"
)

// LoadMP4File 调用 /index/api/loadMP4File ，加载本地的 mp4 文件，主要用于推流
func LoadMP4File(ctx context.Context, ser Server, req *LoadMP4FileReq) error {
	// 请求
	var res CodeMsg
	if err := Request(ctx, ser, LoadMP4FilePath, req, &res); err != nil {
		return err
	}
	if res.Code != CodeOK {
		return &res
	}
	return nil
}
