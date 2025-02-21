package response

import (
	"context"
	"goutil/gb28181"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// DeviceStatus 是 SendDeviceStatus 的参数
type DeviceStatus struct {
	Ser      *sip.Server
	Cascade  request.Request
	SN       string
	DeviceID string
	Result   string
	Online   string
	Status   string
	// 追踪标识
	TraceID string
}

// SendDeviceStatus 设备状态结果应答
func SendDeviceStatus(ctx context.Context, m *DeviceStatus) error {
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdDeviceStatus
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.Result = m.Result
	body.Status = m.Status
	body.Online = m.Online
	body.DeviceTime = gb28181.Time()
	// 发送
	return request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m)
}
