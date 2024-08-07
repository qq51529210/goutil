package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// TeleBoot 是 SendTeleBoot 的参数
type TeleBoot struct {
	Ser    *sip.Server
	Device request.Request
}

// SendTeleBoot 远程启动
func SendTeleBoot(ctx context.Context, m *TeleBoot) error {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	body.DeviceID = m.Device.GetRemoteID()
	body.SN = sip.GetSNString()
	body.TeleBoot = "Boot"
	// 请求
	return request.SendMessage(ctx, m.Ser, m.Device, &body, nil)
}
