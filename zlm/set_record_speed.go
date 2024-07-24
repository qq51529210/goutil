package zlm

import "context"

// SetRecordSpeedReq 是 SetRecordSpeed 的参数
type SetRecordSpeedReq struct {
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
	// 要设置的录像倍速
	Speed string `query:"speed"`
}

// setRecordSpeedRes 是 SetRecordSpeed 返回值
type setRecordSpeedRes struct {
	apiError
}

const (
	apiSetRecordSpeed = "setRecordSpeed"
)

// SetRecordSpeed 调用 /index/api/setRecordSpeed
func SetRecordSpeed(ctx context.Context, req *SetRecordSpeedReq) error {
	// 请求
	var res setRecordSpeedRes
	if err := request(ctx, req.BaseURL, apiSetRecordSpeed, req, &res); err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiSetRecordSpeed
		return &res.apiError
	}
	return nil
}
