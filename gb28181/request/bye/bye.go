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
	return request.SendBye(ctx, m.Ser, m.Device, m.ChannelID, m.Invite, m)
}
