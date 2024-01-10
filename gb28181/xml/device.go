package xml

// Device 表示设备信息
// 用于目录查询/通知
type Device struct {
	// 设备/区域/系统编码(必选)
	DeviceID string `xml:",omitempty"`
	// 设备/区域/系统名称(必选)
	Name string `xml:",omitempty"`
	// 设备厂商(当为设备时必选)
	Manufacturer string `xml:",omitempty"`
	// 设备型号(当为设备时必选)
	Model string `xml:",omitempty"`
	// 设备归属(当为设备时必选)
	Owner string `xml:",omitempty"`
	// 设备固件(当为设备时必选)
	Firmware string `xml:",omitempty"`
	// 行政区域(必选)
	CivilCode string `xml:",omitempty"`
	// 警区(可选)
	Block string `xml:",omitempty"`
	// 安装地址(当为设备时必选)
	Address string `xml:",omitempty"`
	// 是否有子设备(必选)
	// 1: 有
	// 0: 无
	Parental string `xml:",omitempty"`
	// 父设备/区域/系统ID(必选)
	ParentID string `xml:",omitempty"`
	// 信令安全模式(可选)缺省为 0
	// 0: 不采用
	// 2: S/MIME 签名方式
	// 3: S/MIME 加密签名同时采用方式
	// 4: 数字摘要方式
	SafetyWay string `xml:",omitempty"`
	// 注册方式(必选)缺省为 1
	// 1: 符合IETF RFC 3261标准的认证注册模式
	// 2: 基于口令的双向认证注册模式
	// 3: 基于数字证书的双向认证注册模式
	RegisterWay string `xml:",omitempty"`
	// 证书序列号(有证书的设备必选)
	CertNum string `xml:",omitempty"`
	// 证书有效标识(有证书的设备必选)缺省为 0
	// 0: 无效
	// 1: 有效
	Certifiable string `xml:",omitempty"`
	// 证书无效原因码(有证书且证书无效的设备必选)
	ErrCode string `xml:",omitempty"`
	// 证书终止有效期(有证书的设备必选)
	EndTime string `xml:",omitempty"`
	// 保密属性(必选)缺省为 0
	// 0: 不涉密
	// 1: 涉密
	Secrecy string `xml:",omitempty"`
	// 设备/区域/系统 IP地址(可选)
	IPAddress string `xml:",omitempty"`
	// 设备/区域/系统端口(可选)
	Port string `xml:",omitempty"`
	// 设备口令(可选)
	Password string `xml:",omitempty"`
	// 设备状态(必选)，ON/OFF
	Status string `xml:",omitempty"`
	// 经度(可选)
	Longitude string `xml:",omitempty"`
	// 纬度(可选)
	Latitude string `xml:",omitempty"`
	// Info
	Info *DeviceInfo
}

// DeviceInfo 是 Device 的 Info 字段
type DeviceInfo struct {
	// 摄像机类型扩展，标识摄像机类型
	// 1: 球机
	// 2: 半球
	// 3: 固定枪机
	// 4: 遥控枪机
	PTZType string `xml:",omitempty"`
	// 摄像机位置类型扩展
	// 1: 省际检查站
	// 2: 党政机关
	// 3: 车站码头
	// 4: 中心广场
	// 5: 体育场馆
	// 6: 商业中心
	// 7: 宗教场所
	// 8: 校园周边
	// 9: 治安复杂区域
	// 10: 交通干线
	PositionType string `xml:",omitempty"`
	// 摄像机安装位置室外、室内属性
	// 1: 室外
	// 2: 室内
	RoomType string `xml:",omitempty"`
	// 摄像机用途属性
	// 1: 治安
	// 2: 交通
	// 3: 重点
	UseType string `xml:",omitempty"`
	// 摄像机补光属性
	// 1: 无补光
	// 2: 红外补光
	// 3: 白光补光
	SupplyLightType string `xml:",omitempty"`
	// 摄像机方位属性
	// 1: 东
	// 2: 南
	// 3: 西
	// 4: 北
	// 5: 东南
	// 6: 东北
	// 7: 西南
	// 8: 西北
	DirectionType string `xml:",omitempty"`
	// 摄像机支持的分辨率
	Resolution string `xml:",omitempty"`
	// 虚拟组织所属的业务分组 ID
	BusinessGroupID string `xml:",omitempty"`
	// 下载倍速范围，各可选参数以"/"分隔，如设备支持 1、2、4 倍速下载则应写为 "1/2/4"
	DownloadSpeed string `xml:",omitempty"`
	// 空域编码能力
	// 0: 不支持
	// 1: 1 级增强
	// 2: 2 级增强
	// 3: 3 级增强
	SVCSpaceSupportMode string `xml:",omitempty"`
	// 时域编码能力
	// 0: 不支持
	// 1: 1 级增强
	// 2: 2 级增强
	// 3: 3 级增强
	SVCTimeSupportMode string `xml:",omitempty"`
}

// FilterDir 过滤掉目录，只要通道
// 修改 ParentID 和 Parental
func FilterDir(parentID string, ms []*Device) []*Device {
	// 过滤
	data := make(map[string]*Device)
	for _, m := range ms {
		// 没有编号？
		if m.DeviceID == "" {
			continue
		}
		data[m.DeviceID] = m
		// 如果 d.ParentID 存在 data 中，说明是目录，过滤掉
		if m.DeviceID != m.ParentID {
			delete(data, m.ParentID)
		}
		m.Parental = "0"
		m.ParentID = parentID
	}
	// 返回
	ms = ms[:0]
	for _, d := range data {
		ms = append(ms, d)
	}
	return ms
}
