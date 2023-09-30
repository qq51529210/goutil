package zlm

import (
	"fmt"
)

type mediaInfoKey struct {
	App    string
	Stream string
}

// MediaInfo 表示某一个媒体流的信息
type MediaInfo struct {
	// 流的标识
	App string `json:"app"`
	// 流的标识
	Stream string `json:"stream"`
	// 音频轨道信息
	Audio []*MediaInfoAudioTrack `json:"audioTracks"`
	// 视频轨道信息
	Video []*MediaInfoVideoTrack `json:"videoTracks"`
	// 原始的综合轨道信息
	Tracks []map[string]any `json:"-"`
	// 服务
	Ser *Server `json:"-"`
	// 是否正在录制 hls
	IsRecordingHLS bool `json:"-"`
	// 是否正在录制 mp4
	IsRecordingMP4 bool `json:"-"`
	// 创建的时间
	Timestamp     int64  `json:"-"`
	OriginTypeStr string `json:"-"`
	OriginURL     string `json:"-"`
}

// PlayInfo 返回播放信息，public 决定使用内网还是外网 ip
func (m *MediaInfo) PlayInfo(public bool) *PlayInfo {
	info := new(PlayInfo)
	info.init(m, public, _playToken.new())
	return info
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

func (m *MediaInfo) init(ser *Server, data *MediaListData) {
	m.App = data.App
	m.Stream = data.Stream
	m.Tracks = data.Tracks
	m.initTracks()
	m.Ser = ser
	m.IsRecordingHLS = data.IsRecordingHLS
	m.IsRecordingMP4 = data.IsRecordingMP4
	m.Timestamp = data.CreateStamp
}

// initTracks 区分出音/视频轨道
func (m *MediaInfo) initTracks() {
	// 这里假定，返回的数据格式没有错误
	for _, track := range m.Tracks {
		switch int64(track["codec_type"].(float64)) {
		case 0:
			t := new(MediaInfoVideoTrack)
			t.CodecName = track["codec_id_name"].(string)
			t.FPS = int(track["fps"].(float64))
			t.Height = int(track["height"].(float64))
			t.Width = int(track["width"].(float64))
			m.Video = append(m.Video, t)
		case 1:
			t := new(MediaInfoAudioTrack)
			t.CodecName = track["codec_id_name"].(string)
			t.SampleBit = int(track["sample_bit"].(float64))
			t.SampleRate = int(track["sample_rate"].(float64))
			m.Audio = append(m.Audio, t)
		}
	}
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
	*MediaInfo
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsp://ip:port/app/stream
	RTSP string `json:"rtsp"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtsps://ip:port/app/stream
	RTSPs string `json:"rtsps"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmp://ip:port/app/stream
	RTMP string `json:"rtmp"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// rtmps://ip:port/app/stream
	RTMPs string `json:"rtmps"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.flv
	Flv string `json:"flv"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.flv
	Flvs string `json:"flvs"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream/hls.m3u8
	HLS string `json:"hls"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream/hls.m3u8
	HLSs string `json:"hlss"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.ts
	TS string `json:"ts"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.ts
	TSs string `json:"tss"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// http://ip:port/app/stream.live.mp4
	FMP4 string `json:"fmp4"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// https://ip:port/app/stream.live.mp4
	FMP4s string `json:"fmp4s"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.flv
	WsFlv string `json:"wsFlv"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.flv
	WssFlv string `json:"wssFlv"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream/hls.m3u8
	WsHLS string `json:"wsHLS"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream/hls.m3u8
	WssHLS string `json:"wssHLS"`
	// ws://ip:port/app/stream.live.ts
	WsTS string `json:"wsTS"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.ts
	WssTS string `json:"wssTS"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// ws://ip:port/app/stream.live.mp4
	WsFMP4 string `json:"wsFMP4"`
	// 播放地址，可能为空(主要看是否开启了服务端口)
	// wss://ip:port/app/stream.live.mp4
	WssFMP4 string `json:"wssFMP4"`
}

// init 主要初始化 url
func (m *PlayInfo) init(i *MediaInfo, public bool, token string) {
	m.MediaInfo = i
	// url
	var ip string
	if public {
		ip = m.Ser.PublicIP
	} else {
		ip = m.Ser.PrivateIP
	}
	// RTMP
	if m.Ser.RTMPPort != "" {
		if token != "" {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s?token=%s", ip, m.Ser.RTMPPort, m.App, m.Stream, token)
		} else {
			m.RTMP = fmt.Sprintf("rtmp://%s:%s/%s/%s", ip, m.Ser.RTMPPort, m.App, m.Stream)
		}
	}
	if m.Ser.RTMPSSLPort != "" {
		if token != "" {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s?token=%s", ip, m.Ser.RTMPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTMPs = fmt.Sprintf("rtmps://%s:%s/%s/%s", ip, m.Ser.RTMPSSLPort, m.App, m.Stream)
		}
	}
	if m.Ser.HTTPPort != "" {
		if token != "" {
			m.Flv = fmt.Sprintf("http://%s:%s/%s/%s.live.flv?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
			m.WsFlv = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.Flv = fmt.Sprintf("http://%s:%s/%s/%s.live.flv", ip, m.Ser.HTTPPort, m.App, m.Stream)
			m.WsFlv = fmt.Sprintf("ws://%s:%s/%s/%s.live.flv", ip, m.Ser.HTTPPort, m.App, m.Stream)
		}
	}
	if m.Ser.HTTPSSLPort != "" {
		if token != "" {
			m.Flvs = fmt.Sprintf("https://%s:%s/%s/%s.live.flv?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssFlv = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.Flvs = fmt.Sprintf("https://%s:%s/%s/%s.live.flv", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
			m.WssFlv = fmt.Sprintf("wss://%s:%s/%s/%s.live.flv", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// RTSP
	if m.Ser.RTSPPort != "" {
		if token != "" {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s?token=%s", ip, m.Ser.RTSPPort, m.App, m.Stream, token)
		} else {
			m.RTSP = fmt.Sprintf("rtsp://%s:%s/%s/%s", ip, m.Ser.RTSPPort, m.App, m.Stream)
		}
	}
	if m.Ser.RTSPSSLPort != "" {
		if token != "" {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s?token=%s", ip, m.Ser.RTSPSSLPort, m.App, m.Stream, token)
		} else {
			m.RTSPs = fmt.Sprintf("rtsps://%s:%s/%s/%s", ip, m.Ser.RTSPSSLPort, m.App, m.Stream)
		}
	}
	// HLS
	if m.Ser.HTTPPort != "" {
		if token != "" {
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
			m.WsHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.HLS = fmt.Sprintf("http://%s:%s/%s/%s/hls.m3u8", ip, m.Ser.HTTPPort, m.App, m.Stream)
			m.WsHLS = fmt.Sprintf("ws://%s:%s/%s/%s/hls.m3u8", ip, m.Ser.HTTPPort, m.App, m.Stream)
		}
	}
	if m.Ser.HTTPSSLPort != "" {
		if token != "" {
			m.HLSs = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.HLSs = fmt.Sprintf("https://%s:%s/%s/%s/hls.m3u8", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
			m.WssHLS = fmt.Sprintf("wss://%s:%s/%s/%s/hls.m3u8", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// TS
	if m.Ser.HTTPPort != "" {
		if token != "" {
			m.TS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
			m.WsTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.TS = fmt.Sprintf("http://%s:%s/%s/%s.live.ts", ip, m.Ser.HTTPPort, m.App, m.Stream)
			m.WsTS = fmt.Sprintf("ws://%s:%s/%s/%s.live.ts", ip, m.Ser.HTTPPort, m.App, m.Stream)
		}
	}
	if m.Ser.HTTPSSLPort != "" {
		if token != "" {
			m.TSs = fmt.Sprintf("https://%s:%s/%s/%s.live.ts?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.TSs = fmt.Sprintf("https://%s:%s/%s/%s.live.ts", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
			m.WssTS = fmt.Sprintf("wss://%s:%s/%s/%s.live.ts", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
	// MP4
	if m.Ser.HTTPPort != "" {
		if token != "" {
			m.FMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
			m.WsFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4?token=%s", ip, m.Ser.HTTPPort, m.App, m.Stream, token)
		} else {
			m.FMP4 = fmt.Sprintf("http://%s:%s/%s/%s.live.mp4", ip, m.Ser.HTTPPort, m.App, m.Stream)
			m.WsFMP4 = fmt.Sprintf("ws://%s:%s/%s/%s.live.mp4", ip, m.Ser.HTTPPort, m.App, m.Stream)
		}
	}
	if m.Ser.HTTPSSLPort != "" {
		if token != "" {
			m.FMP4s = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
			m.WssFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4?token=%s", ip, m.Ser.HTTPSSLPort, m.App, m.Stream, token)
		} else {
			m.FMP4s = fmt.Sprintf("https://%s:%s/%s/%s.live.mp4", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
			m.WssFMP4 = fmt.Sprintf("wss://%s:%s/%s/%s.live.mp4", ip, m.Ser.HTTPSSLPort, m.App, m.Stream)
		}
	}
}

// HasMediaInfo 返回内存中 app 和 stream 的媒体流缓存是否存在
func (s *Server) HasMediaInfo(app, stream string) bool {
	key := mediaInfoKey{App: app, Stream: stream}
	s.lock.RLock()
	m := s.mediaInfos[key]
	s.lock.RUnlock()
	return m != nil
}

// GetMediaInfo 返回内存中 app 和 stream 的媒体流缓存
func (s *Server) GetMediaInfo(app, stream string) *MediaInfo {
	key := mediaInfoKey{App: app, Stream: stream}
	s.lock.RLock()
	m := s.mediaInfos[key]
	s.lock.RUnlock()
	return m
}

// GetAppMediaInfo 返回内存中 app 的媒体流缓存
func (s *Server) GetAppMediaInfo(app string) []*MediaInfo {
	var ms []*MediaInfo
	s.lock.RLock()
	for _, m := range s.mediaInfos {
		if m.App == app {
			ms = append(ms, m)
		}
	}
	s.lock.RUnlock()
	return ms
}

// IsRecording 返回，是否有流，是否在录像
func (s *Server) IsRecording(app, stream string) (bool, bool) {
	// 查询
	key := mediaInfoKey{App: app, Stream: stream}
	s.lock.RLock()
	m := s.mediaInfos[key]
	s.lock.RUnlock()
	// 返回
	if m == nil {
		return false, false
	}
	return true, m.IsRecordingMP4
}

// // GetMediaInfoTimestamp 返回内存中 app 和 stream 的媒体流缓存是否存在
// func (s *Server) GetMediaInfoTimestamp(app, stream string) int64 {
// 	key := mediaInfoKey{App: app, Stream: stream}
// 	s.lock.RLock()
// 	media := s.mediaInfos[key]
// 	s.lock.RUnlock()
// 	if media == nil {
// 		return -1
// 	}
// 	return media.timestamp
// }
