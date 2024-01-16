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
	// 是一个 any 看不懂
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

type SystemDateTimeExtension string

type Capabilities struct {
	Analytics *AnalyticsCapabilities
	Device    *DeviceCapabilities
	Events    *EventCapabilities
	Imaging   *ImagingCapabilities
	Media     *MediaCapabilities
	PTZ       *PTZCapabilities
	Extension *CapabilitiesExtension
}

type AnalyticsCapabilities struct {
	XAddr                  string
	RuleSupport            bool
	AnalyticsModuleSupport bool
}

type DeviceCapabilities struct {
	XAddr     string
	Network   *NetworkCapabilities
	System    *SystemCapabilities
	IO        *IOCapabilities
	Security  *SecurityCapabilities
	Extension DeviceCapabilitiesExtension `xml:",chardata"`
}

type NetworkCapabilities struct {
	IPFilter          bool
	ZeroConfiguration bool
	IPVersion6        bool
	DynDNS            bool
	Extension         NetworkCapabilitiesExtension
}

type NetworkCapabilitiesExtension struct {
	Dot11Configuration bool
	Extension          NetworkCapabilitiesExtension2
}

type NetworkCapabilitiesExtension2 map[string]any

type SystemCapabilities struct {
	DiscoveryResolve  bool
	DiscoveryBye      bool
	RemoteDiscovery   bool
	SystemBackup      bool
	SystemLogging     bool
	FirmwareUpgrade   bool
	SupportedVersions OnvifVersion
	//
	Extension SystemCapabilitiesExtension
}

type OnvifVersion struct {
	Major int
	Minor int
}

type SystemCapabilitiesExtension struct {
	HttpFirmwareUpgrade    bool
	HttpSystemBackup       bool
	HttpSystemLogging      bool
	HttpSupportInformation bool
	Extension              SystemCapabilitiesExtension2 `xml:",chardata"`
}

type SystemCapabilitiesExtension2 string

type IOCapabilities struct {
	InputConnectors int
	RelayOutputs    int
	Extension       IOCapabilitiesExtension
}

type IOCapabilitiesExtension struct {
	Auxiliary         bool
	AuxiliaryCommands string
	Extension         IOCapabilitiesExtension2 `xml:",chardata"`
}

type IOCapabilitiesExtension2 string

type SecurityCapabilities struct {
	TLS1_1               bool
	TLS1_2               bool
	OnboardKeyGeneration bool
	AccessPolicyConfig   bool
	X509Token            bool `xml:"X.509Token"`
	SAMLToken            bool
	KerberosToken        bool
	RELToken             bool
	Extension            SecurityCapabilitiesExtension
}

type SecurityCapabilitiesExtension struct {
	TLS1_0    bool
	Extension SecurityCapabilitiesExtension2
}

type SecurityCapabilitiesExtension2 struct {
	Dot1X              bool
	SupportedEAPMethod int
	RemoteUserHandling bool
}

type DeviceCapabilitiesExtension string

type EventCapabilities struct {
	XAddr                                         string
	WSSubscriptionPolicySupport                   bool
	WSPullPointSupport                            bool
	WSPausableSubscriptionManagerInterfaceSupport bool
}

type ImagingCapabilities struct {
	XAddr string
}

type MediaCapabilities struct {
	XAddr                 string
	StreamingCapabilities RealTimeStreamingCapabilities
	Extension             MediaCapabilitiesExtension
}

type RealTimeStreamingCapabilities struct {
	RTPMulticast bool
	RTP_TCP      bool
	RTP_RTSP_TCP bool
	Extension    RealTimeStreamingCapabilitiesExtension `xml:",chardata"`
}

type RealTimeStreamingCapabilitiesExtension string

type MediaCapabilitiesExtension struct {
	ProfileCapabilities ProfileCapabilities
}

type ProfileCapabilities struct {
	MaximumNumberOfProfiles int
}

type PTZCapabilities struct {
	XAddr string
}

type CapabilitiesExtension struct {
	DeviceIO        *DeviceIOCapabilities
	Display         *DisplayCapabilities
	Recording       *RecordingCapabilities
	Search          *SearchCapabilities
	Replay          *ReplayCapabilities
	Receiver        *ReceiverCapabilities
	AnalyticsDevice *AnalyticsDeviceCapabilities
	Extensions      CapabilitiesExtension2 `xml:",chardata"`
}

type DeviceIOCapabilities struct {
	XAddr        string
	VideoSources int
	VideoOutputs int
	AudioSources int
	AudioOutputs int
	RelayOutputs int
}

type DisplayCapabilities struct {
	XAddr       string
	FixedLayout bool
}

type RecordingCapabilities struct {
	XAddr              string
	ReceiverSource     bool
	MediaProfileSource bool
	DynamicRecordings  bool
	DynamicTracks      bool
	MaxStringLength    int
}

type SearchCapabilities struct {
	XAddr          string
	MetadataSearch bool
}

type ReplayCapabilities struct {
	XAddr string
}

type ReceiverCapabilities struct {
	XAddr                string
	RTP_Multicast        bool
	RTP_TCP              bool
	RTP_RTSP_TCP         bool
	SupportedReceivers   int
	MaximumRTSPURILength int
}

type AnalyticsDeviceCapabilities struct {
	XAddr       string
	RuleSupport bool
	Extension   AnalyticsDeviceExtension `xml:",chardata"`
}

type AnalyticsDeviceExtension string

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
