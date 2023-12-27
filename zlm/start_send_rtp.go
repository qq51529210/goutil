package zlm

import (
	"context"
)

// StartSendRTPReq 是 StartSendRTP 参数
type StartSendRTPReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
	// 推流的rtp的ssrc,指定不同的ssrc可以同时推流到多个服务器
	SSRC string `query:"ssrc"`
	// 目标ip或域名
	DstURL string `query:"dst_url"`
	// 目标端口
	DstPort string `query:"dst_port"`
	// 是否为udp模式，否则为tcp模式，0/1
	IsUDP string `query:"is_udp"`
	// 使用的本机端口，为0或不传时默认为随机端口
	SrcPort string `query:"src_port"`
	// 发送时，rtp的pt（uint8_t）,不传时默认为96
	PT string `query:"pt"`
	// 发送时，rtp的负载类型。为1时，负载为ps；为0时，为es；不传时默认为1
	UsePS string `query:"use_ps"`
	// rtp es方式打包时，是否只打包音频，该参数非必选参数
	OnlyAudio string `query:"only_audio"`
	// 是否推送本地MP4录像，该参数非必选参数，0/1
	FromMP4 string `query:"from_mp4"`
	// udp方式推流时，是否开启rtcp发送和rtcp接收超时判断，开启后(默认关闭)，如果接收rtcp超时，将导致主动停止rtp发送，0/1
	UDPRtcpTimeout string `query:"udp_rtcp_timeout"`
	// 发送rtp同时接收，一般用于双向语言对讲, 如果不为空，说明开启接收，值为接收流的id
	RecvStreamID string `query:"recv_stream_id"`
}

// startSendRTPRes 是 StartSendRTP 返回值
type startSendRTPRes struct {
	apiError
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

const (
	apiStartSendRtp = "startSendRtp"
)

// StartSendRTP 调用 /index/api/startSendRtp
// 作为GB28181客户端，启动ps-rtp推流，支持rtp/udp方式；该接口支持rtsp/rtmp等协议转ps-rtp推流。
// 第一次推流失败会直接返回错误，成功一次后，后续失败也将无限重试。
// 返回使用的本地端口号
func StartSendRTP(ctx context.Context, req *StartSendRTPReq) (int, error) {
	// 请求
	var res startSendRTPRes
	err := request(ctx, req.BaseURL, apiStartSendRtp, req, &res)
	if err != nil {
		return 0, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiStartSendRtp
		return 0, &res.apiError
	}
	return res.LocalPort, nil
}
