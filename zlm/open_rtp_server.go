package zlm

import (
	"context"
)

// OpenRTPServerReq 是 OpenRTPServer 的参数
type OpenRTPServerReq struct {
	apiCall
	// 接收端口，0则为随机端口
	Port string `query:"port"`
	// tcp模式，0时为不启用tcp监听，1时为启用tcp监听，2时为tcp主动连接模式
	TCPMode string `query:"tcp_mode"`
	// 该端口绑定的流 id
	Stream string `query:"stream_id"`
	// 是否重用端口，默认为0，非必选参数，0/1
	ReusePort string `query:"re_use_port"`
	// 是否指定收流的rtp ssrc, 十进制数字，不指定或指定0时则不过滤rtp，非必选参数
	SSRC string `query:"ssrc"`
	// 是否为单音频track，用于语音对讲
	OnlyAudio string `query:"only_audio"`
}

// openRTPServerRes 是 OpenRTPServer 的返回值
type openRTPServerRes struct {
	apiError
	// 接收端口，0 随机端口号
	Port int `json:"port"`
}

const (
	apiOpenRTPServer = "openRtpServer"
)

// OpenRTPServer 调用 /index/api/openRtpServer
// 创建GB28181 RTP接收端口，如果该端口接收数据超时，则会自动被回收(不用调用closeRtpServer接口)
// 返回使用的端口
func OpenRTPServer(ctx context.Context, req *OpenRTPServerReq) (int, error) {
	// 请求
	var res openRTPServerRes
	err := request(ctx, &req.apiCall, apiOpenRTPServer, req, &res)
	if err != nil {
		return 0, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = req.apiCall.ID
		res.apiError.Path = apiOpenRTPServer
		return 0, &res.apiError
	}
	return res.Port, nil
}
