package zlm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ghttp "goutil/http"
	"net/http"
	"net/url"
)

type Boolean string

// 参数常量
const (
	Zero          = "0"
	One           = "1"
	Two           = "2"
	True  Boolean = One
	False Boolean = Zero
)

// 流的协议
const (
	RTMP = "rtmp"
	RTSP = "rtsp"
	HLS  = "hls"
	TS   = "ts"
	FMP4 = "fmp4"
)

const (
	// CodeOK 正确
	CodeOK = 0
	// VHost 默认的 vhost
	VHost = "__defaultVhost__"
	// api path 前缀
	apiPathPrefix = "/api/index"
)

// 定义一些错误以便全局使用，看名称猜意思
var (
	ErrServerNotAvailable = errors.New("server not available")
	ErrMediaNotFound      = errors.New("media not found")
)

type ResponseData interface {
	GetCode() int
	GetMsg() string
}

// CodeMsg on_xx 的返回值
type CodeMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg,omitempty"`
}

func (e *CodeMsg) GetCode() int {
	return e.Code
}

func (e *CodeMsg) GetMsg() string {
	return e.Msg
}

func (e *CodeMsg) Error() string {
	return fmt.Sprintf("code %d msg %s", e.Code, e.Msg)
}

// rtsp 拉流方式
type RTSPRTPType string

const (
	// RTSPRTPTypeTCP tcp
	RTSPRTPTypeTCP RTSPRTPType = Zero
	// RTSPRTPTypeTCP udp
	RTSPRTPTypeUDP RTSPRTPType = One
	// RTSPRTPTypeTCP 组播
	RTSPRTPTypeMulticast RTSPRTPType = Two
)

// rtp 推流的负载类型
type RTPPayloadType string

const (
	// RTPPayloadTypeES es
	RTPPayloadTypeES RTPPayloadType = Zero
	// RTPPSTypePS ps
	RTPPayloadTypePS RTPPayloadType = One
)

// 录像的类型
type RecordFileType string

const (
	//  hls
	RecordFileTypeHLS RecordFileType = Zero
	//  mp4
	RecordFileTypeMP4 RecordFileType = One
)

// Server 接口，用于请求
type Server interface {
	// http://localhost:8080
	BaseURL() string
	// 访问密钥
	Secret() string
	// 虚拟主机
	VHost() string
}

// Config 返回服务标识
type Config interface {
	ServerID() string
}

var (
	// Request 请求函数，可以替换
	Request = func(ctx context.Context, ser Server, apiPath string, query any, data ResponseData) error {
		return ghttp.JSONRequest(ctx, http.DefaultClient, http.MethodGet,
			ser.BaseURL()+apiPath, ghttp.Query(query, NewRequestQuery(ser)), nil,
			func(res *http.Response) error {
				// 响应码
				if res.StatusCode != http.StatusOK {
					return ghttp.StatusError(res.StatusCode)
				}
				// 解析
				if err := json.NewDecoder(res.Body).Decode(data); err != nil {
					return err
				}
				// 错误码
				if data.GetCode() != CodeOK {
					return fmt.Errorf("code:%d, msg:%s", data.GetCode(), data.GetMsg())
				}
				//
				return nil
			})
	}
)

// NewRequestQuery 填充 secret 和 vhost 返回
func NewRequestQuery(ser Server) url.Values {
	q := make(url.Values)
	q.Add("secret", ser.Secret())
	q.Add("vhost", ser.VHost())
	return q
}
