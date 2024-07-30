package zlm

import (
	"context"
	"errors"
)

var (
	// ErrConfig 错误的服务配置
	ErrConfig = errors.New("error config")
)

// GetServerConfigReq 是 GetServerConfig 参数
type GetServerConfigReq struct {
	// 服务标识，用于筛选配置
	ID string
}

type getServerConfigRes[M Config] struct {
	CodeMsg
	Data []M `json:"data"`
}

// GetServerConfigRes 是 GetServerConfig 的返回值
// M 是结构体
//
//	type config struct {
//		MediaServerID      string `json:"general.mediaServerId"`
//		RTMPPort           string `json:"rtmp.port"`
//		RTMPSSLPort        string `json:"rtmp.sslport"`
//		RTSPPort           string `json:"rtsp.port"`
//		RTSPSSLPort        string `json:"rtsp.sslport"`
//		HTTPPort           string `json:"http.port"`
//		HTTPSSLPort        string `json:"http.sslport"`
//		RTCUDPPort         string `json:"rtc.port"`
//		RTCTCPPort         string `json:"rtc.tcpPort"`
//		AliveInterval      string `json:"hook.alive_interval"`
//		RTPProxyPortRange  string `json:"rtp_proxy.port_range"`
//		RTPProxyTimeoutSec string `json:"rtp_proxy.timeoutSec"`
//		MaxStreamWaitMS    string `json:"general.maxStreamWaitMS"`
//	}
//
//	func (m *config) ID() string {
//		return m.MediaServerID
//	}
type GetServerConfigRes[M Config] struct {
	CodeMsg
	Data M `json:"data"`
}

const (
	GetServerConfigPath = apiPathPrefix + "/getServerConfig"
)

// GetServerConfig 调用 /index/api/getServerConfig ，查询配置
func GetServerConfig[M Config](ctx context.Context, ser Server, req *GetServerConfigReq, res *GetServerConfigRes[M]) error {
	var _res getServerConfigRes[M]
	if err := Request(ctx, ser, GetServerConfigPath, req, &_res); err != nil {
		return err
	}
	res.CodeMsg = _res.CodeMsg
	// 筛选
	for i := 0; i < len(_res.Data); i++ {
		if _res.Data[i].ServerID() == req.ID {
			res.Data = _res.Data[i]
			break
		}
	}
	return nil
}
