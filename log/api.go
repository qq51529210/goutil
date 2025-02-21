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

type LogFunc func(stack int, trace string, cost time.Duration, formatOrText string, args ...any)

// 包接口
var (
	Debug   LogFunc
	Info    LogFunc
	Warn    LogFunc
	Error   LogFunc
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
	Debug = lg.Debug
	Info = lg.Info
	Warn = lg.Warn
	Error = lg.Error
	DefaultLogger = lg
}
