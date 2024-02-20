package zlm

import (
	"context"
)

// AddStreamProxyReq 是 AddStreamProxy 参数
type AddStreamProxyReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
	// 拉流地址，例如rtmp://live.hkstv.hk.lxdns.com/live/hks2
	URL string `query:"url"`
	// rtsp拉流时，拉流方式，0：tcp，1：udp，2：组播
	RTPType string `query:"rtp_type"`
	// 拉流超时时间，单位秒，float类型
	TimeoutSec string `query:"timeout_sec"`
	// 拉流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string `query:"retry_count"`
	// 是否转换成hls协议，0/1
	EnableHLS string `query:"enable_hls"`
	// 是否mp4录制
	EnableMP4 string `query:"enable_mp4"`
	// 是否转换成rtsp协议，0/1
	EnableRTSP string `query:"enable_rtsp"`
	// 是否转换成rtmp/flv协议，0/1
	EnableRTMP string `query:"enable_rtmp"`
	// 是否转换成http-ts/ws-ts协议，0/1
	EnableTS string `query:"enable_ts"`
	// 是否转换成http-fmp4/ws-fmp4协议，0/1
	EnableFMP4 string `query:"enable_fmp4"`
	// 转协议时是否开启音频
	EnableAudio string `query:"enable_audio"`
	// 转协议时，无音频是否添加静音aac音频
	AddMuteAudio string `query:"add_mute_audio"`
	// mp4录制文件保存根目录，置空使用默认
	MP4SavePath string `query:"mp4_save_path"`
	// mp4录制切片大小，单位秒
	MP4MaxSecond string `query:"mp4_max_second"`
	// hls文件保存保存根目录，置空使用默认
	HLSSavePath string `query:"hls_save_path"`
}

// addStreamProxyRes 是 AddStreamProxy 返回值
type addStreamProxyRes struct {
	apiError
	Data AddStreamProxyResData `json:"data"`
}

// AddStreamProxyResData 是 addStreamProxyRes 的 Data 字段
type AddStreamProxyResData struct {
	// 流的唯一标识
	Key string
}

const (
	apiAddStreamProxy = "addStreamProxy"
)

// AddStreamProxy 调用 /index/api/addStreamProxy ，返回 key
func AddStreamProxy(ctx context.Context, req *AddStreamProxyReq) (string, error) {
	// 请求
	var res addStreamProxyRes
	if err := request(ctx, req.BaseURL, apiAddStreamProxy, req, &res); err != nil {
		return "", err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiAddStreamProxy
		return "", &res.apiError
	}
	return res.Data.Key, nil
}
