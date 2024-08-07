package info

import (
	"bytes"
	"context"
	"fmt"
	"goutil/gb28181/request"
	"goutil/sip"
	"io"
	"net"
)

const (
	InfoMethodPlay     = "PLAY"
	InfoMethodPause    = "PAUSE"
	InfoMethodTeardown = "TEARDOWN"
)

// Info 是 SendInfo 的参数
type Info struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	Invite    request.Invite
}

func (m *Info) Message() (*sip.Message, net.Addr, error) {
	// 消息
	msg, addr, err := request.New(m.Device, m.ChannelID, sip.MethodInfo, request.ContentTypeMANSRTSP)
	if err != nil {
		return nil, nil, err
	}
	// 恢复
	msg.Header.From.Tag = m.Invite.GetFromTag()
	msg.Header.To.Tag = m.Invite.GetToTag()
	msg.Header.CallID = m.Invite.GetCallID()
	//
	return msg, addr, nil
}

func (m *Info) encStartLline(w *bytes.Buffer, cmd string) {
	fmt.Fprintf(w, "%s RTSP/1.0\r\nCSeq: %d\r\n", cmd, sip.GetSN())
}

// SendInfoRaw 用于级联转发，因为 body 的数据不变
func SendInfoRaw(ctx context.Context, m *Info, body io.Reader) error {
	// 消息
	msg, addr, err := m.Message()
	if err != nil {
		return err
	}
	// body
	if _, err := io.Copy(&msg.Body, body); err != nil {
		return err
	}
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
