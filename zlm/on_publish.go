package zlm

import (
	"context"
)

// OnPublishReq 表示 on_publish 提交的数据
type OnPublishReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 推流的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 推流url参数
	Params string `json:"params"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnPublishRes 表示 OnPublish 返回的数据
type OnPublishRes struct {
	// 错误代码，0代表允许推流
	Code int `json:"code"`
	// 是否转换成hls协议
	EnableHLS bool `json:"enable_hls,omitempty"`
	// 是否允许mp4录制
	EnableMP4 bool `json:"enable_mp4,omitempty"`
	// 是否转rtsp协议
	EnableRTSP bool `json:"enable_rtsp,omitempty"`
	// 是否转rtmp/flv协议
	EnableRTMP bool `json:"enable_rtmp,omitempty"`
	// 是否转http-ts/ws-ts协议
	EnableTS bool `json:"enable_ts,omitempty"`
	// 是否转http-fmp4/ws-fmp4协议
	EnableFMP4 bool `json:"enable_fmp4,omitempty"`
	// 转协议时是否开启音频
	EnableAudio bool `json:"enable_audio,omitempty"`
	// 转协议时，无音频是否添加静音aac音频
	EnableAddMuteAudio bool `json:"add_mute_audio,omitempty"`
	// mp4录制文件保存根目录，置空使用默认
	MP4SavePath string `json:"mp4_save_path,omitempty"`
	// hls文件保存保存根目录，置空使用默认
	HLSSavePath string `json:"hls_save_path,omitempty"`
	// 断连续推延时，单位毫秒，置空使用配置文件默认值
	ContinuePushMS bool `json:"continue_push_ms,omitempty"`
	// MP4录制是否当作观看者参与播放人数计数
	MP4AsPlayer bool `json:"mp4_as_player,omitempty"`
	// 该流是否开启时间戳覆盖
	ModifyStamp bool `json:"modify_stamp,omitempty"`
}

// OnPublish 处理 zlm 的 on_publish 回调
func OnPublish(ctx context.Context, req *OnPublishReq, res *OnPublishRes) {
}
