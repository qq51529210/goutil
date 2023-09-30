package zlm

import (
	"context"
	"mms/db"
	"sync/atomic"
	"time"
)

// OnServerKeepaliveReq 表示 on_server_keepalive 提交的数据
type OnServerKeepaliveReq struct {
	Data          *OnServerKeepaliveDataModel `json:"data"`
	MediaServerID string                      `json:"mediaServerId"`
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
func OnServerKeepalive(ctx context.Context, req *OnServerKeepaliveReq) {
	// 获取实例
	ser := GetServer(req.MediaServerID)
	if ser == nil || atomic.LoadInt32(&ser.ok) != 1 {
		return
	}
	now := time.Now()
	timestamp := now.Unix()
	ser.keepaliveTime = &now
	ser.KeepaliveTime = timestamp
	*ser.Online = db.True
	atomic.StoreInt32(&ser.updateDB, 1)
}
