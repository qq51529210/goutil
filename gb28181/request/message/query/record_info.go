package query

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// RecordInfo 是 SendRecord 的参数
type RecordInfo struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	// 录像起始时间
	StartTime string
	// 录像终止时间
	EndTime string
	// 录像产生类型， time/alarm/manual/all
	Type string
	// 已接收的个数
	Item []*xml.Record
	// 总数
	Total int64
	// 因为可能需要在接收到设备的数据后
	// 立即转发给上级，这里可以带相关的数据
	Data any
}

// SendRecordInfo 查询录像文件
func SendRecordInfo(ctx context.Context, m *RecordInfo) ([]*xml.Record, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdRecordInfo
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.StartTime = m.StartTime
	body.EndTime = m.EndTime
	body.Type = m.Type
	if body.Type == "" {
		body.Type = "all"
	}
	// 请求
	return m.Item, request.SendReplyMessage(ctx, m.Ser, m.Device, &body, m)
}
