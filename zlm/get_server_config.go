package zlm

import (
	"context"
	"errors"
)

var (
	// ErrConfig 错误的服务配置
	ErrConfig = errors.New("error config")
)

// getServerConfigRes 是 GetServerConfig 的返回值
type getServerConfigRes struct {
	apiError
	Data []map[string]string `json:"data"`
}

const (
	apiGetServerConfig = "getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig
func (s *Server) GetServerConfig(ctx context.Context) (map[string]string, error) {
	// 请求
	var res getServerConfigRes
	err := httpCallRes[any](ctx, s, apiGetServerConfig, nil, &res)
	if err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiGetServerConfig
		return nil, &res.apiError
	}
	// 找到自己的配置
	for _, d := range res.Data {
		if d["general.mediaServerId"] == s.ID {
			return d, nil
		}
	}
	// 没有配置，流媒体服务有问题
	return nil, ErrConfig
}

// // 端口
// s.RTMPPort = d["rtmp.port"]
// s.RTMPPort = d["rtmp.port"]
// s.RTMPSSLPort = d["rtmp.sslport"]
// s.RTSPPort = d["rtsp.port"]
// s.RTSPSSLPort = d["rtsp.sslport"]
// s.HTTPPort = d["http.port"]
// s.HTTPSSLPort = d["http.sslport"]
// s.RTCUDPPort = d["rtc.port"]
// s.RTCTCPPort = d["rtc.tcpPort"]
// s.RTPProxyPort = d["rtp_proxy.port"]
// p := strings.Split(d["rtp_proxy.port_range"], "-")
// if len(p) > 0 {
// 	s.RTPProxyPortMin = p[0]
// }
// if len(p) > 1 {
// 	s.RTPProxyPortMax = p[1]
// }
// // ffmpeg 命令
// var ffmpegCmd []string
// for k := range d {
// 	if strings.HasPrefix(k, "ffmpeg.cmd") {
// 		ffmpegCmd = append(ffmpegCmd, k)
// 	}
// }
// if len(ffmpegCmd) > 0 {
// 	_d, _ := json.Marshal(ffmpegCmd)
// 	s.FFMPEGCmd = string(_d)
// }
// // rtp 超时
// n, err := strconv.ParseInt(d["rtp_proxy.timeoutSec"], 10, 32)
// if err != nil {
// 	return fmt.Errorf("rtp_proxy.timeoutSec %w", err)
// }
// s.RTPTimeout = int32(n)
//
