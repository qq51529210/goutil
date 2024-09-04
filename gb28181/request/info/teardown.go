package info

import (
	"context"
	"goutil/gb28181/request"
)

// SendInfoTeardown 发送 Request-Info-Teardown 请求消息
func SendInfoTeardown(ctx context.Context, m *Info) error {
	// 网络地址
	addr, err := m.Device.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := request.NewInfo(m.Device, addr.Network(), m.ChannelID, m.Invite)
	// body
	m.encStartLline(&msg.Body, InfoMethodTeardown)
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
