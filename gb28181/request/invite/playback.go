package invite

import (
	"context"
	"goutil/sdp"
)

// SendPlayback 录像回放
func SendPlayback(ctx context.Context, m *Invite, startTime, endTime int64) error {
	// sdp
	var sdp sdp.Session
	// 消息
	msg, addr, err := m.Message(&sdp, InvitePlayback, "")
	if err != nil {
		return err
	}
	sdp.T.Start = startTime
	sdp.T.Stop = endTime
	sdp.FormatTo(&msg.Body)
	// 请求
	return m.Ser.RequestWithContext(ctx, m.TraceID, msg, addr, m)
}
