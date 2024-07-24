package zlm

import "context"

// SeekRecordStampReq 是 SeekRecordStamp 的参数
type SeekRecordStampReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 这个 postman 上没有，但是不添不行
	Schema string `query:"schema"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流 id ，例如 test
	Stream string `query:"stream"`
	// 要设置的录像播放位置
	Stamp string `query:"stamp"`
}

// seekRecordStampRes 是 SeekRecordStamp 返回值
type seekRecordStampRes struct {
	apiError
}

const (
	apiSeekRecordStamp = "seekRecordStamp"
)

// SeekRecordStamp 调用 /index/api/seekRecordStamp
func SeekRecordStamp(ctx context.Context, req *SeekRecordStampReq) error {
	// 请求
	var res seekRecordStampRes
	if err := request(ctx, req.BaseURL, apiSeekRecordStamp, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiSeekRecordStamp
		return &res.apiError
	}
	return nil
}
