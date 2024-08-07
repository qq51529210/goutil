package invite

import (
	"context"
	"goutil/sdp"
)

// SendPlay 直播
func SendPlay(ctx context.Context, m *Invite) error {
	// sdp
	var sdp sdp.Session
	// 消息
	msg, addr, err := m.Message(&sdp, InvitePlay, "")
	if err != nil {
		return err
	}
	sdp.FormatTo(&msg.Body)
	// 请求
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
