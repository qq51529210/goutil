package devicecontrol

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Alarm 是 SendAlarm 的参数
type Alarm struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	// 复位报警的报警方式
	// 1: 电话报警
	// 2: 设备报警
	// 3: 短信报警
	// 4: GPS 报警
	// 5: 视频报警
	// 6: 设备故障报警
	// 7: 其他报警
	AlarmMethod string `json:"alarmMethod" binding:"omitempty,oneof=1 2 3 4 5 6 7"`
	// 复位报警的报警类型
	// 报警方式为2时，不携带 AlarmType 为默认的报警设备报警
	// 携带 AlarmType 取值及对应报警类型如下:
	// 1: 视频丢失报警
	// 2: 设备防拆报警
	// 3: 存储设备磁盘满报警
	// 4: 设备高温报警
	// 5: 设备低温报警
	// 报警方式为5时,取值如下:
	// 1: 人工视频报警
	// 2: 运动目标检测报警
	// 3: 遗留物检测报警
	// 4: 物体移除检测报警
	// 5: 绊线检测报警
	// 6: 入侵检测报警
	// 7: 逆行检测报警
	// 8: 徘徊检测报警
	// 9: 流量统计报警
	// 10: 密度检测报警
	// 11: 视频异常检测报警
	// 12: 快速移动报警
	// 报警方式为6时,取值如下:
	// 1: 存储设备磁盘故障报警
	// 2: 存储设备风扇故障报警
	AlarmType string `json:"alarmType" binding:"omitempty,oneof=1 2 3 4 5 6 7 8 9 10 11 12"`
}

// SendAlarm 报警复位
func SendAlarm(ctx context.Context, m *Alarm) (string, error) {
	// 消息
	var body xml.Message
	body.XMLName.Local = xml.TypeControl
	body.CmdType = xml.CmdDeviceControl
	body.DeviceID = m.ChannelID
	body.SN = sip.GetSNString()
	body.AlarmCmd = "ResetAlarm"
	body.Info = new(xml.MessageInfo)
	body.Info.AlarmMethod = m.AlarmMethod
	body.Info.AlarmType = m.AlarmType
	// 请求
	var res request.XMLResult
	return res.Result, request.SendReplyMessage(ctx, m.Ser, m.Device, &body, &res)
}
