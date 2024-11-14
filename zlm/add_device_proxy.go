package zlm

import (
	"context"
)

// AddDeviceProxyReq 是 AddDeviceProxy 参数
type AddDeviceProxyReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
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
	// 设备IP
	IP string `query:"ip"`
	// 设备端口
	Port int `query:"port"`
	// 设备用户名
	Username string `query:"username"`
	// 设备密码
	Password string `query:"password"`
	// 通道索引，1-x
	Channel int `query:"channel"`
	// 通道索引，1-海康，2-大华
	Manufacturer int `query:"manufacturer"`
	// 码流，0-主，1-辅
	Subtype int `query:"subtype"`
}

// addDeviceProxyRes 是 AddDeviceProxy 返回值
type AddDeviceProxyRes struct {
	CodeMsg
	Data struct {
		// 流的唯一标识
		Key string
	} `json:"data"`
}

const (
	AddDeviceProxyPath = apiPathPrefix + "/addDeviceProxy"
)

// AddDeviceProxy 调用 /index/api/addDeviceProxy ，返回 key
func AddDeviceProxy(ctx context.Context, ser Server, req *AddDeviceProxyReq, res *AddDeviceProxyRes) error {
	return Request(ctx, ser, AddDeviceProxyPath, req, res)
}
