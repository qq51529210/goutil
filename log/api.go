package log

import (
	"os"
	"time"
)

var (
	// 级别
	levels = []string{"[D] ", "[I] ", "[W] ", "[E] ", "[P] "}
	// DefaultLogger 默认
	DefaultLogger *Logger
)

// 包接口
var (
	Debug func(trace string, cost time.Duration, args ...any)
	Info  func(trace string, cost time.Duration, args ...any)
	Warn  func(trace string, cost time.Duration, args ...any)
	Error func(trace string, cost time.Duration, args ...any)
	//
	Debugf func(trace string, cost time.Duration, format string, args ...any)
	Infof  func(trace string, cost time.Duration, format string, args ...any)
	Warnf  func(trace string, cost time.Duration, format string, args ...any)
	Errorf func(trace string, cost time.Duration, format string, args ...any)
	//
	DebugStack func(depth int, trace string, cost time.Duration, args ...any)
	InfoStack  func(depth int, trace string, cost time.Duration, args ...any)
	WarnStack  func(depth int, trace string, cost time.Duration, args ...any)
	ErrorStack func(depth int, trace string, cost time.Duration, args ...any)
	//
	DebugfStack func(depth int, trace string, cost time.Duration, format string, args ...any)
	InfofStack  func(depth int, trace string, cost time.Duration, format string, args ...any)
	WarnfStack  func(depth int, trace string, cost time.Duration, format string, args ...any)
	ErrorfStack func(depth int, trace string, cost time.Duration, format string, args ...any)
	//
	Recover func(recover any) bool
)

const (
	_DebugLevel = iota
	_InfoLevel
	_WarnLevel
	_ErrorLevel
	_PanicLevel
)

func init() {
	SetLogger(NewLogger(os.Stdout, "", ""))
}

// SetLogger 设置默认的 Logger 和所有的包函数
func SetLogger(lg *Logger) {
	Recover = lg.Recover
	//
	Debug = lg.Debug
	Info = lg.Info
	Warn = lg.Warn
	Error = lg.Error
	//
	Debugf = lg.Debugf
	Infof = lg.Infof
	Warnf = lg.Warnf
	Errorf = lg.Errorf
	//
	DebugStack = lg.DebugStack
	InfoStack = lg.InfoStack
	WarnStack = lg.WarnStack
	ErrorStack = lg.ErrorStack
	//
	DebugfStack = lg.DebugfStack
	InfofStack = lg.InfofStack
	WarnfStack = lg.WarnfStack
	ErrorfStack = lg.ErrorfStack
	//
	DefaultLogger = lg
}
