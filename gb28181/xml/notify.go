package xml

import (
	"encoding/xml"
)

// NotifyDeviceListItem.Event
const (
	NotifyEventON     = "ON"
	NotifyEventOFF    = "OFF"
	NotifyEventVLOST  = "VLOST"
	NotifyEventDEFECT = "DEFECT"
	NotifyEventADD    = "ADD"
	NotifyEventDEL    = "DEL"
	NotifyEventUPDATE = "UPDATE"
)

// Notify 包含了所有的字段，用于 NOTIFY 消息
type Notify struct {
	// 必须的
	XMLName  xml.Name
	CmdType  string
	SN       string
	DeviceID string
	// Response-Alarm
	// 执行结果标志
	Result string `xml:",omitempty"`
	// Notify-Catalog
	// 查询结果总数
	SumNum int64 `xml:",omitempty"`
	// Notify-Catalog
	// 设备目录项列表
	DeviceList *NotifyDeviceList `xml:",omitempty"`
	// Query-Alarm
	// 报警起始级别
	// 0:为全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	StartAlarmPriority string `xml:",omitempty"`
	// Query-Alarm
	// 报警终止级别
	// 0:为全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	EndAlarmPriority string `xml:",omitempty"`
	// Query-Alarm
	// Notify-Alarm
	// 报警方式
	// 1为电话报警
	// 2为设备报警
	// 3为短信报警
	// 4为GPS报警
	// 5为视频报警
	// 6为设备故障报警
	// 7其他报警
	AlarmMethod string `xml:",omitempty"`
	// Query-Alarm
	// 报警发生起止时间
	StartTime string `xml:",omitempty"`
	// Query-Alarm
	// 报警发生起止时间
	EndTime string `xml:",omitempty"`
	// Notify-Alarm
	// 报警终止级别
	// 0:为全部
	// 1:一级警情
	// 2:二级警情
	// 3:三级警情
	// 4:四级警情
	AlarmPriority string `xml:",omitempty"`
	// Notify-Alarm
	// 报警时间
	AlarmTime string `xml:",omitempty"`
	// Notify-Alarm
	// 报警内容描述
	AlarmDescription string `xml:",omitempty"`
	// Notify-Alarm
	// Notify-MobilePosition
	// 经纬度信息
	Longitude string `xml:",omitempty"`
	// Notify-Alarm
	// Notify-MobilePosition
	// 经纬度信息
	Latitude string `xml:",omitempty"`
	// Notify-MobilePosition
	// 产生通知时间
	Time string `xml:",omitempty"`
	// Notify-MobilePosition
	// 速度，单位  km/h
	Speed string `xml:",omitempty"`
	// 方向, 取值为当前摄像头方向与正北方的顺时针夹角，取值范围 0°~360°
	Direction string `xml:",omitempty"`
	// 海拔，单位 m
	Altitude string `xml:",omitempty"`
}

// NotifyDeviceList 是 Notify 的 DeviceList 字段
type NotifyDeviceList struct {
	Num  int64 `xml:"Num,attr"`
	Item []*NotifyDeviceListItem
}

// NotifyDeviceListItem 是 NotifyDeviceList 的 Item 字段
type NotifyDeviceListItem struct {
	Device
	// 通知事件
	Event string
}
