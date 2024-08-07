package info

import (
	"context"
	"goutil/gb28181/request"
)

// SendInfoPlay 发送 Request-Info-Play 请求消息
func SendInfoPlay(ctx context.Context, m *Info) error {
	// 消息
	msg, addr := request.NewInfo(m.Device, m.ChannelID, m.Invite)
	// body
	m.encStartLline(&msg.Body, InfoMethodPlay)
	msg.Body.WriteString("Range: npt=now-\r\n")
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
