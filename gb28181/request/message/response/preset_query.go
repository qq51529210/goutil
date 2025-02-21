package response

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// PresetQuery 是 SendPresetQuery 的参数
type PresetQuery struct {
	Ser      *sip.Server
	Cascade  request.Request
	SN       string
	DeviceID string
	Items    []*xml.MessagePresetListItem
	// 追踪标识
	TraceID string
}

// SendPresetQuery 预置位查询结果应答
func SendPresetQuery(ctx context.Context, m *PresetQuery) error {
	items := m.Items
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdPresetQuery
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.SumNum = int64(len(items))
	body.PresetList = new(xml.MessagePresetList)
	body.PresetList.Num = int64(len(items))
	body.PresetList.Item = items
	// 发送
	return request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m)
}
