package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

type RecordCmd string

// 命令
const (
	RecordCmdStart RecordCmd = "Record"
	RecordCmdStop  RecordCmd = "StopRecord"
)

// Record 是 SendRecord 的参数
type Record struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	Cmd       RecordCmd
	// 追踪标识
	TraceID string
}

// SendRecord 录像控制
func SendRecord(ctx context.Context, m *Record) (string, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	// 这个应该是用的通道编号
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.RecordCmd = string(m.Cmd)
	// 请求
	var res request.XMLResult
	return res.Result, request.SendReplyMessage(ctx, m.TraceID, m.Ser, m.Device, &body, &res)
}
