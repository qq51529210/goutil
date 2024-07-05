package log

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// FormatTime 格式化 "2006-01-02 15:04:05.000000"
func FormatTime(log *Log) {
	// 不使用 time 标准库，快一点
	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()
	// Date
	log.b = append(log.b, '[')
	log.IntLeftAlign(year, 4)
	log.b = append(log.b, '-')
	log.IntRightAlign(int(month), 2)
	log.b = append(log.b, '-')
	log.IntRightAlign(day, 2)
	log.b = append(log.b, ' ')
	// Time
	log.IntRightAlign(hour, 2)
	log.b = append(log.b, ':')
	log.IntRightAlign(minute, 2)
	log.b = append(log.b, ':')
	log.IntRightAlign(second, 2)
	// Nanosecond
	log.b = append(log.b, '.')
	log.IntLeftAlign(t.Nanosecond(), 9)
	log.b = append(log.b, ']')
}

// FormatHeader 用于格式化日志头
// depth < 0 用于表示没有堆栈，用于 panic 的
type FormatHeader func(log *Log, name, module string, level, depth int)

// DefaultHeader 输出 [level] [2006-01-02 15:04:05.000000000] [name] [module]
func DefaultHeader(log *Log, name, module string, level, depth int) {
	// 级别
	log.b = append(log.b, levels[level]...)
	// 时间
	FormatTime(log)
	// 名称
	if name != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, name...)
	}
	// 模块
	if module != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, module...)
	}
}

// FileNameHeader 输出 [level] [2006-01-02 15:04:05.000000000] [name] [module] [fileName:fileLine]
func FileNameHeader(log *Log, name, module string, level, depth int) {
	// 级别
	log.b = append(log.b, levels[level]...)
	// 时间
	FormatTime(log)
	// 名称
	if name != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, name...)
	}
	// 模块
	if module != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, module...)
	}
	if depth < 0 {
		return
	}
	// [fileName:fileLine]
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if os.IsPathSeparator(path[i]) {
				path = path[i+1:]
				break
			}
		}
	}
	log.b = append(log.b, ' ')
	log.b = append(log.b, '[')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.Int(line)
	log.b = append(log.b, ']')
}

// FilePathHeader 输出 [level] [2006-01-02 15:04:05.000000000] [name] [module] [filePath:fileLine]
func FilePathHeader(log *Log, name, module string, level, depth int) {
	// 级别
	log.b = append(log.b, levels[level]...)
	// 时间
	FormatTime(log)
	// 名称
	if name != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, name...)
	}
	// 模块
	if module != "" {
		log.b = append(log.b, ' ')
		log.b = append(log.b, module...)
	}
	if depth < 0 {
		return
	}
	// [filePath:fileLine]
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	log.b = append(log.b, ' ')
	log.b = append(log.b, '[')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.Int(line)
	log.b = append(log.b, ']')
}

type FileNameError[T any] struct {
	// 追踪
	Trace string
	// 文件名
	Name string
	// 行号
	Line int
	// 错误字符串
	Err string
	// 自定义数据
	Data T
}

// Error 实现接口
func (e *FileNameError[T]) Error() string {
	return e.Err
}

// Error 返回 [Name:Line] [Trace] Err
func (e *FileNameError[T]) Log() string {
	if e.Trace != "" {
		return fmt.Sprintf("[%s:%d] [%s] %s", e.Name, e.Line, e.Trace, e.Err)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Name, e.Line, e.Err)
}

func NewFileNameError[T any](depth int, trace string, data T, err error) *FileNameError[T] {
	_, path, line, ok := runtime.Caller(depth + 1)
	if !ok {
		path = "???"
		line = -1
	} else {
		for i := len(path) - 1; i > 0; i-- {
			if os.IsPathSeparator(path[i]) {
				path = path[i+1:]
				break
			}
		}
	}
	return &FileNameError[T]{
		Trace: trace,
		Name:  path,
		Line:  line,
		Err:   err.Error(),
		Data:  data,
	}
}

type FilePathError[T any] struct {
	// 追踪
	Trace string
	// 文件路径
	Path string
	// 行号
	Line int
	// 错误字符串
	Err string
	// 自定义数据
	Data T
}

// Error 返回 [Path:Line] Err
func (e *FilePathError[T]) Error() string {
	return e.Err
}

// Error 返回 [Path:Line] [Trace] Err
func (e *FilePathError[T]) Log() string {
	if e.Trace != "" {
		return fmt.Sprintf("[%s:%d] [%s] %s", e.Path, e.Line, e.Trace, e.Err)
	}
	return fmt.Sprintf("[%s:%d] %s", e.Path, e.Line, e.Err)
}

func NewFilePathError[T any](depth int, trace string, data T, err error) *FilePathError[T] {
	_, path, line, ok := runtime.Caller(depth + 1)
	if !ok {
		path = "???"
		line = -1
	}
	return &FilePathError[T]{
		Trace: trace,
		Path:  path,
		Line:  line,
		Err:   err.Error(),
		Data:  data,
	}
}
