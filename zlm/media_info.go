package zlm

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
