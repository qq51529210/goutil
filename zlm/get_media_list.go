package zlm

import (
	"context"
)

const (
	OriginTypeUnknown = iota
	OriginTypeRtmpPush
	OriginTypeRtspPush
	OriginTypeRtpPush
	OriginTypePull
	OriginTypeFFmpegPull
	OriginTypeMp4Vod
	OriginTypeDeviceChn
	OriginTypeRtcPush
)

// MediaListDataSock 是 MediaListData.OriginSock
type MediaListDataSock struct {
	LocalIP   string `json:"local_ip"`
	LocalPort int    `json:"local_port"`
	PeerIP    string `json:"peer_ip"`
	PeerPort  int    `json:"peer_port"`
}

// MediaListData 是 GetMediaList 的返回值
type MediaListData struct {
	// 虚拟主机名
	VHost string `json:"vhost"`
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// unix 系统时间戳，单位秒
	CreateStamp int64 `json:"createStamp"`
	// 是否正在录制 hls
	IsRecordingHLS bool `json:"isRecordingHLS"`
	// 是否正在录制 mp4
	IsRecordingMP4 bool `json:"isRecordingMP4"`
	// 自定义数据
	UserData string `json:"userdata"`
	// 音视频轨道
	Tracks []map[string]any `json:"tracks"`
	// 连接信息
	OriginSock *MediaListDataSock `json:"originSock"`
	// 观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv
	TotalReaderCount int64  `json:"totalReaderCount"`
	OriginType       int64  `json:"originType"`
	OriginURL        string `json:"originUrl"`
}

// ParseTrack 区分出音/视频轨道
func ParseTrack(tracks []map[string]any) (vs []*MediaInfoVideoTrack, as []*MediaInfoAudioTrack) {
	// 这里假定，返回的数据格式没有错误
	for _, track := range tracks {
		switch int64(track["codec_type"].(float64)) {
		case 0:
			t := new(MediaInfoVideoTrack)
			t.CodecName = track["codec_id_name"].(string)
			t.FPS = int(track["fps"].(float64))
			t.Height = int(track["height"].(float64))
			t.Width = int(track["width"].(float64))
			vs = append(vs, t)
		case 1:
			t := new(MediaInfoAudioTrack)
			t.CodecName = track["codec_id_name"].(string)
			t.SampleBit = int(track["sample_bit"].(float64))
			t.SampleRate = int(track["sample_rate"].(float64))
			as = append(as, t)
		}
	}
	return
}

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
	// 协议
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
}

// GetMediaListRes 是 GetMediaList 的返回值
type GetMediaListRes struct {
	CodeMsg
	Data []*MediaListData `json:"data"`
}

const (
	GetMediaListPath = apiPathPrefix + "/getMediaList"
)

// GetMediaList 调用 /index/api/getMediaList ，返回媒体流列表
func GetMediaList(ctx context.Context, ser Server, req *GetMediaListReq, res *GetMediaListRes) error {
	return Request(ctx, ser, GetMediaListPath, req, res)
}
