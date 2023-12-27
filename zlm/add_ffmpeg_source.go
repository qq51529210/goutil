package zlm

import (
	"context"
)

// AddFFMPEGSourceReq 是 AddFFMPEGSource 的参数
type AddFFMPEGSourceReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// FFmpeg拉流地址,支持任意协议或格式(只要FFmpeg支持即可)
	SrcURL string `query:"src_url"`
	// FFmpeg rtmp推流地址，一般都是推给自己，例如rtmp://127.0.0.1/live/stream_form_ffmpeg
	DstURL string `query:"dst_url"`
	// FFmpeg推流成功超时时间
	TimeoutMS string `query:"timeout_ms"`
	// 是否开启hls录制，0/1
	EnableHLS string `query:"enable_hls"`
	// 是否开启mp4录制，0/1
	EnableMP4 string `query:"enable_mp4"`
	// 配置文件中FFmpeg命令参数模板key(非内容)，置空则采用默认模板:ffmpeg.cmd
	// ffmpeg_cmd_key
	CmdKey string `query:"ffmpeg_cmd_key"`
}

// addFFMPEGSourceRes 是 AddFFMPEGSource 返回值
type addFFMPEGSourceRes struct {
	apiError
	Data AddFFMPEGSourceResData `json:"data"`
}

// AddFFMPEGSourceResData 是 addFFMPEGSourceRes 的 Data 字段
type AddFFMPEGSourceResData struct {
	// 唯一标识
	Key string
}

const (
	apiAddFFmpegSource = "addFFmpegSource"
)

// AddFFMPEGSource 调用 /index/api/addFFmpegSource
// 通过fork FFmpeg进程的方式拉流代理，支持任意协议
// 返回 key
func AddFFMPEGSource(ctx context.Context, req *AddFFMPEGSourceReq) (string, error) {
	// 请求
	var res addFFMPEGSourceRes
	err := request(ctx, req.BaseURL, apiAddFFmpegSource, req, &res)
	if err != nil {
		return "", err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiAddFFmpegSource
		return "", &res.apiError
	}
	return res.Data.Key, nil
}
