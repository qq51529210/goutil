package info

import (
	"context"
	"fmt"
	"goutil/gb28181/request"
)

// SendInfoScale 发送 Request-Info-Scale 请求消息
func SendInfoScale(ctx context.Context, m *Info, scale string) error {
	// 消息
	msg, addr := request.NewInfo(m.Device, m.ChannelID, m.Invite)
	// body
	m.encStartLline(&msg.Body, InfoMethodPlay)
	fmt.Fprintf(&msg.Body, "Scale: %s\r\n", scale)
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
