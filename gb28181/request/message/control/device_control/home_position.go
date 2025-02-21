package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// HomePosition 是 SendHomePosition 的参数
type HomePosition struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	//
	Data *xml.MessageHomePosition
	// 追踪标识
	TraceID string
}

// SendHomePosition 设置看守位
func SendHomePosition(ctx context.Context, m *HomePosition) (string, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	// 通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.HomePosition = m.Data
	// 请求
	var res request.XMLResult
	return res.Result, request.SendReplyMessage(ctx, m.TraceID, m.Ser, m.Device, &body, &res)
}
