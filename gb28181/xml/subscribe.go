package xml

import "encoding/xml"

// Subscribe 包含了所有的字段，用于 SUBSCRIBE 消息
type Subscribe struct {
	// 基本
	XMLName  xml.Name
	CmdType  string
	SN       string
	DeviceID string
	// Query-Alarm
	// 报警起始级别(可选)
	// 0:全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	StartAlarmPriority string `xml:",omitempty"`
	// Query-Alarm
	// 报警终止级别(可选)
	// 0:全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	EndAlarmPriority string `xml:",omitempty"`
	// Query-Alarm
	// 报警方式(必选)
	// 0:全部
	// 1为电话报警
	// 2为设备报警
	// 3为短信报警
	// 4为GPS报警
	// 5为视频报警
	// 6为设备故障报警
	// 7其他报警
	AlarmMethod string `xml:",omitempty"`
	// Query-Alarm
	// 起止时间(可选)
	// Query-Catalog
	// 起止时间(可选)
	StartTime string `xml:",omitempty"`
	// Query-Alarm
	// 起止时间(可选)
	// Query-Catalog
	// 起止时间(可选)
	EndTime string `xml:",omitempty"`
	// Query-Position
	// 移动设备位置信息上报时间间隔，单位:秒，默认 5
	Interval int64 `xml:",omitempty"`
}
