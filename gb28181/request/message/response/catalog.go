package response

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Catalog 是 SendCatalog 的参数
type Catalog struct {
	Ser      *sip.Server
	Cascade  request.Request
	SN       string
	DeviceID string
	Items    []*xml.Device
	// 追踪标识
	TraceID string
}

// SendCatalog 目录查询结果应答
func SendCatalog(ctx context.Context, m *Catalog) error {
	items := m.Items
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdCatalog
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.SumNum = int64(len(items))
	body.DeviceList = new(xml.MessageDeviceList)
	// 发送
	if body.SumNum == 0 {
		return request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m)
	}
	// 一条一条的发送，2 条就可能超 1500 了
	body.DeviceList.Item = make([]*xml.Device, 1)
	body.DeviceList.Num = int64(len(body.DeviceList.Item))
	for len(items) > 0 {
		body.DeviceList.Item[0] = items[0]
		if err := request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m); err != nil {
			return err
		}
		items = items[1:]
	}
	return nil
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
