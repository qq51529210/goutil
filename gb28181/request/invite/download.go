package invite

import (
	"context"
	"goutil/sdp"
)

// SendDownload 录像下载
func SendDownload(ctx context.Context, m *Invite, startTime, endTime int64, speed string) error {
	// sdp
	var sdp sdp.Session
	// 消息
	msg, addr, err := m.Message(&sdp, InviteDownload, speed)
	if err != nil {
		return err
	}
	sdp.T.Start = startTime
	sdp.T.Stop = endTime
	sdp.FormatTo(&msg.Body)
	// 请求
	return m.Ser.RequestWithContext(ctx, m.TraceID, msg, addr, m)
}
