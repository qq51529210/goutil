package log

// Level 日志级别
type Level string

// 日志级别
const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	WarnLevel  = "warn"
	ErrorLevel = "error"
)

// Disable 返回这个级别以下禁用的
func (s Level) Disable() []string {
	var disableLevels []string
	switch s {
	case "all":
		disableLevels = []string{DebugLevel, InfoLevel, WarnLevel, ErrorLevel}
	case ErrorLevel:
		disableLevels = []string{DebugLevel, InfoLevel, WarnLevel}
	case WarnLevel:
		disableLevels = []string{DebugLevel, InfoLevel}
	case InfoLevel:
		disableLevels = []string{DebugLevel}
	default:
	}
	return disableLevels
}
