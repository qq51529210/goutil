package info

import (
	"bytes"
	"context"
	"fmt"
	"goutil/gb28181/request"
	"goutil/sip"
	"io"
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

func (m *Info) encStartLline(w *bytes.Buffer, cmd string) {
	fmt.Fprintf(w, "%s RTSP/1.0\r\nCSeq: %d\r\n", cmd, sip.GetSN())
}

// SendInfoRaw 用于级联转发，因为 body 的数据不变
func SendInfoRaw(ctx context.Context, m *Info, body io.Reader) error {
	return request.SendInfo(ctx, m.Ser, m.Device, m.ChannelID, m.Invite, body, m)
}
