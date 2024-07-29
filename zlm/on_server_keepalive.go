package zlm

import (
	"context"
)

// OnServerKeepaliveReq 表示 on_server_keepalive 提交的数据
type OnServerKeepaliveReq struct {
	// 数据
	Data *OnServerKeepaliveDataModel `json:"data"`
	// 服务标识
	MediaServerID string `json:"mediaServerId"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnServerKeepaliveDataModel 是 OnServerKeepaliveReq 的 data 字段
type OnServerKeepaliveDataModel struct {
	Buffer                int
	BufferLikeString      int
	BufferList            int
	BufferRaw             int
	Frame                 int
	FrameImp              int
	MediaSource           int
	MultiMediaSourceMuxer int
	RTMPPacket            int `json:"RtmpPacket"`
	RTPPacket             int `json:"RtpPacket"`
	Socket                int
	TCPClient             int `json:"TcpClient"`
	TCPServer             int `json:"TcpServer"`
	TCPSession            int `json:"TcpSession"`
	UDPServer             int `json:"UdpServer"`
	UDPSession            int `json:"UdpSession"`
}

// OnServerKeepalive 处理 zlm 的 on_server_keepalive 回调
func OnServerKeepalive(ctx context.Context, req *OnServerKeepaliveReq, res *CodeMsg) {
}
