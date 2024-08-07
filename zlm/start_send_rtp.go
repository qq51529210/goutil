package zlm

import (
	"context"
)

// StartSendRTPReq 是 StartSendRTP 参数
type StartSendRTPReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// ssrc
	SSRC string `query:"ssrc"`
	// 目标 ip
	DstIP string `query:"dst_url"`
	// 目标端口
	DstPort string `query:"dst_port"`
	// 是否为udp模式，否则为tcp模式，0/1
	IsUDP string `query:"is_udp"`
	// 使用的本机端口，默认为随机端口
	SrcPort string `query:"src_port"`
	// 默认为 96
	PT string `query:"pt"`
	// 负载类型，默认为 1
	UsePS RTPPayloadType `query:"use_ps"`
	// es 方式是否只打包音频
	OnlyAudio Boolean `query:"only_audio"`
}

// StartSendRTPRes 是 StartSendRTP 返回值
type StartSendRTPRes struct {
	CodeMsg
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

const (
	StartSendRtpPath = apiPathPrefix + "/startSendRtp"
)

// StartSendRTP 调用 /index/api/startSendRtp ，开始推流
func StartSendRTP(ctx context.Context, ser Server, req *StartSendRTPReq, res *StartSendRTPRes) error {
	return Request(ctx, ser, StartSendRtpPath, req, res)
}
