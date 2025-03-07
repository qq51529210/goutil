package notify

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Keepalive 是 SendKeepalive 的参数
type Keepalive struct {
	Ser     *sip.Server
	Cascade request.Request
	// 追踪标识
	TraceID string
}

// SendKeepalive 心跳
func SendKeepalive(ctx context.Context, m *Keepalive) error {
	// 消息体
	var body xml.Message
	body.XMLName.Local = xml.TypeNotify
	body.CmdType = xml.CmdKeepalive
	body.SN = sip.GetSNString()
	body.DeviceID = m.Cascade.GetFromID()
	body.Status = "OK"
	// 请求
	return request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m)
}
