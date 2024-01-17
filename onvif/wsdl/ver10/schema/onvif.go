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

// SystemDateTime 系统日期时间大全
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
	Extension SystemDateTimeExtension `xml:",chardata"`
}

// UTC 返回 utc 的 go time
func (t *SystemDateTime) UTC() time.Time {
	return t.UTCDateTime.ToTime(time.UTC)
}

// Local 返回 local 的 go time
func (t *SystemDateTime) Local() time.Time {
	return t.LocalDateTime.ToTime(time.Local)
}

// SystemDateTimeExtension 是 SystemDateTime 的 Extension 字段
type SystemDateTimeExtension string

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
	Extension DeviceCapabilitiesExtension `xml:",chardata"`
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

// NetworkCapabilitiesExtension 是 NetworkCapabilities 的 Extension 字段
type NetworkCapabilitiesExtension struct {
	// 没有说明
	Dot11Configuration bool
	// 没有说明
	Extension NetworkCapabilitiesExtension2
}

// NetworkCapabilitiesExtension2 是 NetworkCapabilitiesExtension 的 Extension 字段
type NetworkCapabilitiesExtension2 map[string]any

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

// OnvifVersion 是 SystemCapabilities 的 SupportedVersions 字段
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

// SystemCapabilitiesExtension 是 SystemCapabilities 的 Extension 字段
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
	Extension SystemCapabilitiesExtension2 `xml:",chardata"`
}

// SystemCapabilitiesExtension2 是 SystemCapabilitiesExtension 的 Extension 字段
type SystemCapabilitiesExtension2 string

// IOCapabilities IO 能力
type IOCapabilities struct {
	// 输入的连接器数量
	InputConnectors int
	// 输出的转发数量
	RelayOutputs int
	// 没有说明
	Extension IOCapabilitiesExtension
}

// IOCapabilitiesExtension 是 IOCapabilities 的 Extension 字段
type IOCapabilitiesExtension struct {
	// 没有说明
	Auxiliary bool
	// 没有说明
	AuxiliaryCommands string
	// 没有说明
	Extension IOCapabilitiesExtension2 `xml:",chardata"`
}

// IOCapabilitiesExtension2 是 IOCapabilitiesExtension 的 Extension 字段
type IOCapabilitiesExtension2 string

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

// SecurityCapabilitiesExtension 是 SecurityCapabilities 的 Extension 字段
type SecurityCapabilitiesExtension struct {
	// 没有说明
	TLS10 bool `xml:"TLS1.0"`
	// 没有说明
	Extension SecurityCapabilitiesExtension2
}

// SecurityCapabilitiesExtension2 是 SecurityCapabilitiesExtension 的 Extension 字段
type SecurityCapabilitiesExtension2 struct {
	// 没有说明
	Dot1X bool
	// EAP 方法
	SupportedEAPMethod int
	// 没有说明
	RemoteUserHandling bool
}

// DeviceCapabilitiesExtension 是 DeviceCapabilities 的 Extension 字段
type DeviceCapabilitiesExtension string

// EventCapabilities 时间能力
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
	StreamingCapabilities RealTimeStreamingCapabilities
	// 没有说明
	Extension MediaCapabilitiesExtension
}

// RealTimeStreamingCapabilities 是 MediaCapabilities 的 StreamingCapabilities 字段
type RealTimeStreamingCapabilities struct {
	// rtp 多播
	RTPMulticast bool
	// rtp over tcp
	RTPOverTCP bool `xml:"RTP_TCP"`
	// rtp/rtsp/tcp
	RTPRTSPTCP bool `xml:"RTP_RTSP_TCP"`
	// 没有说明
	Extension RealTimeStreamingCapabilitiesExtension `xml:",chardata"`
}

// RealTimeStreamingCapabilitiesExtension 是 RealTimeStreamingCapabilities 的 Extension 字段
type RealTimeStreamingCapabilitiesExtension string

// MediaCapabilitiesExtension 是 MediaCapabilities 的 Extension 字段
type MediaCapabilitiesExtension struct {
	ProfileCapabilities ProfileCapabilities
}

// ProfileCapabilities 是 MediaCapabilitiesExtension 的 ProfileCapabilities 字段
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
	Extensions CapabilitiesExtension2 `xml:",chardata"`
}

// DeviceIOCapabilities 是 CapabilitiesExtension 的 DeviceIO 字段
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

// DisplayCapabilities 是 CapabilitiesExtension 的 Display 字段
type DisplayCapabilities struct {
	// 没有说明
	XAddr string
	// SetLayout 命令仅支持预定义布局
	FixedLayout bool
}

// RecordingCapabilities 是 CapabilitiesExtension 的 Recording 字段
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

// SearchCapabilities 是 CapabilitiesExtension 的 Search 字段
type SearchCapabilities struct {
	// 没有说明
	XAddr string
	// 没有说明
	MetadataSearch bool
}

// ReplayCapabilities 是 CapabilitiesExtension 的 Replay 字段
type ReplayCapabilities struct {
	// 服务地址
	XAddr string
}

// ReceiverCapabilities 是 CapabilitiesExtension 的 Receiver 字段
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

// AnalyticsDeviceCapabilities 是 CapabilitiesExtension 的 AnalyticsDevice 字段
type AnalyticsDeviceCapabilities struct {
	// 没有说明
	XAddr string
	// 过时的
	RuleSupport bool
	// 没有说明
	Extension AnalyticsDeviceExtension `xml:",chardata"`
}

// AnalyticsDeviceExtension 是 AnalyticsDeviceCapabilities 的 Extension 字段
type AnalyticsDeviceExtension string

// CapabilitiesExtension2 是 CapabilitiesExtension 的 Extensions 字段
type CapabilitiesExtension2 string

type IntRectangle struct {
	X      int `xml:"x,attr"`
	Y      int `xml:"y,attr"`
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

type Profile struct {
	Name                        string
	VideoSourceConfiguration    *VideoSourceConfiguration
	AudioSourceConfiguration    *AudioSourceConfiguration
	VideoEncoderConfiguration   *VideoEncoderConfiguration
	AudioEncoderConfiguration   *AudioEncoderConfiguration
	VideoAnalyticsConfiguration *VideoAnalyticsConfiguration
	PTZConfiguration            *PTZConfiguration
	MetadataConfiguration       *MetadataConfiguration
	Extension                   *ProfileExtension
	Token                       string `xml:"token,attr"`
	Fixed                       bool   `xml:"fixed,attr"`
}

type VideoSourceConfiguration struct {
	ConfigurationEntity
	ViewMode    string `xml:",attr"`
	SourceToken string
	Bounds      IntRectangle
	Extension   *VideoSourceConfigurationExtension
}

type ConfigurationEntity struct {
	Name     string
	UseCount int
	Token    string `xml:"token,attr"`
}

type VideoSourceConfigurationExtension struct {
	Rotate    Rotate
	Extension VideoSourceConfigurationExtension2
}

type VideoSourceConfigurationExtension2 struct {
	LensDescription  LensDescription
	SceneOrientation SceneOrientation
}

type LensDescription struct {
	Offset      LensOffset
	Projection  LensProjection
	XFactor     float64
	FocalLength float64 `xml:",attr"`
}

type LensOffset struct {
	X float64 `xml:"x,attr"`
	Y float64 `xml:"y,attr"`
}

type LensProjection struct {
	Angle         float64
	Radius        float64
	Transmittance float64
}

type SceneOrientationMode string

const (
	SceneOrientationModeMANUAL SceneOrientationMode = "MANUAL"
	SceneOrientationModeAUTO   SceneOrientationMode = "AUTO"
)

type SceneOrientation struct {
	Mode        SceneOrientationMode
	Orientation string
}

type RotateMode string

const (
	RotateModeON   RotateMode = "ON"
	RotateModeOFF  RotateMode = "OFF"
	RotateModeAUTO RotateMode = "AUTO"
)

type Rotate struct {
	Mode      RotateMode
	Degree    int
	Extension RotateExtension `xml:",chardata"`
}

type RotateExtension string

type AudioSourceConfiguration struct {
	ConfigurationEntity
	SourceToken string
}

type VideoEncoding string

const (
	VideoEncodingJPEG  VideoEncoding = "JPEG"
	VideoEncodingMPEG4 VideoEncoding = "MPEG4"
	VideoEncodingH264  VideoEncoding = "H264"
)

type VideoEncoderConfiguration struct {
	ConfigurationEntity
	Encoding       VideoEncoding
	Resolution     VideoResolution
	Quality        float64
	RateControl    *VideoRateControl
	MPEG4          *Mpeg4Configuration
	H264           *H264Configuration
	Multicast      MulticastConfiguration
	SessionTimeout string
}

type VideoResolution struct {
	Width  int
	Height int
}

type VideoRateControl struct {
	FrameRateLimit   int
	EncodingInterval int
	BitrateLimit     int
}

type Mpeg4Profile string

const (
	Mpeg4ProfileSP  Mpeg4Profile = "SP"
	Mpeg4ProfileASP Mpeg4Profile = "ASP"
)

type Mpeg4Configuration struct {
	GovLength    int
	Mpeg4Profile Mpeg4Profile
}

type H264Profile string

const (
	H264ProfileBaseline H264Profile = "Baseline"
	H264ProfileMain     H264Profile = "Main"
	H264ProfileExtended H264Profile = "Extended"
	H264ProfileHigh     H264Profile = "High"
)

type H264Configuration struct {
	GovLength   int
	H264Profile H264Profile
}

type MulticastConfiguration struct {
	Address   IPAddress
	Port      int
	TTL       int
	AutoStart bool
}

type IPType string

const (
	IPTypeIPv4 IPType = "IPv4"
	IPTypeIPv6 IPType = "IPv6"
)

type IPAddress struct {
	Type        IPType
	IPv4Address string
	IPv6Address string
}

type AudioEncoding string

const (
	AudioEncodingG711 AudioEncoding = "G711"
	AudioEncodingG726 AudioEncoding = "G726"
	AudioEncodingACC  AudioEncoding = "ACC"
)

type AudioEncoderConfiguration struct {
	ConfigurationEntity
	Encoding       AudioEncoding
	Bitrate        int
	SampleRate     int
	Multicast      MulticastConfiguration
	SessionTimeout string
}

type VideoAnalyticsConfiguration struct {
	ConfigurationEntity
	AnalyticsEngineConfiguration AnalyticsEngineConfiguration
	RuleEngineConfiguration      RuleEngineConfiguration
}

type AnalyticsEngineConfiguration struct {
	AnalyticsModule Config
	Extension       AnalyticsEngineConfigurationExtension `xml:",chardata"`
}

type Config struct {
	Parameters ItemList
	Name       string `xml:",attr"`
	Type       string `xml:",attr"`
}

type ItemList struct {
	SimpleItem  []SimpleItem
	ElementItem []ElementItem
	Extension   ItemListExtension `xml:",chardata"`
}

type SimpleItem struct {
	Name  string `xml:",attr"`
	Value string `xml:",attr"`
}

type ElementItem struct {
	// Item name.
	Name  string `xml:",attr"`
	Value string `xml:",chardata"`
}

type ItemListExtension string

type AnalyticsEngineConfigurationExtension string

type RuleEngineConfiguration struct {
	Rule      Config
	Extension RuleEngineConfigurationExtension `xml:",chardata"`
}

type RuleEngineConfigurationExtension string

type PTZConfiguration struct {
	ConfigurationEntity
	NodeToken                              string
	DefaultAbsolutePantTiltPositionSpace   string
	DefaultAbsoluteZoomPositionSpace       string
	DefaultRelativePanTiltTranslationSpace string
	DefaultRelativeZoomTranslationSpace    string
	DefaultContinuousPanTiltVelocitySpace  string
	DefaultContinuousZoomVelocitySpace     string
	DefaultPTZSpeed                        PTZSpeed
	DefaultPTZTimeout                      string
	PanTiltLimits                          PanTiltLimits
	ZoomLimits                             ZoomLimits
	Extension                              PTZConfigurationExtension
	MoveRamp                               int `xml:",attr"`
	PresetRamp                             int `xml:",attr"`
	PresetTourRamp                         int `xml:",attr"`
}

type PTZSpeed struct {
	PanTilt Vector2D
	Zoom    Vector1D
}

type Vector2D struct {
	X     float64 `xml:"x,attr"`
	Y     float64 `xml:"y,attr"`
	Space string  `xml:"space,attr"`
}

type Vector1D struct {
	X     float64 `xml:"x,attr"`
	Space string  `xml:"space,attr"`
}

type PanTiltLimits struct {
	Range Space2DDescription
}

type Space2DDescription struct {
	URI    string
	XRange FloatRange
	YRange FloatRange
}

type FloatRange struct {
	Min float64
	Max float64
}

type ZoomLimits struct {
	Range Space1DDescription
}

type Space1DDescription struct {
	URI    string
	XRange FloatRange
}

type PTZConfigurationExtension struct {
	PTControlDirection PTControlDirection
	Extension          PTZConfigurationExtension2 `xml:",chardata"`
}

type PTControlDirection struct {
	EFlip     EFlip
	Reverse   Reverse
	Extension PTControlDirectionExtension `xml:",chardata"`
}

type EFlipMode string

const (
	EFlipModeOFF      EFlipMode = "OFF"
	EFlipModeON       EFlipMode = "ON"
	EFlipModeExtended EFlipMode = "Extended"
)

type EFlip struct {
	Mode EFlipMode
}

type ReverseMode string

const (
	ReverseModeOFF      ReverseMode = "OFF"
	ReverseModeON       ReverseMode = "ON"
	ReverseModeAUTO     ReverseMode = "AUTO"
	ReverseModeExtended ReverseMode = "Extended"
)

type Reverse struct {
	Mode ReverseMode
}

type PTControlDirectionExtension string

type PTZConfigurationExtension2 string

type MetadataConfiguration struct {
	ConfigurationEntity
	PTZStatus                    *PTZFilter
	Events                       *EventSubscription
	Analytics                    bool
	Multicast                    MulticastConfiguration
	SessionTimeout               string
	AnalyticsEngineConfiguration *AnalyticsEngineConfiguration
	Extension                    *MetadataConfigurationExtension `xml:",chardata"`
	CompressionType              string                          `xml:",attr"`
	GeoLocation                  bool                            `xml:",attr"`
	ShapePolygon                 bool                            `xml:",attr"`
}

type PTZFilter struct {
	Status   bool
	Position bool
}

type EventSubscription struct {
	Filter             FilterType
	SubscriptionPolicy SubscriptionPolicy `xml:",chardata"`
}

type FilterType string

type SubscriptionPolicy string

type MetadataConfigurationExtension string

type ProfileExtension struct {
	AudioOutputConfiguration  *AudioOutputConfiguration
	AudioDecoderConfiguration *AudioDecoderConfiguration
	Extension                 *ProfileExtension2 `xml:",chardata"`
}

type AudioOutputConfiguration struct {
	ConfigurationEntity
	OutputToken string
	SendPrimacy string
	OutputLevel int
}

type AudioDecoderConfiguration struct {
	ConfigurationEntity
}

type ProfileExtension2 string

type IntRectangleRange struct {
	XRange      IntRange
	YRange      IntRange
	WidthRange  IntRange
	HeightRange IntRange
}

type VideoSourceConfigurationOptions struct {
	BoundsRange                IntRectangleRange
	VideoSourceTokensAvailable string
	Extension                  VideoSourceConfigurationOptionsExtension
	MaximumNumberOfProfiles    int `xml:",attr"`
}

type VideoSourceConfigurationOptionsExtension struct {
	Rotate    RotateOptions
	Extension VideoSourceConfigurationOptionsExtension2
}

type VideoSourceConfigurationOptionsExtension2 struct {
	SceneOrientationMode SceneOrientationMode
}

type RotateOptions struct {
	Mode       RotateMode
	DegreeList IntList
	Extension  RotateOptionsExtension `xml:",chardata"`
	Reboot     bool                   `xml:",attr"`
}

type IntList struct {
	Items []int
}

type RotateOptionsExtension string

type IntRange struct {
	Min int
	Max int
}

type VideoEncoderConfigurationOptions struct {
	QualityRange IntRange
	JPEG         JpegOptions
	MPEG4        Mpeg4Options
	H264         H264Options
	Extension    VideoEncoderOptionsExtension
}

type JpegOptions struct {
	ResolutionsAvailable  VideoResolution
	FrameRateRange        IntRange
	EncodingIntervalRange IntRange
}

type Mpeg4Options struct {
	ResolutionsAvailable   VideoResolution
	GovLengthRange         IntRange
	FrameRateRange         IntRange
	EncodingIntervalRange  IntRange
	Mpeg4ProfilesSupported Mpeg4Profile
}

type H264Options struct {
	ResolutionsAvailable  VideoResolution
	GovLengthRange        IntRange
	FrameRateRange        IntRange
	EncodingIntervalRange IntRange
	H264ProfilesSupported H264Profile
}

type VideoEncoderOptionsExtension struct {
	JPEG      JpegOptions2
	MPEG4     Mpeg4Options2
	H264      H264Options2
	Extension VideoEncoderOptionsExtension2
}

type JpegOptions2 struct {
	JpegOptions
	BitrateRange IntRange
}

type Mpeg4Options2 struct {
	Mpeg4Options
	BitrateRange IntRange
}

type H264Options2 struct {
	H264Options
	BitrateRange IntRange
}

type VideoEncoderOptionsExtension2 string

type AudioSourceConfigurationOptions struct {
	InputTokensAvailable string
	Extension            AudioSourceOptionsExtension
}

type AudioSourceOptionsExtension string

type AudioEncoderConfigurationOptions struct {
	Options AudioEncoderConfigurationOption
}

type AudioEncoderConfigurationOption struct {
	Encoding       AudioEncoding
	BitrateList    IntList
	SampleRateList IntList
}
