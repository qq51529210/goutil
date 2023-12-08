package log

import (
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
type FormatHeader func(log *Log, depth int)

// DefaultHeader 输出 2006-01-02 15:04:05.000000000
func DefaultHeader(log *Log, depth int) {
	FormatTime(log)
}

// FileNameHeader 输出 2006-01-02 15:04:05.000000000 [fileName:fileLine]
func FileNameHeader(log *Log, depth int) {
	// 2006-01-02 15:04:05.000000000
	FormatTime(log)
	log.b = append(log.b, ' ')
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
	log.b = append(log.b, '[')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.Int(line)
	log.b = append(log.b, ']')
}

// FilePathHeader 输出 2006-01-02 15:04:05.000000 [filePath:fileLine]
func FilePathHeader(log *Log, depth int) {
	// 2006-01-02 15:04:05.000000000
	FormatTime(log)
	log.b = append(log.b, ' ')
	// [filePath:fileLine]
	_, path, line, ok := runtime.Caller(depth)
	if !ok {
		path = "???"
		line = -1
	}
	log.b = append(log.b, '[')
	log.b = append(log.b, path...)
	log.b = append(log.b, ':')
	log.Int(line)
	log.b = append(log.b, ']')
}
