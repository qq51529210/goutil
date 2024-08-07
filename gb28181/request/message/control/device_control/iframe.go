package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	gbxml "goutil/gb28181/xml"
	"goutil/sip"
)

// IFame 是 SendIFame 的参数
type IFame struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
}

// SendIFame 强制关键帧
func SendIFame(ctx context.Context, m *IFame) error {
	// 消息
	var body gbxml.Message
	body.XMLName.Local = gbxml.TypeControl
	body.CmdType = gbxml.CmdDeviceControl
	// 通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.IFameCmd = "Send"
	// 请求
	return request.SendMessage(ctx, m.Ser, m.Device, &body, nil)
}
