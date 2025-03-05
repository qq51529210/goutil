package log

import (
	"fmt"
	"io"
	"runtime"
	"time"
)

const (
	loggerDepth = 3
)

// Logger 默认实现，修改字段注意并发
type Logger struct {
	// 输出
	io.Writer
	// 头格式
	FormatHeader      FormatHeader
	FormatStackHeader FormatStackHeader
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
func NewLogger(writer io.Writer, name, module string, disableLevels ...string) *Logger {
	lg := new(Logger)
	lg.Writer = writer
	lg.FormatHeader = DefaultHeader
	lg.FormatStackHeader = FilePathHeader
	//
	if name != "" {
		lg.name = fmt.Sprintf("[%s]", name)
	}
	if module != "" {
		lg.module = fmt.Sprintf("[%s]", module)
	}
	// 禁用级别
	if len(disableLevels) > 0 {
		lg.DisableLevels(disableLevels...)
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
func (lg *Logger) DisableLevels(levels ...string) {
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
	lg.FormatHeader(l, lg.name, lg.module, _PanicLevel)
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
	_, _ = lg.Writer.Write(l.b)
	// 回收
	logPool.Put(b)
	logPool.Put(l)
	//
	return true
}

func (lg *Logger) log(level int, trace string, cost time.Duration, text string) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.FormatHeader(l, lg.name, lg.module, level)
	// trace
	if trace != "" {
		l.b = append(l.b, ' ')
		l.b = append(l.b, '[')
		l.b = append(l.b, trace...)
		l.b = append(l.b, ']')
	}
	// cost
	if cost > 0 {
		l.b = append(l.b, ' ')
		l.b = append(l.b, '[')
		l.b = append(l.b, cost.String()...)
		l.b = append(l.b, ']')
	}
	l.b = append(l.b, ' ')
	l.b = append(l.b, text...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	_, _ = lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) stackLog(depth, level int, trace string, cost time.Duration, text string) {
	l := logPool.Get().(*Log)
	l.b = l.b[:0]
	// 头
	lg.FormatStackHeader(l, lg.name, lg.module, level, loggerDepth+depth)
	// trace
	if trace != "" {
		l.b = append(l.b, ' ')
		l.b = append(l.b, '[')
		l.b = append(l.b, trace...)
		l.b = append(l.b, ']')
	}
	// cost
	if cost > 0 {
		l.b = append(l.b, ' ')
		l.b = append(l.b, '[')
		l.b = append(l.b, cost.String()...)
		l.b = append(l.b, ']')
	}
	l.b = append(l.b, ' ')
	l.b = append(l.b, text...)
	// 换行
	l.b = append(l.b, '\n')
	// 输出
	_, _ = lg.Writer.Write(l.b)
	// 回收
	logPool.Put(l)
}

func (lg *Logger) Debug(stack int, trace string, cost time.Duration, text string) {
	if !lg.DisableDebug {
		if stack < 0 {
			lg.log(_DebugLevel, trace, cost, text)
		} else {
			lg.stackLog(stack, _DebugLevel, trace, cost, text)
		}
	}
}

func (lg *Logger) Debugf(stack int, trace string, cost time.Duration, formatOrText string, args ...any) {
	if !lg.DisableDebug {
		var s string
		if len(args) > 0 {
			s = fmt.Sprintf(formatOrText, args...)
		} else {
			s = formatOrText
		}
		if stack < 0 {
			lg.log(_DebugLevel, trace, cost, s)
		} else {
			lg.stackLog(stack, _DebugLevel, trace, cost, s)
		}
	}
}

func (lg *Logger) Info(stack int, trace string, cost time.Duration, text string) {
	if !lg.DisableInfo {
		if stack < 0 {
			lg.log(_InfoLevel, trace, cost, text)
		} else {
			lg.stackLog(stack, _InfoLevel, trace, cost, text)
		}
	}
}

func (lg *Logger) Infof(stack int, trace string, cost time.Duration, formatOrText string, args ...any) {
	if !lg.DisableInfo {
		var s string
		if len(args) > 0 {
			s = fmt.Sprintf(formatOrText, args...)
		} else {
			s = formatOrText
		}
		if stack < 0 {
			lg.log(_InfoLevel, trace, cost, s)
		} else {
			lg.stackLog(stack, _InfoLevel, trace, cost, s)
		}
	}
}

func (lg *Logger) Warn(stack int, trace string, cost time.Duration, text string) {
	if !lg.DisableWarn {
		if stack < 0 {
			lg.log(_WarnLevel, trace, cost, text)
		} else {
			lg.stackLog(stack, _WarnLevel, trace, cost, text)
		}
	}
}

func (lg *Logger) Warnf(stack int, trace string, cost time.Duration, formatOrText string, args ...any) {
	if !lg.DisableWarn {
		var s string
		if len(args) > 0 {
			s = fmt.Sprintf(formatOrText, args...)
		} else {
			s = formatOrText
		}
		if stack < 0 {
			lg.log(_WarnLevel, trace, cost, s)
		} else {
			lg.stackLog(stack, _WarnLevel, trace, cost, s)
		}
	}
}

func (lg *Logger) Errorf(stack int, trace string, cost time.Duration, formatOrText string, args ...any) {
	if !lg.DisableError {
		var s string
		if len(args) > 0 {
			s = fmt.Sprintf(formatOrText, args...)
		} else {
			s = formatOrText
		}
		if stack < 0 {
			lg.log(_ErrorLevel, trace, cost, s)
		} else {
			lg.stackLog(stack, _ErrorLevel, trace, cost, s)
		}
	}
}

func (lg *Logger) Error(stack int, trace string, cost time.Duration, text string) {
	if !lg.DisableError {
		if stack < 0 {
			lg.log(_ErrorLevel, trace, cost, text)
		} else {
			lg.stackLog(stack, _ErrorLevel, trace, cost, text)
		}
	}
}
