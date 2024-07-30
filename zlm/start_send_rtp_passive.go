package zlm

import (
	"context"
)

// StartSendRTPPassiveReq 是 StartSendRTPPassive 参数
type StartSendRTPPassiveReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// ssrc
	SSRC string `query:"ssrc"`
	// 使用的本机端口，默认为随机端口
	SrcPort string `query:"src_port"`
	// 默认为 96
	PT string `query:"pt"`
	// 负载类型，默认为 1
	UsePS RTPPayloadType `query:"use_ps"`
	// es 方式打包是否只打包音频
	OnlyAudio Boolean `query:"only_audio"`
	// 是否推送本地 MP4 录像
	FromMP4 Boolean `query:"from_mp4"`
	// 接收流的标识，发送同时接收，一般用于双向语言对讲
	RecvStreamID string `query:"recv_stream_id"`
}

// StartSendRTPPassiveRes 是 StartSendRTPPassive 返回值
type StartSendRTPPassiveRes struct {
	CodeMsg
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

const (
	StartSendRtpPassivePath = apiPathPrefix + "/startSendRtpPassive"
)

// StartSendRTPPassive 调用 /index/api/startSendRtpPassive ，被动推流
func StartSendRTPPassive(ctx context.Context, ser Server, req *StartSendRTPPassiveReq, res *StartSendRTPPassiveRes) error {
	return Request(ctx, ser, StartSendRtpPassivePath, req, res)
}
