package response

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// RecordInfo 是 SendRecordInfo 的参数
type RecordInfo struct {
	Ser      *sip.Server
	Cascade  request.Request
	SN       string
	DeviceID string
	Items    []*xml.Record
	// 追踪标识
	TraceID string
}

// SendRecordInfo 设备录像文件查询结果应答
func SendRecordInfo(ctx context.Context, m *RecordInfo) error {
	items := m.Items
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = xml.CmdRecordInfo
	body.SN = m.SN
	body.DeviceID = m.DeviceID
	body.SumNum = int64(len(items))
	body.RecordList = new(xml.MessageRecordList)
	// 发送
	if body.SumNum == 0 {
		return request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m)
	}
	// 3 条一起吧，差不多 1500 了
	body.RecordList.Item = make([]*xml.Record, 0, 3)
	body.RecordList.Num = int64(len(body.RecordList.Item))
	for len(items) > 0 {
		body.RecordList.Item = body.RecordList.Item[:0]
		body.RecordList.Num = 0
		for i := 0; i < cap(body.RecordList.Item); i++ {
			if len(items) < 1 {
				break
			}
			body.RecordList.Item = append(body.RecordList.Item, items[0])
			body.RecordList.Num++
			items = items[1:]
			// 这个要替换一下
			body.RecordList.Item[i].DeviceID = m.DeviceID
		}
		// 这个要替换一下
		if err := request.SendMessage(ctx, m.TraceID, m.Ser, m.Cascade, &body, m); err != nil {
			return err
		}
	}
	return nil
}
