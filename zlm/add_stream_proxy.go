package zlm

import (
	"context"
)

// AddStreamProxyReq 是 AddStreamProxy 参数
type AddStreamProxyReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 拉流的源地址
	URL string `query:"url"`
	// rtsp 拉流方式
	RTPType RTSPRTPType `query:"rtp_type"`
	// 拉流超时时间，单位秒，float类型
	Timeout string `query:"timeout_sec"`
	// 拉流重试次数,不传此参数或传值<=0时，则无限重试
	RetryCount string `query:"retry_count"`
	// 是否启动 mp4 录制
	EnableMP4Record Boolean `query:"enable_mp4"`
	// 是否转换成 hls 协议
	EnableHLS Boolean `query:"enable_hls"`
	// 是否转换成 rtsp 协议
	EnableRTSP Boolean `query:"enable_rtsp"`
	// 是否转换成 rtmp/flv 协议
	EnableRTMP Boolean `query:"enable_rtmp"`
	// 是否转换成 http-ts/ws-ts 协议
	EnableTS Boolean `query:"enable_ts"`
	// 是否转换成 http-fmp4/ws-fmp4 协议
	EnableFMP4 Boolean `query:"enable_fmp4"`
	// 转协议时是否开启音频
	EnableAudio Boolean `query:"enable_audio"`
	// 转协议时，无音频是否添加静音 aac 音频
	AddMuteAudio Boolean `query:"add_mute_audio"`
	// mp4 录制文件保存根目录，置空使用默认
	MP4SavePath string `query:"mp4_save_path"`
	// mp4 录制切片大小，单位秒
	MP4MaxSecond string `query:"mp4_max_second"`
	// hls 文件保存保存根目录，置空使用默认
	HLSSavePath string `query:"hls_save_path"`
}

// addStreamProxyRes 是 AddStreamProxy 返回值
type addStreamProxyRes struct {
	CodeMsg
	Data AddStreamProxyResData `json:"data"`
}

// AddStreamProxyResData 是 addStreamProxyRes 的 Data 字段
type AddStreamProxyResData struct {
	// 流的唯一标识
	Key string
}

const (
	AddStreamProxyPath = apiPathPrefix + "/addStreamProxy"
)

// AddStreamProxy 调用 /index/api/addStreamProxy ，rtsp/rtmp/hls/http-ts/http-flv 拉流，返回 key
func AddStreamProxy(ctx context.Context, ser Server, req *AddStreamProxyReq) (string, error) {
	// 请求
	var res addStreamProxyRes
	if err := Request(ctx, ser, AddStreamProxyPath, req, &res); err != nil {
		return "", err
	}
	if res.Code != CodeOK {
		return "", &res.CodeMsg
	}
	return res.Data.Key, nil
}
