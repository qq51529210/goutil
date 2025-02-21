package query

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// PresetQuery 是 SendPresetQuery 的参数
type PresetQuery struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	// 已接收的个数
	Item []*xml.MessagePresetListItem
	// 总数
	Total int64
	// 追踪标识
	TraceID string
}

// SendPresetQuery 查询预置位
func SendPresetQuery(ctx context.Context, m *PresetQuery) ([]*xml.MessagePresetListItem, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdPresetQuery
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	// 请求
	return m.Item, request.SendReplyMessage(ctx, m.TraceID, m.Ser, m.Device, &body, m)
}
