package subscribe

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Alarm 是 SendAlarm 的参数
type Alarm struct {
	Ser    *sip.Server
	Device request.Request
	// 过期时间
	Expire int64
	// 报警终止级别
	// 0:全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	StartAlarmPriority string `json:"startAlarmPriority" binding:"omitempty,oneof=0 1 2 3 4"`
	// 报警终止级别
	// 0:全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	EndAlarmPriority string `json:"endAlarmPriority" binding:"omitempty,oneof=0 1 2 3 4"`
	// 报警方式条件
	// 0:全部
	// 1:电话报警
	// 2:为设备报警
	// 3:短信报警
	// 4:GPS报警
	// 5:视频报警
	// 6:设备故障报警
	// 7:他报警
	AlarmMethod string `json:"alarmMethod" binding:"omitempty,oneof=0 1 2 3 4 5 6 7"`
	// 结果
	result string
	// 追踪标识
	TraceID string
}

func (m *Alarm) SetResult(s string) {
	m.result = s
}

// SendAlarm 订阅报警
func SendAlarm(ctx context.Context, m *Alarm) (string, error) {
	// 消息
	var body xml.Subscribe
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdAlarm
	body.DeviceID = m.Device.GetToID()
	body.SN = sip.GetSNString()
	body.StartAlarmPriority = m.StartAlarmPriority
	body.EndAlarmPriority = m.EndAlarmPriority
	body.AlarmMethod = m.AlarmMethod
	//
	var result request.XMLResult
	return result.Result, request.SendSubscribe(ctx, m.TraceID, m.Ser, m.Device, &body, m.Expire, &result)
}
