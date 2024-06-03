package xml

import (
	"encoding/xml"
)

// Message 包含了所有的字段，用于 MESSAGE 消息
type Message struct {
	// 必须的
	XMLName  xml.Name
	CmdType  string
	SN       string
	DeviceID string
	// Control-DeviceControl-PTZCmd
	// 球机/云台控制命令
	PTZCmd string `xml:",omitempty"`
	// Control-DeviceControl-TeleBoot
	// 远程启动控制命令
	// 只有一个命令 Boot
	TeleBoot string `xml:",omitempty"`
	// Control-DeviceControl-IFameCmd
	// 强制关键帧命令,设备收到此命令应立刻发送一个IDR 帧
	// 只有一个命令 Send
	IFameCmd string `xml:",omitempty"`
	// Control-DeviceControl-DragZoomOut
	// 拉框放大控制命令
	DragZoomOut *MessageDragZoom `xml:",omitempty"`
	// Control-DeviceControl-DragZoomIn
	// 拉框缩小控制命令
	DragZoomIn *MessageDragZoom `xml:",omitempty"`
	// Control-DeviceControl-GuardCmd
	// 报警布防/撤防命令
	// SetGuard / ResetGuard
	GuardCmd string `xml:",omitempty"`
	// Control-DeviceControl-AlarmCmd
	// 报警复位命令
	// 只有一个命令 ResetAlarm
	AlarmCmd string `xml:",omitempty"`
	// Control-DeviceControl-RecordCmd
	// 录像控制命令
	// Record / StopRecord
	RecordCmd string `xml:",omitempty"`
	// Control-DeviceControl-HomePosition
	// 看守位控制命令
	HomePosition *MessageHomePosition `xml:",omitempty"`
	// Control-DeviceConfig
	// 基本参数配置
	BasicParam *MessageBasicParam `xml:",omitempty"`
	// Control-DeviceConfig
	// Response-ConfigDownload
	// SVAC 编码配置
	SVACEncodeConfig *MessageSVACEncodeConfig `xml:",omitempty"`
	// Control-DeviceConfig
	// Response-ConfigDownload
	// SVAC 解码配置
	SVACDecodeConfig *MessageSVACDecodeConfig `xml:",omitempty"`
	// Response-ConfigDownload
	// 视频参数范围，各可选参数以 '/' 分隔
	VideoParamOpt *MessageVideoParamOpt `xml:",omitempty"`
	// Query-Catalog
	// Query-RecordInfo
	// 录像起始时间
	StartTime string `xml:",omitempty"`
	// Query-Catalog
	// Query-RecordInfo
	// 录像终止时间
	EndTime string `xml:",omitempty"`
	// Query-RecordInfo
	// 文件路径名
	FilePath string `xml:",omitempty"`
	// Query-RecordInfo
	// 录像地址支持不完全查询
	Address string `xml:",omitempty"`
	// Query-RecordInfo
	// 保密属性
	// 0: 不涉密
	// 1: 涉密
	Secrecy string `xml:",omitempty"`
	// Query-RecordInfo
	// 录像产生类型，time/alarm/manual/all
	Type string `xml:",omitempty"`
	// Query-RecordInfo
	// 录像触发者
	RecorderID string `xml:",omitempty"`
	// Query-RecordInfo
	// 录像模糊查询属性
	// 0: 不进行模糊查询，此时根据 SIP 消息中 To 头域 URI 中的 ID 值确定查询录像位置，
	// 若 ID 值为本域系统 ID 则进行中心历史记录检索，若为前端设备 ID 则进行前端设备历史记录检索
	// 1: 进行模糊查询，此时设备所在域应同时进行中心检索和前端检索并将结果统一返回
	IndistinctQuery string `xml:",omitempty"`
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
	// 报警方式条件
	// 0:为全部
	// 1:电话报警
	// 2:为设备报警
	// 3:短信报警
	// 4:GPS报警
	// 5:视频报警
	// 6:设备故障报警
	// 7:他报警
	// 可以为直接组合，如 12 为电话报警或设备报警
	//
	// Notify-Alarm
	// 报警方式
	// 1为电话报警
	// 2为设备报警
	// 3为短信报警
	// 4为GPS报警
	// 5为视频报警
	// 6为设备故障报警
	// 7其他报警
	// 注:设备发送报警方式为2的“设备报警”通知后,平台需进行 A.2.3a)“报警复位”控制操作,设备才能发送新的“设备报警”通知。
	AlarmMethod string `xml:",omitempty"`
	// Query-Alarm
	// 报警类型
	AlarmType string `xml:",omitempty"`
	// Query-Alarm
	// 报警发生起止时间
	StartAlarmTime string `xml:",omitempty"`
	// Query-Alarm
	// 报警发生起止时间
	EndAlarmTime string `xml:",omitempty"`
	// Query-ConfigDownload
	// 查询配置参数类型
	// 可查询的配置类型包括基本参数配置: BasicParam
	// 视频参数范围: VideoParamOpt
	// SVAC 编码配置: SVACEncodeConfig
	// SVAC 解码配置: SVACDecodeConfig
	// 可同时查询多个配置类型，各类型以 '/' 分隔，
	// 可返回与查询 SN 值相同的多个响应，
	// 每个响应对应一个配置类型。
	ConfigType string `xml:",omitempty"`
	// Query-MobilePosition
	// 移动设备位置信息上报时间间隔，单位:秒
	Interval int64 `xml:",omitempty"`
	// Notify-Alarm
	// 报警级别
	// 1为一级警情
	// 2为二级警情
	// 3为三级警情
	// 4为四级警情
	AlarmPriority string `xml:",omitempty"`
	// Notify-Alarm
	// 报警时间，国标时间格式
	AlarmTime string `xml:",omitempty"`
	// Notify-Alarm
	// 报警内容描述
	AlarmDescription string `xml:",omitempty"`
	// Notify-Alarm
	// Notify-MobilePosition
	// 经度
	Longitude string `xml:",omitempty"`
	// Notify-Alarm
	// Notify-MobilePosition
	// 纬度
	Latitude string `xml:",omitempty"`
	// Notify-MediaStatus
	// 通知事件类型，取值 "121" 表示历史媒体文件发送结束
	NotifyType string `xml:",omitempty"`
	// Notify-Broadcast
	// 语音输入设备的设备编码
	SourceID string `xml:",omitempty"`
	// Notify-Broadcast
	// 语音输出设备的设备编码
	TargetID string `xml:",omitempty"`
	// Notify-MobilePosition
	// 产生通知时间，国标格式
	Time string `xml:",omitempty"`
	// Notify-MobilePosition
	// 速度，单位: km/h
	Speed string `xml:",omitempty"`
	// Notify-MobilePosition
	// 方向，取值为当前摄像头方向与正北方的顺时针夹角，取值范围0°~360°
	Direction string `xml:",omitempty"`
	// Notify-MobilePosition
	//海拔高度，单位: m
	Altitude string `xml:",omitempty"`
	// Response-DeviceControl
	// Response-Alarm
	// Response-Catalog
	// Response-DeviceInfo
	// Response-DeviceConfig
	// Response-ConfigDownload
	// Response-Broadcast
	// 执行结果标志
	Result string `xml:",omitempty"`
	// Response-Catalog
	// Response-RecordInfo
	// 查询结果总数
	SumNum int64 `xml:",omitempty"`
	// Response-Catalog
	// 设备目录项列表
	DeviceList *MessageDeviceList `xml:",omitempty"`
	// Response-DeviceInfo
	// 目标设备/区域/系统的名称
	DeviceName string `xml:",omitempty"`
	// Response-DeviceInfo
	// 设备生产商
	Manufacturer string `xml:",omitempty"`
	// Response-DeviceInfo
	// 设备型号
	Model string `xml:",omitempty"`
	// Response-DeviceInfo
	// 设备固件版本
	Firmware string `xml:",omitempty"`
	// Response-DeviceInfo
	// 视频输入通道数
	Channel int64 `xml:",omitempty"`
	// Response-DeviceStatus
	// 是否在线，ONLINE/OFFLINE
	Online string `xml:",omitempty"`
	// Response-DeviceStatus
	// 是否正常工作，ON/OFF
	Status string `xml:",omitempty"`
	// Response-DeviceStatus
	// 不正常工作原因
	Reason string `xml:",omitempty"`
	// Response-DeviceStatus
	// 是否编码，ON/OFF
	Encode string `xml:",omitempty"`
	// Response-DeviceStatus
	// 是否录像，ON/OFF
	Record string `xml:",omitempty"`
	// Response-DeviceStatus
	// 设备时间和日期
	DeviceTime string `xml:",omitempty"`
	// Response-DeviceStatus
	// 报警设备状态列表，num 表示列表项个数
	Alarmstatus *MessageAlarmstatus `xml:",omitempty"`
	// Response-RecordInfo
	// 设备/区域名称
	Name string `xml:",omitempty"`
	// Response-RecordInfo
	// 文件目录项列表,Num表示目录项个数
	RecordList *MessageRecordList `xml:",omitempty"`
	// Response-PresetQuery
	PresetList *MessagePresetList `xml:",omitempty"`
	// Control-DeviceConfig
	// Notify-Alarm
	Info *MessageInfo `xml:",omitempty"`
}

// MessagePresetList 是 Message 的 PresetList 字段
type MessagePresetList struct {
	Num  int64                    `xml:"Num,attr"`
	Item []*MessagePresetListItem `xml:",omitempty" json:"item"`
}

// MessagePresetListItem 是 MessagePresetList 的 Item 字段
type MessagePresetListItem struct {
	// 预置位编码
	PresetID string `xml:",omitempty" json:"presetID"`
	// 预置位名称
	PresetName string `xml:",omitempty" json:"presetName"`
}

// MessageVideoParamOpt 是 Message 的 VideoParamOpt 字段
type MessageVideoParamOpt struct {
	//下载倍速范围,各可选参数以“/”分隔,如设备支持1,2,4倍速下载则应写为“1/2/4”
	DownloadSped string `xml:",omitempty"`
	//摄像机支持的分辨率,可有多个分辨率值,各个取值间以“/”分隔。分辨率取值参见附录F中SDPf字段规定
	Resolution string `xml:",omitempty"`
}

// MessageRecordList 是 Message 的 RecordList 字段
type MessageRecordList struct {
	Num  int64     `xml:"Num,attr"`
	Item []*Record `xml:",omitempty"`
}

// MessageAlarmstatus 是 Message 的 Alarmstatus 字段
type MessageAlarmstatus struct {
	Num  int64                     `xml:"Num,attr"`
	Item []*MessageAlarmstatusItem `xml:",omitempty"`
}

// MessageAlarmstatusItem 是 MessageAlarmstatus 的 Item 字段
type MessageAlarmstatusItem struct {
	// 报警设备编码
	DeviceID string `xml:",omitempty"`
	// 报警设备状态
	// ONDUTY / OFFDUTY / ALARM
	DutyStatus string `xml:",omitempty"`
}

// MessageDeviceList 是 Message 的 DeviceList 字段
type MessageDeviceList struct {
	Num  int64     `xml:"Num,attr"`
	Item []*Device `xml:",omitempty"`
}

// MessageInfo 是 Message 的 Info 字段
type MessageInfo struct {
	// Control-DeviceControl
	// 文档上没有说明，就是消息范例那里有出现，我觉得应该是整数，越大越优先
	ControlPriority string `xml:",omitempty"`
	// Notify-Alarm
	// 报警类型
	// 报警方式为2时,不携带AlarmType为默认的报警设备报警,携带AlarmType取值及对应报警类型如下:
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
	//
	// Control-DeviceControl-AlarmCmd
	// 复位报警的报警方式属性
	AlarmType string `xml:",omitempty"`
	// Notify-Alarm
	// 报警类型扩展参数，在入侵检测报警时可携带
	AlarmTypeParam *MessageInfoAlarmTypeParam `xml:",omitempty"`
	// Control-DeviceControl-AlarmCmd
	// 复位报警的报警类型属性
	// 1: 电话报警
	// 2: 设备报警
	// 3: 短信报警
	// 4: GPS 报警
	// 5: 视频报警
	// 6: 设备故障报警
	// 7: 其他报警
	AlarmMethod string `xml:",omitempty"`
	// Notify-Keepalive
	DeviceID []string `xml:",omitempty"`
}

// MessageInfoAlarmTypeParam 是 XMLInfo 的 AlarmTypeParam 字段
type MessageInfoAlarmTypeParam struct {
	// 事件类型
	// 1: 进入区域
	// 2: 离开区域
	EventType string `xml:",omitempty"`
}

// MessageDragZoom 是 Message 的 DragZoomOut/DragZoomIn 字段
type MessageDragZoom struct {
	// 播放窗口长度像素值
	Length int64 `xml:",omitempty"`
	// 播放窗口宽度像素值
	Width int64 `xml:",omitempty"`
	// 拉框中心的横轴坐标像素值
	MidPointX int64 `xml:",omitempty"`
	// 拉框中心的纵轴坐标像素值
	MidPointY int64 `xml:",omitempty"`
	// 拉框长度像素值
	LengthX int64 `xml:",omitempty"`
	// 拉框宽度像素值
	LengthY int64 `xml:",omitempty"`
}

// MessageHomePosition 是 Message 的 HomePosition 字段
type MessageHomePosition struct {
	// 看守位使能
	// 1: 开启
	// 0: 关闭
	Enabled string `xml:",omitempty" json:"enabled,omitempty"`
	// 自动归位时间间隔，开启看守位时使用，单位:秒
	ResetTime int64 `xml:",omitempty" json:"resetTime,omitempty"`
	// 调用预置位编号，开启看守位时使用，取值范围0~255
	PresetIndex int64 `xml:",omitempty" json:"presetIndex"`
}

// MessageBasicParam 是 Message 的 BasicParam 字段
type MessageBasicParam struct {
	// 设备名称
	Name string `xml:",omitempty"`
	// 注册过期时间
	Expiration int64 `xml:",omitempty"`
	// 心跳间隔时间
	HeartBeatInterval int64 `xml:",omitempty"`
	// 心跳超时次数
	HeartBeatCount int64 `xml:",omitempty"`
	// Response-ConfigDownload
	// 定位功能支持情况
	// 0: 不支持
	// 1: 支持 GPS定位
	// 2: 支持北斗定位
	PositionCapability string `xml:",omitempty"`
	// Response-ConfigDownload
	// 经度
	Longitude string `xml:",omitempty"`
	// Response-ConfigDownload
	// 纬度
	Latitude string `xml:",omitempty"`
}

// MessageSVACEncodeConfig 是 XML 的 SVACEncodeConfig 字段
type MessageSVACEncodeConfig struct {
	// 感兴趣区域参数
	ROIParam *MessageSVACEncodeConfigROIParam `xml:",omitempty"`
	// 音频参数
	AudioParam *MessageSVACEncodeConfigAudioParam `xml:",omitempty"`
}

// MessageSVACEncodeConfigROIParam 是 MessageSVACEncodeConfig 的 ROIParam 字段
type MessageSVACEncodeConfigROIParam struct {
	// 感兴趣区域开关
	// 0: 关闭
	// 1: 打开
	ROIFlag string `xml:",omitempty"`
	// 感兴趣区域数量，取值范围 0~16
	ROINumber string `xml:",omitempty"`
	// 感兴趣区域
	Item []*MessageSVACEncodeConfigROIParamItem `xml:",omitempty"`
	// 背景区域编码质量等级可选)
	// 0: 一般
	// 1: 较好
	// 2: 好
	// 3: 很好
	BackGroundQP string `xml:",omitempty"`
	// 背景跳过开关
	// 0: 关闭
	// 1: 打开
	BackGroundSkipFlag string `xml:",omitempty"`
}

// MessageSVACEncodeConfigROIParamItem 是 MessageSVACEncodeConfigROIParam 的 Item 字段
type MessageSVACEncodeConfigROIParamItem struct {
	// 感兴趣区域编号，取值范围 1~16
	ROISeq int8 `xml:",omitempty"`
	// 感兴趣区域左上角坐标，取值范围 0~19683
	TopLeft int64 `xml:",omitempty"`
	// 感兴趣区域右下角坐标，取值范围 0~19683
	BottomRight int64 `xml:",omitempty"`
	// ROI区域编码质量等级
	// 0: 一般
	// 1: 较好
	// 2: 好
	// 3: 很好
	ROIQP string `xml:",omitempty"`
}

// MessageSVACEncodeConfigAudioParam 是 MessageSVACEncodeConfig 的 AudioParam 字段
type MessageSVACEncodeConfigAudioParam struct {
	// 声音识别特征参数开关
	// 0: 关闭
	// 1: 打开
	AudioRecognitionFlag string `xml:",omitempty"`
}

// MessageSVACDecodeConfig 是 XML 的 SVACDecodeConfig 字段
type MessageSVACDecodeConfig struct {
	// SVC 参数
	SVCParam *MessageSVACDecodeConfigSVCParam `xml:",omitempty"`
	// 监控专用信息参数
	SurveilanceParam *MessageSVACDecodeConfigSurveilanceParam `xml:",omitempty"`
}

// MessageSVACDecodeConfigSVCParam 是 MessageSVACDecodeConfig 的 SVCParam 字段
type MessageSVACDecodeConfigSVCParam struct {
	// 感兴趣区域参数
	// 码流显示模式
	// 0: 基本层码流单独显示方式
	// 1: 基本层 +1 个增强层码流方式
	// 2: 基本层 +2 个增强层码流方式
	// 3: 基本层 +3 个增强层码流方式
	SVCSTMMode string `xml:",omitempty"`
	// Response-ConfigDownload
	// 空域编码能力
	// 0: 不支持
	// 1: 1级增强(1个增强层)
	// 2: 2级增强 (2个增强层)
	// 3: 3级增强(3个增强层)
	SVCSpaceSupportMode string `xml:",omitempty"`
	// Response-ConfigDownload
	// 时域编码能力
	// 0: 不支持
	// 1: 1级增强
	// 2: 2级增强
	// 3: 3级增强
	SVCTimeSupportMode string `xml:",omitempty"`
}

// MessageSVACDecodeConfigSurveilanceParam 是 MessageSVACDecodeConfig 的 SurveilanceParam 字段
type MessageSVACDecodeConfigSurveilanceParam struct {
	// 绝对时间信息显示开关
	// 0: 关闭
	// 1: 打开
	TimeShowFlag string `xml:",omitempty"`
	// 监控事件信息显示开关
	// 0: 关闭
	// 1: 打开
	EventShowFlag string `xml:",omitempty"`
	// 报警信息显示开关
	// 0:关闭
	// 1:打开
	AlerShowtFlag string `xml:",omitempty"`
}
