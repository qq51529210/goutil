package zlm

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
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

// Config 用于保存流媒体服务配置
type Config struct {
	// 配置
	d map[string]string
	// 配置的 rtmp 端口
	RTMPPort string
	// 配置的 rtmps 端口
	RTMPSSLPort string
	// 配置的 rtsp 端口
	RTSPPort string
	// 配置的 rtsps 端口
	RTSPSSLPort string
	// 配置的 http 端口
	HTTPPort string
	// 配置的 https 端口
	HTTPSSLPort string
	// 配置的 rtc 端口
	RTCUDPPort string
	// 配置的 rtc 端口
	RTCTCPPort string
	// 配置的 rtp proxy 端口
	RTPProxyPort string
	// 配置的 rtp proxy 端口范围
	RTPProxyPortRange string
	// ffmpeg cmd
	FFMPEGCMD map[string]string
}

// GetServerConfig 调用 /index/api/getServerConfig
// 获取服务器配置
func (s *Server) GetServerConfig(ctx context.Context) error {
	// 请求
	var res getServerConfigRes
	err := httpCallRes[any](ctx, s, apiGetServerConfig, nil, &res)
	if err != nil {
		return err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.SerID = s.ID
		res.apiError.Path = apiGetServerConfig
		return &res.apiError
	}
	// 找到自己的配置
	for _, d := range res.Data {
		if d["general.mediaServerId"] == s.ID {
			s.updateConfig(d)
			break
		}
	}
	// 没有配置，流媒体服务有问题
	if s.Config == nil {
		return ErrConfig
	}
	//
	return nil
}

// updateConfig 更新配置
func (s *Server) updateConfig(data map[string]string) {
	cfg := new(Config)
	cfg.d = data
	cfg.FFMPEGCMD = map[string]string{}
	// ffmpeg cmd
	for k, v := range data {
		if strings.HasPrefix(k, "ffmpeg.cmd") {
			cfg.FFMPEGCMD[k] = v
		}
	}
	// 端口
	cfg.RTMPPort = data["rtmp.port"]
	cfg.RTMPSSLPort = data["rtmp.sslport"]
	cfg.RTSPPort = data["rtsp.port"]
	cfg.RTSPSSLPort = data["rtsp.sslport"]
	cfg.HTTPPort = data["http.port"]
	cfg.HTTPSSLPort = data["http.sslport"]
	cfg.RTCUDPPort = data["rtc.port"]
	cfg.RTCTCPPort = data["rtc.tcpPort"]
	cfg.RTPProxyPort = data["rtp_proxy.port"]
	cfg.RTPProxyPortRange = data["rtp_proxy.port_range"]
	// 心跳
	s.keepaliveTimeout = time.Minute
	str := data["hook.alive_interval"]
	if str != "" {
		n, err := strconv.ParseFloat(str, 64)
		if err == nil {
			s.keepaliveTimeout = time.Duration(n * float64(time.Second))
		}
	}
	//
	s.Config = cfg
}
