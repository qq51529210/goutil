package info

import (
	"context"
	"goutil/gb28181/request"
)

// SendInfoPause 回放暂停
func SendInfoPause(ctx context.Context, m *Info) error {
	// 网络地址
	addr, err := m.Device.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := request.NewInfo(m.Device, addr.Network(), m.ChannelID, m.Invite)
	// body
	m.encStartLline(&msg.Body, InfoMethodPause)
	msg.Body.WriteString("PauseTime: now\r\n")
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, m.TraceID, msg, addr, m)
}
