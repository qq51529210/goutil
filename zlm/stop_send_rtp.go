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

const (
	StopSendRTPPath = apiPathPrefix + "/stopSendRtp"
)

// StopSendRTP 调用 /index/api/stopSendRtp
// 停止 GB28181 rtp 推流。
func StopSendRTP(ctx context.Context, ser Server, req *StopSendRTPReq) error {
	var res CodeMsg
	if err := Request(ctx, ser, StopSendRTPPath, req, &res); err != nil {
		return err
	}
	// 经过测试，-500 是找不到流，-1 是已经停止
	if res.Code != CodeOK && res.Code != -500 && res.Code != -1 {
		return &res
	}
	return nil
}
