package zlm

import (
	"context"
)

// StopSendRTPReq 是 StopSendRTP 参数
type StopSendRTPReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 停止指定 ssrc 的推流，默认关闭所有推流
	SSRC string `query:"ssrc"`
}

// StopSendRTPRes 是 StopSendRTP 的返回值
type StopSendRTPRes struct {
	CodeMsg
}

const (
	StopSendRTPPath = apiPathPrefix + "/stopSendRtp"
)

// StopSendRTP 调用 /index/api/stopSendRtp ，停止推流
// 经过测试，-500 是找不到流，-1 是已经停止
func StopSendRTP(ctx context.Context, ser Server, req *StopSendRTPReq, res *StopSendRTPRes) error {
	return Request(ctx, ser, StopSendRTPPath, req, res)
}
