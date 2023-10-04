package zlm

import (
	"context"
)

// GetMediaListReq 是 GetMediaList 的参数
type GetMediaListReq struct {
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
	// OriginSock       *mediaInfoOriginSock `json:"originSock"`
	// AliveSecond      int                  `json:"aliveSecond"`
	// ReaderCount      int                  `json:"readerCount"`
	// OriginType       int                  `json:"originType"`
	// BytesSpeed       int                  `json:"bytesSpeed"`
}

const (
	apiGetMediaList = "getMediaList"
)

// GetMediaList 调用 /index/api/getMediaList
// 返回媒体流列表
func (s *Server) GetMediaList(ctx context.Context, req *GetMediaListReq) ([]*MediaListData, error) {
	// 请求
	var res getMediaListRes
	err := httpCallRes(ctx, s, apiGetMediaList, req, &res)
	if err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiGetMediaList
		return nil, &res.apiError
	}
	// 更新内存
	player := 0
	infos := make(map[mediaInfoKey]*MediaInfo)
	for _, d := range res.Data {
		// 使用 rtmp 的即可
		if d.Schema != RTMP {
			continue
		}
		info := new(MediaInfo)
		info.init(s, d)
		infos[mediaInfoKey{
			App:    d.App,
			Stream: d.Stream,
		}] = info
		//
		player += int(d.TotalReaderCount)
	}
	// 赋值
	s.mediaInfos.Lock()
	s.mediaInfos.D = infos
	s.mediaInfos.ResetSlice()
	s.player = player
	s.mediaInfos.Unlock()
	//
	return res.Data, nil
}
