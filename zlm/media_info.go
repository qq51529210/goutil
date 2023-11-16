package zlm

import (
	"fmt"
)

// MediaInfoKey 媒体流标识
type MediaInfoKey struct {
	// 流的唯一标识
	App string `json:"app"`
	// 流的唯一标识
	Stream string `json:"stream"`
}

// MediaInfo 表示某一个媒体流的信息
type MediaInfo struct {
	MediaInfoKey
	// 音频轨道信息
	Audio []*MediaInfoAudioTrack `json:"audioTracks,omitempty" copy:"-"`
	// 视频轨道信息
	Video []*MediaInfoVideoTrack `json:"videoTracks,omitempty" copy:"-"`
	// 原始的综合轨道信息
	Tracks []map[string]any `json:"-"`
	// 服务
	Ser *Server `json:"-"`
	// 是否正在录制 hls
	IsRecordingHLS bool `json:"-"`
	// 是否正在录制 mp4
	IsRecordingMP4 bool `json:"-"`
	// 创建的时间
	Timestamp int64 `json:"-"`
	// 观看总人数
	TotalReaderCount int64  `json:"-"`
	OriginTypeStr    string `json:"-"`
	OriginURL        string `json:"-"`
}

// HasH265 返回是否包含 h265 视频流
func (m *MediaInfo) HasH265() bool {
	for _, v := range m.Video {
		if v.CodecName == "H265" {
			return true
		}
	}
	return false
}

// MediaInfoVideoTrack 是 MediaInfo 的 Video 字段
type MediaInfoVideoTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 视频fps
	FPS int `json:"fps"`
	// 视频高
	Height int `json:"height"`
	// 视频宽
	Width int `json:"width"`
}

// MediaInfoAudioTrack 是 MediaInfo 的 Audio 字段
type MediaInfoAudioTrack struct {
	// 编码类型名称
	CodecName string `json:"codecName"`
	// 音频采样位数
	SampleBit int `json:"sampleBit"`
	// 音频采样率
	SampleRate int `json:"sampleRate"`
}

// PlayInfo 用于返回播放地址
type PlayInfo struct {
	// 媒体流标识
	App string `json:"app,omitempty"`
	// 媒体流标识
	Stream string `json:"stream,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsp://ip:port/app/stream
	RTSP string `json:"rtsp,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsps://ip:port/app/stream
	RTSPs string `json:"rtsps,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmp://ip:port/app/stream
	RTMP string `json:"rtmp,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmps://ip:port/app/stream
	RTMPs string `json:"rtmps,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.flv
	Flv string `json:"flv,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.flv
	Flvs string `json:"flvs,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream/hls.m3u8
	HLS string `json:"hls,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream/hls.m3u8
	HLSs string `json:"hlss,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.ts
	TS string `json:"ts,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.ts
	TSs string `json:"tss,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.mp4
	FMP4 string `json:"fmp4,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.mp4
	FMP4s string `json:"fmp4s,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.flv
	WsFlv string `json:"wsFlv,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.flv
	WssFlv string `json:"wssFlv,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream/hls.m3u8
	WsHLS string `json:"wsHLS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream/hls.m3u8
	WssHLS string `json:"wssHLS,omitempty"`
	// ws://ip:port/app/stream.live.ts
	WsTS string `json:"wsTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.ts
	WssTS string `json:"wssTS,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.mp4
	WsFMP4 string `json:"wsFMP4,omitempty"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.mp4
	WssFMP4 string `json:"wssFMP4,omitempty"`
}

// init 主要初始化 url
func (m *PlayInfo) init(ser *Server, public bool, token string) {
	// url
	var ip string
	if public {
		ip = ser.PublicIP
	} else {
		ip = ser.PrivateIP
	}
	// RTMP
	if ser.RTMPPort != "" {
		if token != "" {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s?token=%s", ip, ser.RTMPPort, m.App, m.Stream, token)
		} else {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, ser.RTMPPort, m.App, m.Stream)
		}
	}
	if ser.RTMPSSLPort != "" {
		if token != "" {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s?token=%s", ip, ser.RTMPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, ser.RTMPSSLPort, m.App, m.Stream)
		}
	}
	if ser.HTTPPort != "" {
		if token != "" {
			m.Flv = fmt.Sprintf("http://%s:%s/%s/%s.live.flv?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
			m.WsFlv = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.Flv = fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, ser.HTTPPort, m.App, m.Stream)
			m.WsFlv = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, ser.HTTPPort, m.App, m.Stream)
		}
	}
	if ser.HTTPSSLPort != "" {
		if token != "" {
			m.Flvs = fmt.Sprintf("https://%s:%s/%s/%s.live.flv?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssFlv = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.Flvs = fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, ser.HTTPSSLPort, m.App, m.Stream)
			m.WssFlv = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// RTSP
	if ser.RTSPPort != "" {
		if token != "" {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s?token=%s", ip, ser.RTSPPort, m.App, m.Stream, token)
		} else {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, ser.RTSPPort, m.App, m.Stream)
		}
	}
	if ser.RTSPSSLPort != "" {
		if token != "" {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s?token=%s", ip, ser.RTSPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, ser.RTSPSSLPort, m.App, m.Stream)
		}
	}
	// HLS
	if ser.HTTPPort != "" {
		if token != "" {
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
			m.WsHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, ser.HTTPPort, m.App, m.Stream)
			m.WsHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8", ip, ser.HTTPPort, m.App, m.Stream)
		}
	}
	if ser.HTTPSSLPort != "" {
		if token != "" {
			m.HLSs = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.HLSs = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, ser.HTTPSSLPort, m.App, m.Stream)
			m.WssHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8", ip, ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// TS
	if ser.HTTPPort != "" {
		if token != "" {
			m.TS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
			m.WsTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.TS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, ser.HTTPPort, m.App, m.Stream)
			m.WsTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts", ip, ser.HTTPPort, m.App, m.Stream)
		}
	}
	if ser.HTTPSSLPort != "" {
		if token != "" {
			m.TSs = fmt.Sprintf("https://%s:%s/%s/%s.live.ts?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.TSs = fmt.Sprintf("https://%s:%s/%s/%s.live.ts", ip, ser.HTTPSSLPort, m.App, m.Stream)
			m.WssTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts", ip, ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// MP4
	if ser.HTTPPort != "" {
		if token != "" {
			m.FMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
			m.WsFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4?token=%s", ip, ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.FMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, ser.HTTPPort, m.App, m.Stream)
			m.WsFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4", ip, ser.HTTPPort, m.App, m.Stream)
		}
	}
	if ser.HTTPSSLPort != "" {
		if token != "" {
			m.FMP4s = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4?token=%s", ip, ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.FMP4s = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4", ip, ser.HTTPSSLPort, m.App, m.Stream)
			m.WssFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4", ip, ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
}
