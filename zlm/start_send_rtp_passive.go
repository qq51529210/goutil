package zlm

import (
	"context"
)

// StartSendRTPPassiveReq 是 StartSendRTPPassive 参数
type StartSendRTPPassiveReq struct {
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
	// 发送rtp同时接收，一般用于双向语言对讲, 如果不为空，说明开启接收，值为接收流的id
	RecvStreamID string `query:"recv_stream_id"`
}

// startSendRTPPassiveRes 是 StartSendRTPPassive 返回值
type startSendRTPPassiveRes struct {
	apiError
	// 使用的本地端口号
	LocalPort int `json:"local_port"`
}

const (
	apiStartSendRtpPassive = "startSendRtpPassive"
)

// StartSendRTPPassive 调用 /index/api/startSendRtpPassive
// 作为GB28181 Passive TCP服务器；该接口支持rtsp/rtmp等协议转ps-rtp被动推流。调用该接口，zlm会启动tcp服务器等待连接请求，连接建立后，zlm会关闭tcp服务器，然后源源不断的往客户端推流。
// 第一次推流失败会直接返回错误，成功一次后，后续失败也将无限重试(不停地建立tcp监听，超时后再关闭)。
// 返回使用的本地端口号
func StartSendRTPPassive(ctx context.Context, req *StartSendRTPPassiveReq) (int, error) {
	// 请求
	var res startSendRTPPassiveRes
	err := request(ctx, req.BaseURL, apiStartSendRtpPassive, req, &res)
	if err != nil {
		return 0, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiStartSendRtpPassive
		return 0, &res.apiError
	}
	return res.LocalPort, nil
}
