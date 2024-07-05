package log

import (
	"fmt"
	"io"
	"runtime"
)

const (
	loggerDepth = 3
)

// Logger 默认实现，修改字段注意并发
type Logger struct {
	// 输出
	io.Writer
	// 头格式
	Header FormatHeader
	// 是否禁止 debug
	DisableDebug bool
	// 是否禁止 info
	DisableInfo bool
	// 是否禁止 warn
	DisableWarn bool
	// 是否禁止 error
	DisableError bool
	// 名称
	name string
	// 模块
	module string
}

// NewLogger 返回默认的 Logger
// 格式 "[name] [level] Header [tracID] text"
func NewLogger(writer io.Writer, header FormatHeader, name, module string, disableLevels []string) *Logger {
	lg := new(Logger)
	lg.Writer = writer
	lg.Header = header
	//
	if name != "" {
		lg.name = fmt.Sprintf("[%s]", name)
	}
	if module != "" {
		lg.module = fmt.Sprintf("[%s]", module)
	}
	// 禁用级别
	if len(disableLevels) > 0 {
		lg.DisableLevels(disableLevels)
	}
	//
	return lg
}

// DisableLevelBelow 禁用 level 以下的级别，[all,debug,info,warn,error]
func (lg *Logger) DisableLevelBelow(level string) {
	switch level {
	case "info":
		lg.DisableDebug = true
	case "warn":
		lg.DisableDebug = true
		lg.DisableInfo = true
	case "error":
		lg.DisableDebug = true
		lg.DisableInfo = true
		lg.DisableWarn = true
	case "all":
		lg.DisableDebug = true
		lg.DisableInfo = true
		lg.DisableWarn = true
		lg.DisableError = true
	}
}

// DisableLevels 禁用级别，[debug,info,warn,error]
func (lg *Logger) DisableLevels(levels []string) {
	for i := 0; i < len(levels); i++ {
		switch levels[i] {
		case "debug":
			lg.DisableDebug = true
		case "info":
			lg.DisableInfo = true
		case "warn":
			lg.DisableWarn = true
		case "error":
			lg.DisableError = true
		}
	}
}

func (lg *Logger) print(depth, level int, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.Header(l, lg.name, lg.module, level, depth)
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printf(depth, level int, format string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.Header(l, lg.name, lg.module, level, depth)
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprintf(l, format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printTrace(depth, level int, trace string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.Header(l, lg.name, lg.module, level, depth)
	l.b = append(l.b, ' ')
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprint(l, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) printfTrace(depth, level int, trace, format string, args ...any) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.Header(l, lg.name, lg.module, level, depth)
	l.b = append(l.b, ' ')
	// 追踪
	l.b = append(l.b, '[')
	l.b = append(l.b, trace...)
	l.b = append(l.b, ']')
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprintf(l, format, args...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

// Recover 如果 recover 不为 nil，输出堆栈
func (lg *Logger) Recover(v any) bool {
	if v == nil {
		return false
	}
	// get stack info l.line
	b := logPool.Get().(*Log)
	b.b = b.b[:cap(b.b)]
	for {
		n := runtime.Stack(b.b, false)
		if n < len(b.b) {
			b.b = b.b[:n]
			break
		}
		b.b = make([]byte, len(b.b)+1024)
	}
	// 缓存
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.Header(l, lg.name, lg.module, panicLevel, -1)
	l.b = append(l.b, ' ')
	// 日志
	fmt.Fprintf(l, "%v", v)
	// 换行
	l.b = append(l.b, '\n')
	// 找到 panic.go
	p := b.b
	found := false
	for len(p) > 0 {
		// find new line
		i := 0
		for ; i < len(p); i++ {
			if p[i] == '\n' {
				i++
				break
			}
		}
		line := p[:i]
		p = p[i:]
		// find file line
		if line[0] != '\t' {
			continue
		}
		if !found {
			found = hasPanicGO(line)
			continue
		}
		// \t filepath/file.go:line +0x622
		i = len(line) - 1
		for i > 0 {
			if line[i] == ' ' {
				//
				line = line[:i]
				break
			}
			i--
		}
		// write
		l.b = append(l.b, line[1:]...)
		l.b = append(l.b, '\n')
	}
	// 输出
	lg.Writer.Write(l.b)
	// 回收
	logPool.Put(b)
	logPool.Put(l)
	//
	return true
}

func hasPanicGO(line []byte) bool {
	for i := len(line) - 1; i > 1; i-- {
		if line[i] == '/' {
			for j := i; j < len(line); j++ {
				if line[j] == 'p' &&
					line[j+1] == 'a' &&
					line[j+2] == 'n' &&
					line[j+3] == 'i' &&
					line[j+4] == 'c' &&
					line[j+5] == '.' &&
					line[j+6] == 'g' &&
					line[j+7] == 'o' {
					return true
				}
			}
			return false
		}
	}
	return false
}

// Debug 输出日志
func (lg *Logger) Debug(args ...any) {
	if !lg.DisableDebug {
		lg.print(loggerDepth, debugLevel, args...)
	}
}

// Debugf 输出日志
func (lg *Logger) Debugf(format string, args ...any) {
	if !lg.DisableDebug {
		lg.printf(loggerDepth, debugLevel, format, args...)
	}
}

// DebugDepth 输出日志
func (lg *Logger) DebugDepth(depth int, args ...any) {
	if !lg.DisableDebug {
		lg.print(loggerDepth+depth, debugLevel, args...)
	}
}

// DebugfDepth 输出日志
func (lg *Logger) DebugfDepth(depth int, format string, args ...any) {
	if !lg.DisableDebug {
		lg.printf(loggerDepth+depth, debugLevel, format, args...)
	}
}

// DebugTrace 输出日志
func (lg *Logger) DebugTrace(traceID string, args ...any) {
	if !lg.DisableDebug {
		lg.printTrace(loggerDepth, debugLevel, traceID, args...)
	}
}

// DebugfTrace 输出日志
func (lg *Logger) DebugfTrace(traceID, format string, args ...any) {
	if !lg.DisableDebug {
		lg.printfTrace(loggerDepth, debugLevel, traceID, format, args...)
	}
}

// DebugDepthTrace 输出日志
func (lg *Logger) DebugDepthTrace(depth int, traceID string, args ...any) {
	if !lg.DisableDebug {
		lg.printTrace(loggerDepth+depth, debugLevel, traceID, args...)
	}
}

// DebugfDepthTrace 输出日志
func (lg *Logger) DebugfDepthTrace(depth int, traceID, format string, args ...any) {
	if !lg.DisableDebug {
		lg.printfTrace(loggerDepth+depth, debugLevel, traceID, format, args...)
	}
}

// Info 输出日志
func (lg *Logger) Info(args ...any) {
	if !lg.DisableInfo {
		lg.print(loggerDepth, infoLevel, args...)
	}
}

// Infof 输出日志
func (lg *Logger) Infof(format string, args ...any) {
	if !lg.DisableInfo {
		lg.printf(loggerDepth, infoLevel, format, args...)
	}
}

// InfoDepth 输出日志
func (lg *Logger) InfoDepth(depth int, args ...any) {
	if !lg.DisableInfo {
		lg.print(loggerDepth+depth, infoLevel, args...)
	}
}

// InfofDepth 输出日志
func (lg *Logger) InfofDepth(depth int, format string, args ...any) {
	if !lg.DisableInfo {
		lg.printf(loggerDepth+depth, infoLevel, format, args...)
	}
}

// InfoTrace 输出日志
func (lg *Logger) InfoTrace(traceID string, args ...any) {
	if !lg.DisableInfo {
		lg.printTrace(loggerDepth, infoLevel, traceID, args...)
	}
}

// InfofTrace 输出日志
func (lg *Logger) InfofTrace(traceID, format string, args ...any) {
	if !lg.DisableInfo {
		lg.printfTrace(loggerDepth, infoLevel, traceID, format, args...)
	}
}

// InfoDepthTrace 输出日志
func (lg *Logger) InfoDepthTrace(depth int, traceID string, args ...any) {
	if !lg.DisableInfo {
		lg.printTrace(loggerDepth+depth, infoLevel, traceID, args...)
	}
}

// InfofDepthTrace 输出日志
func (lg *Logger) InfofDepthTrace(depth int, traceID, format string, args ...any) {
	if !lg.DisableInfo {
		lg.printfTrace(loggerDepth+depth, infoLevel, traceID, format, args...)
	}
}

// Warn 输出日志
func (lg *Logger) Warn(args ...any) {
	if !lg.DisableWarn {
		lg.print(loggerDepth, warnLevel, args...)
	}
}

// Warnf 输出日志
func (lg *Logger) Warnf(format string, args ...any) {
	if !lg.DisableWarn {
		lg.printf(loggerDepth, warnLevel, format, args...)
	}
}

// WarnDepth 输出日志
func (lg *Logger) WarnDepth(depth int, args ...any) {
	if !lg.DisableWarn {
		lg.print(loggerDepth+depth, warnLevel, args...)
	}
}

// WarnfDepth 输出日志
func (lg *Logger) WarnfDepth(depth int, format string, args ...any) {
	if !lg.DisableWarn {
		lg.printf(loggerDepth+depth, warnLevel, format, args...)
	}
}

// WarnTrace 输出日志
func (lg *Logger) WarnTrace(traceID string, args ...any) {
	if !lg.DisableWarn {
		lg.printTrace(loggerDepth, warnLevel, traceID, args...)
	}
}

// WarnfTrace 输出日志
func (lg *Logger) WarnfTrace(traceID, format string, args ...any) {
	if !lg.DisableWarn {
		lg.printfTrace(loggerDepth, warnLevel, traceID, format, args...)
	}
}

// WarnDepthTrace 输出日志
func (lg *Logger) WarnDepthTrace(depth int, traceID string, args ...any) {
	if !lg.DisableWarn {
		lg.printTrace(loggerDepth+depth, warnLevel, traceID, args...)
	}
}

// WarnfDepthTrace 输出日志
func (lg *Logger) WarnfDepthTrace(depth int, traceID, format string, args ...any) {
	if !lg.DisableWarn {
		lg.printfTrace(loggerDepth+depth, warnLevel, traceID, format, args...)
	}
}

// Error 输出日志
func (lg *Logger) Error(args ...any) {
	if !lg.DisableError {
		lg.print(loggerDepth, errorLevel, args...)
	}
}

// Errorf 输出日志
func (lg *Logger) Errorf(format string, args ...any) {
	if !lg.DisableError {
		lg.printf(loggerDepth, errorLevel, format, args...)
	}
}

// ErrorDepth 输出日志
func (lg *Logger) ErrorDepth(depth int, args ...any) {
	if !lg.DisableError {
		lg.print(loggerDepth+depth, errorLevel, args...)
	}
}

// ErrorfDepth 输出日志
func (lg *Logger) ErrorfDepth(depth int, format string, args ...any) {
	if !lg.DisableError {
		lg.printf(loggerDepth+depth, errorLevel, format, args...)
	}
}

// ErrorTrace 输出日志
func (lg *Logger) ErrorTrace(traceID string, args ...any) {
	if !lg.DisableError {
		lg.printTrace(loggerDepth, errorLevel, traceID, args...)
	}
}

// ErrorfTrace 输出日志
func (lg *Logger) ErrorfTrace(traceID, format string, args ...any) {
	if !lg.DisableError {
		lg.printfTrace(loggerDepth, errorLevel, traceID, format, args...)
	}
}

// ErrorDepthTrace 输出日志
func (lg *Logger) ErrorDepthTrace(depth int, traceID string, args ...any) {
	if !lg.DisableError {
		lg.printTrace(loggerDepth+depth, errorLevel, traceID, args...)
	}
}

// ErrorfDepthTrace 输出日志
func (lg *Logger) ErrorfDepthTrace(depth int, traceID, format string, args ...any) {
	if !lg.DisableError {
		lg.printfTrace(loggerDepth+depth, errorLevel, traceID, format, args...)
	}
}
