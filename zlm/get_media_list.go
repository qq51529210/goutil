package zlm

import (
	"context"
)

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 筛选虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 筛选协议，例如 rtsp或rtmp
	Schema string `query:"schema"`
	// 筛选应用名，例如 live
	App string `query:"app"`
	// 筛选流id，例如 test
	Stream string `query:"stream"`
}

// getMediaListRes 是 GetMediaList 的返回值
type getMediaListRes struct {
	apiError
	Data []*MediaListData `json:"data"`
}

// MediaListData 是 GetMediaList 的返回值
type MediaListData struct {
	// 虚拟主机名
	VHost string `json:"vhost"`
	// 协议
	Schema string `json:"schema"`
	// 应用名
	App string `json:"app"`
	// 协议
	Stream string `json:"stream"`
	// unix 系统时间戳，单位秒
	CreateStamp int64 `json:"createStamp"`
	// 是否正在录制 hls
	IsRecordingHLS bool `json:"isRecordingHLS"`
	// 是否正在录制 mp4
	IsRecordingMP4 bool `json:"isRecordingMP4"`
	// 音视频轨道
	Tracks []map[string]any `json:"tracks"`
	// 观看总人数，包括hls/rtsp/rtmp/http-flv/ws-flv
	TotalReaderCount int64  `json:"totalReaderCount"`
	OriginTypeStr    string `json:"originTypeStr"`
	OriginURL        string `json:"originUrl"`
}

// InitMediaInfo 填充 m
func (d *MediaListData) InitMediaInfo(m *MediaInfo) {
	m.App = d.App
	m.Stream = d.Stream
	m.Tracks = d.Tracks
	m.IsRecordingHLS = d.IsRecordingHLS
	m.IsRecordingMP4 = d.IsRecordingMP4
	m.Timestamp = d.CreateStamp
	m.TotalReaderCount = d.TotalReaderCount
	m.OriginTypeStr = d.OriginTypeStr
	m.OriginURL = d.OriginURL
	m.Video, m.Audio = ParseTrack(d.Tracks)
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

const (
	apiGetMediaList = "getMediaList"
)

// GetMediaList 调用 /index/api/getMediaList
// 返回媒体流列表
func GetMediaList(ctx context.Context, req *GetMediaListReq) ([]*MediaListData, error) {
	// 请求
	var res getMediaListRes
	if err := request(ctx, req.BaseURL, apiGetMediaList, req, &res); err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiGetMediaList
		return nil, &res.apiError
	}
	//
	return res.Data, nil
}
