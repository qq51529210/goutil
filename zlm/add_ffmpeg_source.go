package zlm

import (
	"context"
)

// AddFFMPEGSourceReq 是 AddFFMPEGSource 的参数
type AddFFMPEGSourceReq struct {
	// 拉流源地址
	SrcURL string `query:"src_url"`
	// 转推流地址
	DstURL string `query:"dst_url"`
	// 超时时间，单位豪秒
	Timeout string `query:"timeout_ms"`
	// 是否开启 hls 录制
	EnableHLS Boolean `query:"enable_hls"`
	// 是否开启 mp4 录制
	EnableMP4 Boolean `query:"enable_mp4"`
	// 置空则采用默认模板: ffmpeg.cmd
	CmdKey string `query:"ffmpeg_cmd_key"`
}

// addFFMPEGSourceRes 是 AddFFMPEGSource 返回值
type addFFMPEGSourceRes struct {
	CodeMsg
	Data AddFFMPEGSourceResData `json:"data"`
}

// AddFFMPEGSourceResData 是 addFFMPEGSourceRes 的 Data 字段
type AddFFMPEGSourceResData struct {
	// 唯一标识
	Key string
}

const (
	AddFFMPEGSourcePath = apiPathPrefix + "/addFFmpegSource"
)

// AddFFMPEGSource 调用 /index/api/addFFmpegSource ，从指定源地址拉流，返回 key
func AddFFMPEGSource(ctx context.Context, ser Server, req *AddFFMPEGSourceReq) (string, error) {
	// 请求
	var res addFFMPEGSourceRes
	if err := Request(ctx, ser, AddFFMPEGSourcePath, req, &res); err != nil {
		return "", err
	}
	if res.Code != CodeOK {
		return "", &res.CodeMsg
	}
	return res.Data.Key, nil
}
