package info

import (
	"context"
	"fmt"
	"goutil/gb28181/request"
)

// SendInfoRange 发送 Request-Info-Range 请求消息
func SendInfoRange(ctx context.Context, m *Info, sec int64) error {
	// 消息
	msg := request.NewInfo(m.Device, m.ChannelID, m.Invite)
	// 网络地址
	addr, err := m.Device.GetNetAddr()
	if err != nil {
		return err
	}
	// body
	m.encStartLline(&msg.Body, InfoMethodPlay)
	if sec > 0 {
		fmt.Fprintf(&msg.Body, "Range: npt=%d-\r\n", sec)
	} else {
		msg.Body.WriteString("Range: npt=now-\r\n")
	}
	msg.Body.WriteString("\r\n")
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
