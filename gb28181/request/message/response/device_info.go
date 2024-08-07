package response

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// DeviceInfo 是 SendDeviceInfo 的参数
type DeviceInfo struct {
	Ser          *sip.Server
	Cascade      request.Request
	SN           string
	DeviceID     string
	Manufacturer string
	Model        string
	Firmware     string
	Result       string
}

// SendDeviceInfo 设备信息结果应答
func SendDeviceInfo(ctx context.Context, m *DeviceInfo) error {
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdDeviceInfo
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.Manufacturer = m.Manufacturer
	body.Model = m.Model
	body.Firmware = m.Firmware
	body.Result = m.Result
	// 发送
	return request.SendMessage(ctx, m.Ser, m.Cascade, &body, m)
}
