package zlm

import (
	"context"
)

// RTPStreamModel 收流模式
type RTPStreamModel string

const (
	// udp
	RTPStreamModelUDP RTPStreamModel = Zero
	// tcp 被动
	RTPStreamModelPassive RTPStreamModel = One
	// tcp 主动
	RTPStreamModelActive RTPStreamModel = Two
)

// OpenRTPServerReq 是 OpenRTPServer 的参数
type OpenRTPServerReq struct {
	// 接收端口，0 则为随机端口
	Port string `query:"port"`
	// tcp模式，0时为不启用tcp监听，1时为启用tcp监听，2时为tcp主动连接模式
	StreamMode RTPStreamModel `query:"tcp_mode"`
	// 绑定的流标识
	Stream string `query:"stream_id"`
	// 是否重用端口，默认为 0
	ReusePort Boolean `query:"re_use_port"`
	// ssrc
	SSRC string `query:"ssrc"`
	// 是否只有音频
	OnlyAudio Boolean `query:"only_audio"`
	// 自定义上下文数据
	UserData string `query:"userdata"`
}

// OpenRTPServerRes 是 OpenRTPServer 的返回值
type OpenRTPServerRes struct {
	CodeMsg
	// 接收端口，0 随机端口号
	Port int `json:"port"`
}

const (
	OpenRTPServerPath = apiPathPrefix + "/openRtpServer"
)

// OpenRTPServer 调用 /index/api/openRtpServer ，打开端口准备收流
// 经过测试，如果返回端口为 0 但是没有错误，说明有流
func OpenRTPServer(ctx context.Context, ser Server, req *OpenRTPServerReq, res *OpenRTPServerRes) error {
	return Request(ctx, ser, OpenRTPServerPath, req, res)
}
