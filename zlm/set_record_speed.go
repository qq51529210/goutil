package zlm

import "context"

// SetRecordSpeedReq 是 SetRecordSpeed 的参数
type SetRecordSpeedReq struct {
	// 协议，这个 postman 上没有，但是不添不行
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 要设置的录像倍速
	Speed string `query:"speed"`
}

// SetRecordSpeedRes 是 SetRecordSpeed 返回值
type SetRecordSpeedRes struct {
	CodeMsg
}

const (
	SetRecordSpeedPath = apiPathPrefix + "/setRecordSpeed"
)

// SetRecordSpeed 调用 /index/api/setRecordSpeed ，设置录像的播放速度
func SetRecordSpeed(ctx context.Context, ser Server, req *SetRecordSpeedReq, res *SetRecordSpeedRes) error {
	return Request(ctx, ser, SetRecordSpeedPath, req, res)
}
