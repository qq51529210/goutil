package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

type GuardCmd string

// 命令
const (
	SetGuardCmd   GuardCmd = "SetGuard"
	ResetGuardCmd GuardCmd = "ResetGuard"
)

// Guard 是 SendGuard 的参数
type Guard struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	// 命令
	Cmd GuardCmd
	// 追踪标识
	TraceID string
}

// SendGuard 布防/撤防
func SendGuard(ctx context.Context, m *Guard) (string, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	// 这个应该是用的通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.GuardCmd = string(m.Cmd)
	// 请求
	var res request.XMLResult
	return res.Result, request.SendReplyMessage(ctx, m.TraceID, m.Ser, m.Device, &body, &res)
}
