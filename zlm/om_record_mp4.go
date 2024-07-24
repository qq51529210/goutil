package zlm

import (
	"context"
)

// OnRecordMP4Req 表示 on_record_mp4 提交的数据
type OnRecordMP4Req struct {
	// 服务器id，通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 录制的流应用名
	App string `json:"app"`
	// 录制的流ID
	Stream string `json:"stream"`
	// 录像文件名
	FileName string `json:"file_name"`
	// 录像文件绝对路径
	FilePath string `json:"file_path"`
	// 录像文件的大小，单位字节
	FileSize int64 `json:"file_size"`
	// 文件所在目录路径
	Folder string `json:"folder"`
	// 开始录制时间戳
	StartTime int64 `json:"start_time"`
	// 录制时长，单位秒
	TimeLen float64 `json:"time_len"`
	// 点播相对 url 路径
	URL string `json:"url"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnRecordMP4 处理 zlm 的 on_record_mp4 回调
func OnRecordMP4(ctx context.Context, req *OnRecordMP4Req) {
}
