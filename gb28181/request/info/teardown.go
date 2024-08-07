package info

import (
	"context"
)

// SendInfoTeardown 发送 Request-Info-Teardown 请求消息
func SendInfoTeardown(ctx context.Context, m *Info) error {
	// 消息
	msg, addr, err := m.Message()
	if err != nil {
		return err
	}
	// body
	m.encStartLline(&msg.Body, InfoMethodTeardown)
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
