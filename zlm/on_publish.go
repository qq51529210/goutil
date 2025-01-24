package zlm

import (
	"context"
)

// OnPublishReq 表示 on_publish 提交的数据
type OnPublishReq struct {
	// 虚拟主机
	VHost string `json:"vhost"`
	// 服务标识
	MediaServerID string `json:"mediaServerId"`
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// url 查询字符串
	Params string `json:"params"`
	// 自定义上下文数据
	UserData string `query:"userdata"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnPublishRes 表示 OnPublish 返回的数据
type OnPublishRes struct {
	// 错误代码
	Code int `json:"code"`
	// 是否转换成 hls 协议
	EnableHLS bool `json:"enable_hls,omitempty"`
	// 是否允许 mp4 录制
	EnableMP4 bool `json:"enable_mp4,omitempty"`
	// 是否转 rtsp 协议
	EnableRTSP bool `json:"enable_rtsp,omitempty"`
	// 是否转 rtmp/flv 协议
	EnableRTMP bool `json:"enable_rtmp,omitempty"`
	// 是否转 http-ts/ws-ts 协议
	EnableTS bool `json:"enable_ts,omitempty"`
	// 是否转 http-fmp4/ws-fmp4 协议
	EnableFMP4 bool `json:"enable_fmp4,omitempty"`
	// 转协议时是否开启音频
	EnableAudio bool `json:"enable_audio,omitempty"`
	// 转协议时，无音频是否添加静音aac音频
	EnableAddMuteAudio bool `json:"add_mute_audio,omitempty"`
	// mp4 录制文件保存根目录，置空使用默认
	MP4SavePath string `json:"mp4_save_path,omitempty"`
	// hls 文件保存保存根目录，置空使用默认
	HLSSavePath string `json:"hls_save_path,omitempty"`
	// 推流断开后可以在超时时间内重新连接上继续推流，这样播放器会接着播放，单位毫秒，置空使用配置文件默认值
	ContinuePushMS bool `json:"continue_push_ms,omitempty"`
	// mp4 录制是否当作观看者参与播放人数计数
	MP4AsPlayer bool `json:"mp4_as_player,omitempty"`
	// 是否开启时间戳覆盖
	ModifyStamp bool `json:"modify_stamp,omitempty"`
}

// OnPublish 处理 zlm 的 on_publish 回调
func OnPublish(ctx context.Context, req *OnPublishReq, res *OnPublishRes) {
}
