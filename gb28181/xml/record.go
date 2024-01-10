package xml

// Record 用于录像查询
type Record struct {
	// 设备/区域/系统编码
	DeviceID string `xml:",omitempty" json:"deviceID,omitempty"`
	// 设备/区域/系统名称
	Name string `xml:",omitempty" json:"name,omitempty"`
	// 文件路径名
	FilePath string `xml:",omitempty" json:"filePath,omitempty"`
	// 录像地址
	Address string `xml:",omitempty" json:"address,omitempty"`
	// 录像开始时间
	StartTime string `xml:",omitempty" json:"startTime,omitempty"`
	// 录像结束时间
	EndTime string `xml:",omitempty" json:"endTime,omitempty"`
	// 保密属性
	// 0:不涉密
	// 1:涉密
	Secrecy string `xml:",omitempty" json:"secrecy,omitempty"`
	// 录像产生类型
	// time/alarm/manual
	Type string `xml:",omitempty" json:"type,omitempty"`
	// 录像触发者
	RecorderID string `xml:",omitempty" json:"recorderID,omitempty"`
	//录像文件大小，单位: Byte
	FileSize int64 `xml:",omitempty" json:"fileSize,omitempty"`
}
