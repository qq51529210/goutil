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
	Recover func(recover any)
)

const (
	debugLevel = iota
	infoLevel
	warnLevel
	errorLevel
	panicLevel
)

func init() {
	SetLogger(NewLogger(os.Stdout, DefaultHeader, ""))
}

// SetLogger 设置默认的 Logger 和所有的包函数
func SetLogger(lg *Logger) {
	DefaultLogger = lg
	// Debug
	Debug = DefaultLogger.Debug
	Debugf = DefaultLogger.Debugf
	DebugDepth = DefaultLogger.DebugDepth
	DebugfDepth = DefaultLogger.DebugfDepth
	DebugTrace = DefaultLogger.DebugTrace
	DebugfTrace = DefaultLogger.DebugfTrace
	DebugDepthTrace = DefaultLogger.DebugDepthTrace
	DebugfDepthTrace = DefaultLogger.DebugfDepthTrace
	// Info
	Info = DefaultLogger.Info
	Infof = DefaultLogger.Infof
	InfoDepth = DefaultLogger.InfoDepth
	InfofDepth = DefaultLogger.InfofDepth
	InfoTrace = DefaultLogger.InfoTrace
	InfofTrace = DefaultLogger.InfofTrace
	InfoDepthTrace = DefaultLogger.InfoDepthTrace
	InfofDepthTrace = DefaultLogger.InfofDepthTrace
	// Warn
	Warn = DefaultLogger.Warn
	Warnf = DefaultLogger.Warnf
	WarnDepth = DefaultLogger.WarnDepth
	WarnfDepth = DefaultLogger.WarnfDepth
	WarnTrace = DefaultLogger.WarnTrace
	WarnfTrace = DefaultLogger.WarnfTrace
	WarnDepthTrace = DefaultLogger.WarnDepthTrace
	WarnfDepthTrace = DefaultLogger.WarnfDepthTrace
	// Error
	Error = DefaultLogger.Error
	Errorf = DefaultLogger.Errorf
	ErrorDepth = DefaultLogger.ErrorDepth
	ErrorfDepth = DefaultLogger.ErrorfDepth
	ErrorTrace = DefaultLogger.ErrorTrace
	ErrorfTrace = DefaultLogger.ErrorfTrace
	ErrorDepthTrace = DefaultLogger.ErrorDepthTrace
	ErrorfDepthTrace = DefaultLogger.ErrorfDepthTrace
	// Recover
	Recover = DefaultLogger.Recover
}
