package query

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Catalog 是 SendCatalog 的参数
type Catalog struct {
	Ser       *sip.Server
	Device    request.Request
	StartTime string
	EndTime   string
	// 已接收的个数
	Item []*xml.Device
	// 追踪标识
	TraceID string
}

// SendCatalog 目录查询
func SendCatalog(ctx context.Context, m *Catalog) ([]*xml.Device, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdCatalog
	body.DeviceID = m.Device.GetToID()
	body.SN = sip.GetSNString()
	body.StartTime = m.StartTime
	body.EndTime = m.EndTime
	// 请求
	return m.Item, request.SendReplyMessage(ctx, m.TraceID, m.Ser, m.Device, &body, m)
}
