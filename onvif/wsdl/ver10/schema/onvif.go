package schema

import "time"

// Time 时间
type Time struct {
	Hour   int
	Minute int
	Second int
}

// Date 日期
type Date struct {
	Year  int
	Month int
	Day   int
}

// DateTime 日期和时间
type DateTime struct {
	Time Time
	Date Date
}

// ToTime 转换返回 go time
func (t *DateTime) ToTime(location *time.Location) time.Time {
	return time.Date(t.Date.Year,
		time.Month(t.Date.Month),
		t.Date.Day,
		t.Time.Hour,
		t.Time.Minute,
		t.Time.Second,
		0, location)
}

// TimeZone 时区
type TimeZone struct {
	TZ string
}

// DateTimeType 类型
type DateTimeType string

// DateTimeType 枚举
const (
	DateTimeTypeManual DateTimeType = "Manual"
	DateTimeTypeNTP    DateTimeType = "NTP"
)

// SystemDateTime 系统日期时间
type SystemDateTime struct {
	// 类型，Manual / NTP
	DateTimeType DateTimeType
	// 夏令时当前是否打开
	DaylightSavings bool
	// Posix 格式
	TimeZone TimeZone
	// utc 时间
	UTCDateTime DateTime
	// local 时间
	LocalDateTime DateTime
	// 没有说明
	Extension string `xml:",chardata"`
}

// UTC 返回 utc 的 go time
func (t *SystemDateTime) UTC() time.Time {
	return t.UTCDateTime.ToTime(time.UTC)
}

// Local 返回 local 的 go time
func (t *SystemDateTime) Local() time.Time {
	return t.LocalDateTime.ToTime(time.Local)
}

// Capabilities 能力列表
type Capabilities struct {
	// 分析能力，nil 表示不提供
	Analytics *AnalyticsCapabilities
	// 设备能力，nil 表示不提供
	Device *DeviceCapabilities
	// 事件能力，nil 表示不提供
	Events *EventCapabilities
	// 图像能力，nil 表示不提供
	Imaging *ImagingCapabilities
	// 媒体能力，nil 表示不提供
	Media *MediaCapabilities
	// 云台控制能力，nil 表示不提供
	PTZ *PTZCapabilities
	// 没有说明
	Extension *CapabilitiesExtension
}

// AnalyticsCapabilities 分析能力
type AnalyticsCapabilities struct {
	// 服务地址
	XAddr string
	// 规则是否支持
	RuleSupport bool
	// 模块是否支持
	AnalyticsModuleSupport bool
}

// DeviceCapabilities 设备能力
type DeviceCapabilities struct {
	// 服务地址
	XAddr string
	// 网络
	Network *NetworkCapabilities
	// 系统
	System *SystemCapabilities
	// IO
	IO *IOCapabilities
	// 安全
	Security *SecurityCapabilities
	// 没有说明
	Extension string `xml:",chardata"`
}

// NetworkCapabilities 网络能力
type NetworkCapabilities struct {
	// ip 过滤是否支持
	IPFilter bool
	// 零配置是否支持
	ZeroConfiguration bool
	// ipv6
	IPVersion6 bool
	// 动态 dns
	DynDNS bool
	// 没有说明
	Extension NetworkCapabilitiesExtension
}

// NetworkCapabilitiesExtension 网络能力扩展
type NetworkCapabilitiesExtension struct {
	// 没有说明
	Dot11Configuration bool
	// 没有说明
	Extension string `xml:",chardata"`
}

// SystemCapabilities 系统能力
type SystemCapabilities struct {
	// ws-discovery
	DiscoveryResolve bool
	// ws-discovery bye
	DiscoveryBye bool
	// remote discovery
	RemoteDiscovery bool
	// 系统备份是否支持
	SystemBackup bool
	// 系统日志是否支持
	SystemLogging bool
	// 防火墙升级是否支持
	FirmwareUpgrade bool
	// 支持的 onvif 版本
	SupportedVersions OnvifVersion
	// 扩展
	Extension SystemCapabilitiesExtension
}

// OnvifVersion 版本
type OnvifVersion struct {
	// 高版本
	Major int
	// 低版本，两个数字
	// 如果主版本号小于 16
	// X.0.1 映射到 01
	// X.2.1 映射到 21
	// 其中X代表主版本号
	// 否则发布月份，例如 06 表示六月。
	Minor int
}

// SystemCapabilitiesExtension 系统能力扩展
type SystemCapabilitiesExtension struct {
	// 没有说明
	HTTPFirmwareUpgrade bool `xml:"HttpFirmwareUpgrade"`
	// 没有说明
	HTTPSystemBackup bool `xml:"HttpSystemBackup"`
	// 没有说明
	HTTPSystemLogging bool `xml:"HttpSystemLogging"`
	// 没有说明
	HTTPSupportInformation bool `xml:"HttpSupportInformation"`
	// 没有说明
	Extension string `xml:",chardata"`
}

// IOCapabilities IO 能力
type IOCapabilities struct {
	// 输入的连接器数量
	InputConnectors int
	// 输出的转发数量
	RelayOutputs int
	// 没有说明
	Extension IOCapabilitiesExtension
}

// IOCapabilitiesExtension IO 能力扩展
type IOCapabilitiesExtension struct {
	// 没有说明
	Auxiliary bool
	// 没有说明
	AuxiliaryCommands string
	// 没有说明
	Extension string `xml:",chardata"`
}

// SecurityCapabilities 安全能力
type SecurityCapabilities struct {
	// tls 1.1
	TLS11 bool `xml:"TLS1.1"`
	// tls 1.2
	TLS12 bool `xml:"TLS1.2"`
	// key 生成是否支持
	OnboardKeyGeneration bool
	// 访问策略配置
	AccessPolicyConfig bool
	// WS-Security X.509 token
	X509Token bool `xml:"X.509Token"`
	// WS-Security SAML token
	SAMLToken bool
	// WS-Security Kerberos token
	KerberosToken bool
	// WS-Security REL token
	RELToken bool
	// 没有说明
	Extension SecurityCapabilitiesExtension
}

// SecurityCapabilitiesExtension 安全能力扩展
type SecurityCapabilitiesExtension struct {
	// 没有说明
	TLS10 bool `xml:"TLS1.0"`
	// 没有说明
	Extension SecurityCapabilitiesExtension2
}

// SecurityCapabilitiesExtension2 安全能力扩展
type SecurityCapabilitiesExtension2 struct {
	// 没有说明
	Dot1X bool
	// EAP 方法
	SupportedEAPMethod int
	// 没有说明
	RemoteUserHandling bool
}

// EventCapabilities 事件能力
type EventCapabilities struct {
	// 服务地址
	XAddr string
	// WS Subscription policy
	WSSubscriptionPolicySupport bool
	// WS Pull Point
	WSPullPointSupport bool
	// WS Pausable Subscription Manager Interface
	WSPausableSubscriptionManagerInterfaceSupport bool
}

// ImagingCapabilities 图像能力
type ImagingCapabilities struct {
	// 服务地址
	XAddr string
}

// MediaCapabilities 媒体能力
type MediaCapabilities struct {
	// 服务地址
	XAddr string
	// 媒体流能力
	StreamingCapabilities *RealTimeStreamingCapabilities
	// 没有说明
	Extension *MediaCapabilitiesExtension
}

// RealTimeStreamingCapabilities 实时流能力
type RealTimeStreamingCapabilities struct {
	// rtp 多播
	RTPMulticast bool
	// rtp over tcp
	RTPOverTCP bool `xml:"RTP_TCP"`
	// rtp/rtsp/tcp
	RTPRTSPTCP bool `xml:"RTP_RTSP_TCP"`
	// 没有说明
	Extension string `xml:",chardata"`
}

// MediaCapabilitiesExtension 是 MediaCapabilities 的 Extension 字段
type MediaCapabilitiesExtension struct {
	ProfileCapabilities ProfileCapabilities
}

// ProfileCapabilities 属性能力
type ProfileCapabilities struct {
	// profile 数量
	MaximumNumberOfProfiles int
}

// PTZCapabilities 云台能力
type PTZCapabilities struct {
	// 服务地址
	XAddr string
}

// CapabilitiesExtension 是 Capabilities 的 Extension 字段
type CapabilitiesExtension struct {
	// 没有说明
	DeviceIO *DeviceIOCapabilities
	// 没有说明
	Display *DisplayCapabilities
	// 没有说明
	Recording *RecordingCapabilities
	// 没有说明
	Search *SearchCapabilities
	// 没有说明
	Replay *ReplayCapabilities
	// 没有说明
	Receiver *ReceiverCapabilities
	// 没有说明
	AnalyticsDevice *AnalyticsDeviceCapabilities
	// 没有说明
	Extensions string `xml:",chardata"`
}

// DeviceIOCapabilities 设备 io 能力
type DeviceIOCapabilities struct {
	// 没有说明
	XAddr string
	// 没有说明
	VideoSources int
	// 没有说明
	VideoOutputs int
	// 没有说明
	AudioSources int
	// 没有说明
	AudioOutputs int
	// 没有说明
	RelayOutputs int
}

// DisplayCapabilities 显示能力
type DisplayCapabilities struct {
	// 没有说明
	XAddr string
	// SetLayout 命令仅支持预定义布局
	FixedLayout bool
}

// RecordingCapabilities 录像能力
type RecordingCapabilities struct {
	// 没有说明
	XAddr string
	// 没有说明
	ReceiverSource bool
	// 没有说明
	MediaProfileSource bool
	// 没有说明
	DynamicRecordings bool
	// 没有说明
	DynamicTracks bool
	// 没有说明
	MaxStringLength int
}

// SearchCapabilities 查询能力？
type SearchCapabilities struct {
	// 没有说明
	XAddr string
	// 没有说明
	MetadataSearch bool
}

// ReplayCapabilities 转发能力
type ReplayCapabilities struct {
	// 服务地址
	XAddr string
}

// ReceiverCapabilities 接收能力
type ReceiverCapabilities struct {
	// 服务地址
	XAddr string
	// 设备是否能接收 rtp 多播流
	RTPMulticast bool `xml:"RTP_Multicast"`
	// 设备是否能接收 rtp/tcp 流
	RTPTCP bool `xml:"RTP_TCP"`
	// 设备是否能接收 rtp/rtsp/tcp 流
	RTPRTSPTCP bool `xml:"RTP_RTSP_TCP"`
	// 设备支持的最大接收数量
	SupportedReceivers int
	// rtsp url 字符串的最大长度
	MaximumRTSPURILength int
}

// AnalyticsDeviceCapabilities 分析能力
type AnalyticsDeviceCapabilities struct {
	// 没有说明
	XAddr string
	// 过时的
	RuleSupport bool
	// 没有说明
	Extension string `xml:",chardata"`
}

// IntRectangle 由左下角位置和大小定义的矩形，单位是像素
type IntRectangle struct {
	X      int `xml:"x,attr"`
	Y      int `xml:"y,attr"`
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

// Profile 表示一个媒体属性的集合
type Profile struct {
	// 可读性名称
	Name string
	// 唯一标识
	Token string `xml:"token,attr"`
	// 是否可以删除
	Fixed bool `xml:"fixed,attr"`
	// 视频输入可选配置
	VideoSourceConfiguration *VideoSourceConfiguration
	// 音频输入可选配置
	AudioSourceConfiguration *AudioSourceConfiguration
	// 视频编码器可选配置
	VideoEncoderConfiguration *VideoEncoderConfiguration
	// 音频编码器可选配置
	AudioEncoderConfiguration *AudioEncoderConfiguration
	// 视频分析模块和规则引擎
	VideoAnalyticsConfiguration *VideoAnalyticsConfiguration
	// 平移-倾斜-缩放单元
	PTZConfiguration *PTZConfiguration
	// 元数据流可选配置
	MetadataConfiguration *MetadataConfiguration
	// 没有说明
	Extension *ProfileExtension
}

// ConfigurationEntity 是公共字段
type ConfigurationEntity struct {
	// 引用该配置的唯一标识
	Token string `xml:"token,attr"`
	// 可读性行名称
	Name string
	// 使用该配置的内部引用数量
	UseCount int
}

// VideoSourceConfiguration 视频源配置
type VideoSourceConfiguration struct {
	ConfigurationEntity
	// 用于支持不同视图模式的设备
	ViewMode string `xml:",attr"`
	// 物理输入的引用？？
	SourceToken string
	// 视频拍摄的区域
	Bounds IntRectangle
	// 没有说明
	Extension *VideoSourceConfigurationExtension
}

// VideoSourceConfigurationExtension 视频源配置扩展
type VideoSourceConfigurationExtension struct {
	// 用于配置捕获图像旋转的可选元素
	// 设备支持的分辨率不受影响
	Rotate *Rotate
	// 没有说明
	Extension *VideoSourceConfigurationExtension2
}

// VideoSourceConfigurationExtension2 视频源配置扩展
type VideoSourceConfigurationExtension2 struct {
	// 描述几何透镜畸变
	LensDescription []*LensDescription
	// 描述摄影机视野中场景方向
	SceneOrientation *SceneOrientation
}

// LensDescription 镜头描述
type LensDescription struct {
	// 光学系统的可选焦距
	FocalLength float64 `xml:",attr"`
	// 在归一化坐标中，透镜中心到成像器中心的偏移
	Offset *LensOffset
	// 投影特性的径向描述
	Projection []*LensProjection
	// ONVIF标准化坐标系所需x坐标的补偿
	XFactor float64
}

// LensOffset 镜头偏移
type LensOffset struct {
	// 标准化坐标中透镜中心的可选水平偏移
	X float64 `xml:"x,attr"`
	// 标准化坐标中透镜中心的可选垂直偏移
	Y float64 `xml:"y,attr"`
}

// LensProjection 镜头投影
type LensProjection struct {
	// 入射角度
	Angle float64
	// 出射角度的映射半径
	Radius float64
	// 在给定角度下，由于渐晕的光线吸收
	// 值为 1 表示没有吸收
	Transmittance float64
}

// SceneOrientationMode 场景方向枚举
type SceneOrientationMode string

// SceneOrientationMode 场景方向枚举
const (
	SceneOrientationModeMANUAL SceneOrientationMode = "MANUAL"
	SceneOrientationModeAUTO   SceneOrientationMode = "AUTO"
)

// SceneOrientation 场景方向
type SceneOrientation struct {
	// 用于指定摄影机确定场景方向的方式的参数
	Mode *SceneOrientationMode
	// 基于 Mode 指定或确定的场景方向
	// 将模式指定为 AUTO 时，此字段是可选的，将被设备忽略
	// 将模式指定为 MANUAL 时，此字段是必需的
	// 如果丢失，设备将返回InvalidArgs错误
	Orientation string
}

// RotateMode 旋转枚举
type RotateMode string

// RotateMode 旋转枚举
const (
	RotateModeON   RotateMode = "ON"
	RotateModeOFF  RotateMode = "OFF"
	RotateModeAUTO RotateMode = "AUTO"
)

// Rotate 旋转
type Rotate struct {
	// 启用/禁用旋转功能的参数
	Mode RotateMode
	// 用于配置 ON 模式下图像顺时针旋转的度数
	// 在 ON 模式中省略此参数意味着旋转180度
	Degree int
	// 没有说明
	Extension string `xml:",chardata"`
}

// AudioSourceConfiguration 音频源配置
type AudioSourceConfiguration struct {
	ConfigurationEntity
	// 应用该配置的 token
	SourceToken string
}

// VideoEncoding 视频编码类型枚举
type VideoEncoding string

// VideoEncoding 视频编码类型枚举
const (
	VideoEncodingJPEG  VideoEncoding = "JPEG"
	VideoEncodingMPEG4 VideoEncoding = "MPEG4"
	VideoEncodingH264  VideoEncoding = "H264"
)

// VideoEncoderConfiguration 视频编码配置
type VideoEncoderConfiguration struct {
	ConfigurationEntity
	// 使用的视频编码
	Encoding VideoEncoding
	// 视频的分辨率
	Resolution *VideoResolution
	// 视频量化器和视频质量的相对值
	// 支持的质量范围内的高值意味着更高的质量
	Quality float64
	// 配置速率控制相关参数
	RateControl *VideoRateControl
	// 用于配置 mpeg4 相关参数
	MPEG4 *Mpeg4Configuration
	// 用于配置 H.264 相关参数
	H264 *H264Configuration
	// 可用于视频流的多播设置
	Multicast *MulticastConfiguration
	// rtsp 视频流的会话超时
	SessionTimeout string
}

// VideoResolution 视频分辨率
type VideoResolution struct {
	// 宽
	Width int
	// 高
	Height int
}

// VideoRateControl 视频输出率控制
type VideoRateControl struct {
	// 最大输出帧速率，fps
	FrameRateLimit int
	// 对图像进行编码和传输的间隔
	EncodingInterval int
	// 最大输出比特率，单位 kbps
	BitrateLimit int
}

// Mpeg4Profile mpeg4 属性枚举
type Mpeg4Profile string

// Mpeg4Profile mpeg4 属性枚举
const (
	Mpeg4ProfileSP  Mpeg4Profile = "SP"
	Mpeg4ProfileASP Mpeg4Profile = "ASP"
)

// Mpeg4Configuration mpeg4 配置
type Mpeg4Configuration struct {
	// 确定 I 帧的编码间隔
	GovLength int
	// 属性文件名称
	Mpeg4Profile Mpeg4Profile
}

// H264Profile h264 属性枚举
type H264Profile string

// H264Profile h264 属性枚举
const (
	H264ProfileBaseline H264Profile = "Baseline"
	H264ProfileMain     H264Profile = "Main"
	H264ProfileExtended H264Profile = "Extended"
	H264ProfileHigh     H264Profile = "High"
)

// H264Configuration h264 配置
type H264Configuration struct {
	// 确定 I 帧的编码间隔
	GovLength int
	// 属性文件名称
	H264Profile H264Profile
}

// MulticastConfiguration 媒体流多播配置
type MulticastConfiguration struct {
	// 多播地址
	Address IPAddress
	// rtp 多播的端口
	Port int
	// ipv6 下的 ttl 次数
	TTL int
	// 流是否持久的
	AutoStart bool
}

// IPType ip 类型枚举
type IPType string

// IPType ip 类型枚举
const (
	IPTypeIPv4 IPType = "IPv4"
	IPTypeIPv6 IPType = "IPv6"
)

// IPAddress 地址
type IPAddress struct {
	// 类型
	Type IPType
	// ipv4 地址
	IPv4Address string
	// ipv6 地址
	IPv6Address string
}

// AudioEncoding 音频编码类型枚举
type AudioEncoding string

// AudioEncoding 音频编码类型枚举
const (
	AudioEncodingG711 AudioEncoding = "G711"
	AudioEncodingG726 AudioEncoding = "G726"
	AudioEncodingACC  AudioEncoding = "ACC"
)

// AudioEncoderConfiguration 音频编码配置
type AudioEncoderConfiguration struct {
	ConfigurationEntity
	// 编码类型
	Encoding AudioEncoding
	// 输出的比特率，单位 kbps
	Bitrate int
	// 输出的采样率，单位 kHz
	SampleRate int
	// 多播设置
	Multicast *MulticastConfiguration
	// rtsp 音频流的会话超时
	SessionTimeout string
}

// VideoAnalyticsConfiguration 视频分析配置
type VideoAnalyticsConfiguration struct {
	ConfigurationEntity
	// 没有说明
	AnalyticsEngineConfiguration *AnalyticsEngineConfiguration
	// 没有说明
	RuleEngineConfiguration *RuleEngineConfiguration
}

// AnalyticsEngineConfiguration 分析引擎配置
type AnalyticsEngineConfiguration struct {
	// 没有说明
	AnalyticsModule *Config
	// 没有说明
	Extension string `xml:",chardata"`
}

// Config 没有说明
type Config struct {
	// 没有说明
	Parameters *ItemList
	// 配置名称
	Name string `xml:",attr"`
	// 指定规则的类型
	// 并且应等于 GetSupportedRules 和 GetSupportedAnalyticsModules 命令返回的
	// ConfigDescription 元素的 Name 属性之一的值
	Type string `xml:",attr"`
}

// ItemList 数组
type ItemList struct {
	// 没有说明
	SimpleItem []*SimpleItem
	// 没有说明
	ElementItem []*ElementItem
	// 没有说明
	Extension string `xml:",chardata"`
}

// SimpleItem name-value 的键值对
type SimpleItem struct {
	Name  string `xml:",attr"`
	Value string `xml:",attr"`
}

// ElementItem name-value 的键值对
type ElementItem struct {
	Name  string `xml:",attr"`
	Value string `xml:",chardata"`
}

// RuleEngineConfiguration 规则引擎
type RuleEngineConfiguration struct {
	// 没有说明
	Rule *Config
	// 没有说明
	Extension string `xml:",chardata"`
}

// PTZConfiguration 云台控制配置
type PTZConfiguration struct {
	ConfigurationEntity
	// 改配置所属的 ptz 节点的强制性引用
	NodeToken string
	// 绝对的水平/倾斜的默认值
	DefaultAbsolutePantTiltPositionSpace string
	// 绝对的缩放的默认值
	DefaultAbsoluteZoomPositionSpace string
	// 相对的水平/倾斜的默认值
	DefaultRelativePanTiltTranslationSpace string
	// 相对的缩放的默认值
	DefaultRelativeZoomTranslationSpace string
	// 连续的水平/倾斜的默认值
	DefaultContinuousPanTiltVelocitySpace string
	// 连续的缩放的默认值
	DefaultContinuousZoomVelocitySpace string
	// 绝对/相对的速度的默认值
	DefaultPTZSpeed *PTZSpeed
	// 连续的操作，超时后停止
	DefaultPTZTimeout string
	// 水平/倾斜限制
	PanTiltLimits *PanTiltLimits
	// 缩放限制
	ZoomLimits *ZoomLimits
	// 没有说明
	Extension *PTZConfigurationExtension
	// 设备移动的加速？？
	MoveRamp int `xml:",attr"`
	// 在调用 presets 的加速
	PresetRamp int `xml:",attr"`
	// 执行 PresetTours 的加速
	PresetTourRamp int `xml:",attr"`
}

// PTZSpeed 速度
type PTZSpeed struct {
	// 水平/倾斜速度
	PanTilt Vector2D
	// 缩放速度
	Zoom Vector1D
}

// Vector2D ptz 的相关类型
type Vector2D struct {
	// 没有说明
	X float64 `xml:"x,attr"`
	// 没有说明
	Y float64 `xml:"y,attr"`
	// 没有说明
	Space string `xml:"space,attr"`
}

// Vector1D ptz 的相关类型
type Vector1D struct {
	// 没有说明
	X float64 `xml:"x,attr"`
	// 没有说明
	Space string `xml:"space,attr"`
}

// PanTiltLimits 水平/倾斜限制
type PanTiltLimits struct {
	Range *Space2DDescription
}

// Space2DDescription 坐标系描述
type Space2DDescription struct {
	// 坐标系地址？
	URI string
	// x 轴范围
	XRange *FloatRange
	// y 轴范围
	YRange *FloatRange
}

// FloatRange 范围
type FloatRange struct {
	// 最小值
	Min float64
	// 最大值
	Max float64
}

// ZoomLimits 缩放限制
type ZoomLimits struct {
	Range *Space1DDescription
}

// Space1DDescription 坐标系描述
type Space1DDescription struct {
	// 坐标系地址？
	URI string
	// x 轴范围
	XRange *FloatRange
}

// PTZConfigurationExtension 云台扩展
type PTZConfigurationExtension struct {
	// 方向控制
	PTControlDirection *PTControlDirection
	// 没有说明
	Extension string `xml:",chardata"`
}

// PTControlDirection 方向控制配置
type PTControlDirection struct {
	// e-flip 陪孩子
	EFlip *EFlip
	// 方向反转
	Reverse *Reverse
	// 没有说明
	Extension string `xml:",chardata"`
}

// EFlipMode E-Flip 特性枚举
type EFlipMode string

// EFlipMode E-Flip 特性枚举
const (
	EFlipModeOFF      EFlipMode = "OFF"
	EFlipModeON       EFlipMode = "ON"
	EFlipModeExtended EFlipMode = "Extended"
)

// EFlip E-Flip 选项
type EFlip struct {
	// 启用/禁用 E-Flip 的特性
	Mode EFlipMode
}

// ReverseMode 控制方向反转枚举
type ReverseMode string

// EFlipMode 控制方向反转枚举
const (
	ReverseModeOFF      ReverseMode = "OFF"
	ReverseModeON       ReverseMode = "ON"
	ReverseModeAUTO     ReverseMode = "AUTO"
	ReverseModeExtended ReverseMode = "Extended"
)

// Reverse ptz 控制方向反转控制
type Reverse struct {
	Mode ReverseMode
}

// MetadataConfiguration 元数据流
type MetadataConfiguration struct {
	ConfigurationEntity
	// 用于配置元数据流中包含的 ptz 数据
	PTZStatus *PTZFilter
	// 用于配置事件流
	Events *EventSubscription
	// 是否包含来自分析引擎的元数据
	Analytics bool
	// 视频流的多播设置
	Multicast *MulticastConfiguration
	// rtsp 会话超时，Media2 忽略
	SessionTimeout string
	// 指示哪些分析模块应输出元数据
	AnalyticsEngineConfiguration *AnalyticsEngineConfiguration
	// 没有说明
	Extension string `xml:",chardata"`
	// 用于配置元数据负载的压缩类型
	CompressionType string `xml:",attr"`
	// 用于配置元数据流是否应包含每个目标的地理位置坐标
	GeoLocation bool `xml:",attr"`
	// 用于配置生成的元数据流是否应包含多边形形状信息
	ShapePolygon bool `xml:",attr"`
}

// PTZFilter 没有说明
type PTZFilter struct {
	// 是否包含 ptz 状态
	Status bool
	// 是否包含 ptz 位置
	Position bool
}

// EventSubscription 订阅
type EventSubscription struct {
	// 没有说明
	Filter string `xml:",chardata"`
	// 没有说明
	SubscriptionPolicy string `xml:",chardata"`
}

// ProfileExtension 属性扩展
type ProfileExtension struct {
	// 音频输出配置
	AudioOutputConfiguration *AudioOutputConfiguration
	// 音频解码器配置
	AudioDecoderConfiguration *AudioDecoderConfiguration
	// 没有说明
	Extension string `xml:",chardata"`
}

// AudioOutputConfiguration 音频输出配置
type AudioOutputConfiguration struct {
	ConfigurationEntity
	// 物理音频输出 token
	OutputToken string
	// 指定音频流的方向，c->s / s->c
	SendPrimacy string
	// 输出的音量
	OutputLevel int
}

// AudioDecoderConfiguration 音频解码器配置
type AudioDecoderConfiguration struct {
	ConfigurationEntity
}

// IntRectangleRange 矩形的范围
type IntRectangleRange struct {
	// x 轴范围
	XRange IntRange
	// y 轴范围
	YRange IntRange
	// 宽度范围
	WidthRange IntRange
	// 高度范围
	HeightRange IntRange
}

// VideoSourceConfigurationOptions 音频源配置
type VideoSourceConfigurationOptions struct {
	// 捕获区域支持的范围
	BoundsRange *IntRectangleRange
	// 物理输入的列表
	VideoSourceTokensAvailable []string
	// 没有说明
	Extension *VideoSourceConfigurationOptionsExtension
	// Profile 最大数量
	MaximumNumberOfProfiles int `xml:",attr"`
}

// VideoSourceConfigurationOptionsExtension 扩展
type VideoSourceConfigurationOptionsExtension struct {
	Rotate    *RotateOptions
	Extension *VideoSourceConfigurationOptionsExtension2
}

// VideoSourceConfigurationOptionsExtension2 扩展
type VideoSourceConfigurationOptionsExtension2 struct {
	SceneOrientationMode []SceneOrientationMode
}

// RotateOptions 旋转选项
type RotateOptions struct {
	Mode []RotateMode
	// 支持的旋转度数值列表
	DegreeList *IntList
	// 没有说明
	Extension string `xml:",chardata"`
	// 如果设备在更改旋转后，是否需要重新启动
	Reboot bool `xml:",attr"`
}

// IntList int 数组
type IntList struct {
	Items []int
}

// IntRange 范围
type IntRange struct {
	// 最小值
	Min int
	// 最大值
	Max int
}

// VideoEncoderConfigurationOptions 视频编码器配置选项
type VideoEncoderConfigurationOptions struct {
	// 质量范围
	QualityRange *IntRange
	// jpeg 选项
	JPEG *JpegOptions
	// mpeg4 选项
	MPEG4 *Mpeg4Options
	// h264 选项
	H264 *H264Options
	// 没有说明
	Extension *VideoEncoderOptionsExtension
	// 是否支持 GuarantineedFrameRate
	// GuarantineedFrameRate 没有找到
	GuaranteedFrameRateSupported bool `xml:",attr"`
}

// JpegOptions jpeg 选项
type JpegOptions struct {
	// 图像大小的支持列表
	ResolutionsAvailable []*VideoResolution
	// 支持的帧率，fps
	FrameRateRange *IntRange
	// 支持的编码间隔的范围
	EncodingIntervalRange *IntRange
}

// Mpeg4Options mpeg4 选项
type Mpeg4Options struct {
	// 图像大小的支持列表
	ResolutionsAvailable []*VideoResolution
	// 支持的视频帧长度组，通常是 I 帧的距离
	GovLengthRange *IntRange
	// 支持的帧率，fps
	FrameRateRange *IntRange
	// 支持的编码间隔的范围
	EncodingIntervalRange *IntRange
	// 属性的支持列表
	Mpeg4ProfilesSupported []Mpeg4Profile
}

// H264Options h264 选项
type H264Options struct {
	// 图像大小的支持列表
	ResolutionsAvailable []*VideoResolution
	// 支持的视频帧长度组，通常是 I 帧的距离
	GovLengthRange *IntRange
	// 支持的帧率，fps
	FrameRateRange *IntRange
	// 支持的编码间隔的范围
	EncodingIntervalRange *IntRange
	// 属性的支持列表
	H264ProfilesSupported []H264Profile
}

// VideoEncoderOptionsExtension 视频编码器配置选项扩展
type VideoEncoderOptionsExtension struct {
	// jpeg 编码器设置范围
	JPEG []*JpegOptions2
	// mpeg4 编码器设置范围
	MPEG4 []*Mpeg4Options2
	// h264 编码器设置范围
	H264 []*H264Options2
	// 没有说明
	Extension string `xml:",chardata"`
}

// JpegOptions2 jpeg 选项
type JpegOptions2 struct {
	JpegOptions
	// 编码比特率的支持范围
	BitrateRange *IntRange
}

// Mpeg4Options2 mpeg4 选项
type Mpeg4Options2 struct {
	Mpeg4Options
	// 编码比特率的支持范围
	BitrateRange *IntRange
}

// H264Options2 h264 选项
type H264Options2 struct {
	H264Options
	// 编码比特率的支持范围
	BitrateRange *IntRange
}

// AudioSourceConfigurationOptions 音频源配置选项
type AudioSourceConfigurationOptions struct {
	// 音频源配置的 token
	InputTokensAvailable []string
	// 没有说明
	Extension string `xml:",chardata"`
}

// AudioEncoderConfigurationOptions 音频编码器配置选项
type AudioEncoderConfigurationOptions struct {
	// 配置的支持列表
	Options []*AudioEncoderConfigurationOption
}

// AudioEncoderConfigurationOption 音频编码器配置选项
type AudioEncoderConfigurationOption struct {
	// 编码类型
	Encoding AudioEncoding
	// 支持的编码率
	BitrateList IntList
	// 支持的样本率
	SampleRateList IntList
}

// MediaURL 媒体流地址
type MediaURL struct {
	// 用于请求媒体流的地址
	URL string `xml:"Uri"`
	// 在建立连接之前是否有效
	InvalidAfterConnect bool
	// 设备重启后是否无效
	InvalidAfterReboot bool
	// 有效的时间，PT0S 表示永久有效
	Timeout string
}
