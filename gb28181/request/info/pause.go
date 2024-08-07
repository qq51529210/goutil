package info

import (
	"context"
)

// SendInfoPause 回放暂停
func SendInfoPause(ctx context.Context, m *Info) error {
	// 消息
	msg, addr, err := m.Message()
	if err != nil {
		return err
	}
	// body
	m.encStartLline(&msg.Body, InfoMethodPause)
	msg.Body.WriteString("PauseTime: now\r\n")
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
