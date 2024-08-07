package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// DragZoom 是 SendDragZoom 的参数
type DragZoom struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	//
	Data *xml.MessageDragZoom
}

// SendDragZoomIn 拉框放大
func SendDragZoomIn(ctx context.Context, m *DragZoom) error {
	return sendDragZoom(ctx, m, &xml.Message{
		DragZoomIn: m.Data,
	})
}

// SendDragZoomOut 拉框缩小
func SendDragZoomOut(ctx context.Context, m *DragZoom) error {
	return sendDragZoom(ctx, m, &xml.Message{
		DragZoomOut: m.Data,
	})
}

// sendDragZoom 封装代码
func sendDragZoom(ctx context.Context, m *DragZoom, body *xml.Message) error {
	// 消息
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	// 通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	// 请求
	return request.SendMessage(ctx, m.Ser, m.Device, body, nil)
}
