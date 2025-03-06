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

// CatalogOne 是 SendCatalogOne 的参数
type CatalogOne struct {
	Ser      *sip.Server
	Cascade  request.Request
	SN       string
	DeviceID string
	Total    int64
	Item     *xml.Device
	// 追踪标识
	TraceID string
}

// SendCatalogOne 目录查询结果应答
func SendCatalogOne(ctx context.Context, m *CatalogOne) error {
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdCatalog
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.SumNum = m.Total
	body.DeviceList = new(xml.MessageDeviceList)
	// 一条一条的发送，2 条就可能超 1500 了
	body.DeviceList.Item = make([]*xml.Device, 1)
	body.DeviceList.Num = int64(len(body.DeviceList.Item))
	body.DeviceList.Item[0] = m.Item
	if err := request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m); err != nil {
		return err
	}
	return nil
}
