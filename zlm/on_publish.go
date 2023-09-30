package zlm

import (
	"github.com/gin-gonic/gin"
	"github.com/qq51529210/log"
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
	// // TCP链接唯一ID
	// ID string `json:"id"`
	// // 推流器ip
	// IP string `json:"ip"`
	// // 推流器端口号
	// Port int `json:"port"`
}

// OnPublishRes 表示 OnPublish 返回的数据
type OnPublishRes struct {
	// 错误代码，0代表允许推流
	Code int `json:"code"`
	// 是否转换成hls协议
	EnableHLS *bool `json:"enable_hls"`
	// 是否允许mp4录制
	EnableMP4 *bool `json:"enable_mp4"`
	// 是否转rtsp协议
	EnableRTSP *bool `json:"enable_rtsp"`
	// 是否转rtmp/flv协议
	EnableRTMP *bool `json:"enable_rtmp"`
	// 是否转http-ts/ws-ts协议
	EnableTS *bool `json:"enable_ts"`
	// 是否转http-fmp4/ws-fmp4协议
	EnableFMP4 *bool `json:"enable_fmp4"`
	// 转协议时是否开启音频
	EnableAudio *bool `json:"enable_audio"`
	// 转协议时，无音频是否添加静音aac音频
	EnableAddMuteAudio *bool `json:"add_mute_audio"`
	// mp4录制文件保存根目录，置空使用默认
	MP4SavePath *string `json:"mp4_save_path"`
	// hls文件保存保存根目录，置空使用默认
	HLSSavePath *string `json:"hls_save_path"`
	// 断连续推延时，单位毫秒，置空使用配置文件默认值
	ContinuePushMS *bool `json:"continue_push_ms"`
	// MP4录制是否当作观看者参与播放人数计数
	MP4AsPlayer *bool `json:"mp4_as_player"`
	// 该流是否开启时间戳覆盖
	ModifyStamp *bool `json:"modify_stamp"`
}

// OnPublish 处理 zlm 的 on_publish 回调
func OnPublish(ctx *gin.Context, req *OnPublishReq, res *OnPublishRes) {
	// 获取实例
	ser := GetServer(req.MediaServerID)
	if !ser.IsOK() {
		return
	}
	req.TraceID, _ = ctx.Value(CtxKeyTraceID).(string)
	// 处理
	_true := true
	_false := false
	res.Code = -1
	res.EnableRTMP = &_true
	res.EnableRTSP = &_true
	res.MP4AsPlayer = &_true
	res.EnableMP4 = &_false
	// 回调
	HandlePublish(ctx, ser, req, res)
	// 日志
	log.Debugf("%s res code %d", req.TraceID, res.Code)
}
