package log

import "os"

var (
	// 级别
	levels = []string{"[D] ", "[I] ", "[W] ", "[E] ", "[P] "}
	// DefaultLogger 默认
	DefaultLogger *Logger
)

// 包接口
var (
	// Debug
	Debug            func(args ...any)
	Debugf           func(format string, args ...any)
	DebugDepth       func(depth int, args ...any)
	DebugfDepth      func(depth int, format string, args ...any)
	DebugTrace       func(traceID string, args ...any)
	DebugfTrace      func(traceID, format string, args ...any)
	DebugDepthTrace  func(depth int, traceID string, args ...any)
	DebugfDepthTrace func(depth int, traceID, format string, args ...any)
	// Info
	Info            func(args ...any)
	Infof           func(format string, args ...any)
	InfoDepth       func(depth int, args ...any)
	InfofDepth      func(depth int, format string, args ...any)
	InfoTrace       func(traceID string, args ...any)
	InfofTrace      func(traceID, format string, args ...any)
	InfoDepthTrace  func(depth int, traceID string, args ...any)
	InfofDepthTrace func(depth int, traceID, format string, args ...any)
	// Warn
	Warn            func(args ...any)
	Warnf           func(format string, args ...any)
	WarnDepth       func(depth int, args ...any)
	WarnfDepth      func(depth int, format string, args ...any)
	WarnTrace       func(traceID string, args ...any)
	WarnfTrace      func(traceID, format string, args ...any)
	WarnDepthTrace  func(depth int, traceID string, args ...any)
	WarnfDepthTrace func(depth int, traceID, format string, args ...any)
	// Error
	Error            func(args ...any)
	Errorf           func(format string, args ...any)
	ErrorDepth       func(depth int, args ...any)
	ErrorfDepth      func(depth int, format string, args ...any)
	ErrorTrace       func(traceID string, args ...any)
	ErrorfTrace      func(traceID, format string, args ...any)
	ErrorDepthTrace  func(depth int, traceID string, args ...any)
	ErrorfDepthTrace func(depth int, traceID, format string, args ...any)
	// Recover
	Recover func(recover any) bool
)

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
	panicLevel
)

func init() {
	SetLogger(NewLogger(os.Stdout, DefaultHeader, "", "", nil))
}

// SetLogger 设置默认的 Logger 和所有的包函数
func SetLogger(lg *Logger) {
	// Debug
	Debug = lg.Debug
	Debugf = lg.Debugf
	DebugDepth = lg.DebugDepth
	DebugfDepth = lg.DebugfDepth
	DebugTrace = lg.DebugTrace
	DebugfTrace = lg.DebugfTrace
	DebugDepthTrace = lg.DebugDepthTrace
	DebugfDepthTrace = lg.DebugfDepthTrace
	// Info
	Info = lg.Info
	Infof = lg.Infof
	InfoDepth = lg.InfoDepth
	InfofDepth = lg.InfofDepth
	InfoTrace = lg.InfoTrace
	InfofTrace = lg.InfofTrace
	InfoDepthTrace = lg.InfoDepthTrace
	InfofDepthTrace = lg.InfofDepthTrace
	// Warn
	Warn = lg.Warn
	Warnf = lg.Warnf
	WarnDepth = lg.WarnDepth
	WarnfDepth = lg.WarnfDepth
	WarnTrace = lg.WarnTrace
	WarnfTrace = lg.WarnfTrace
	WarnDepthTrace = lg.WarnDepthTrace
	WarnfDepthTrace = lg.WarnfDepthTrace
	// Error
	Error = lg.Error
	Errorf = lg.Errorf
	ErrorDepth = lg.ErrorDepth
	ErrorfDepth = lg.ErrorfDepth
	ErrorTrace = lg.ErrorTrace
	ErrorfTrace = lg.ErrorfTrace
	ErrorDepthTrace = lg.ErrorDepthTrace
	ErrorfDepthTrace = lg.ErrorfDepthTrace
	// Recover
	Recover = lg.Recover
	//
	DefaultLogger = lg
}

// Header 返回日志头格式化函数
// name: fileName/filePath/default
func Header(name string) FormatHeader {
	switch name {
	case "fileName":
		return FileNameHeader
	case "filePath":
		return FilePathHeader
	default:
		return DefaultHeader
	}
}
