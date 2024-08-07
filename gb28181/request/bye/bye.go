package bye

import (
	"context"
	"goutil/gb28181/request"
	"goutil/sip"
)

// Bye 是 SendBye 的参数
type Bye struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	Invite    request.Invite
}

// SendBye 发送 Request-Bye 请求消息
func SendBye(ctx context.Context, m *Bye) error {
	// 消息
	msg, addr, err := request.New(m.Device, m.ChannelID, sip.MethodBye, "")
	if err != nil {
		return err
	}
	// 恢复
	msg.Header.From.Tag = m.Invite.GetFromTag()
	msg.Header.To.Tag = m.Invite.GetToTag()
	msg.Header.CallID = m.Invite.GetCallID()
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
