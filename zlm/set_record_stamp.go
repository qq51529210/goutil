package zlm

import "context"

// SetRecordStampReq 是 SetRecordStamp 的参数
type SetRecordStampReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流 id ，例如 test
	Stream string `query:"stream"`
	// 要设置的录像播放位置
	Stamp string `query:"stamp"`
}

// setRecordStampRes 是 SetRecordStamp 返回值
type setRecordStampRes struct {
	apiError
}

const (
	apiSetRecordStamp = "setRecordStamp"
)

// SetRecordStamp 调用 /index/api/setRecordStamp
func SetRecordStamp(ctx context.Context, req *SetRecordStampReq) error {
	// 请求
	var res setRecordStampRes
	if err := request(ctx, req.BaseURL, apiSetRecordStamp, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiSetRecordStamp
		return &res.apiError
	}
	return nil
}
