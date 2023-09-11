package log

import "sync"

var (
	// 用于格式化整数
	intByte []byte
	// 缓存池
	logPool sync.Pool
)

func init() {
	for i := 0; i < 10; i++ {
		intByte = append(intByte, '0'+byte(i))
	}
	//
	logPool.New = func() any {
		return new(Log)
	}
}

// Log 用于写入一行日志
type Log struct {
	// 缓存
	b []byte
	// 用于整数格式化
	f []byte
}

// Reset 重置缓存
func (l *Log) Reset() {
	l.b = l.b[:0]
}

// Write 实现 io.Writer
func (l *Log) Write(data []byte) (int, error) {
	l.b = append(l.b, data...)
	return len(data), nil
}

// writeReversedInt 写入反转的整数
func (l *Log) writeReversedInt(v int) {
	l.f = l.f[:0]
	// 零
	if v == 0 {
		l.f = append(l.f, '0')
		return
	}
	// 1234 -> buff[4,3,2,1]
	for v > 0 {
		l.f = append(l.f, intByte[v%10])
		v /= 10
	}
}

// reversed 将 l.f 的值反转，然后写入 b
func (l *Log) reversed() {
	for i, j := 0, len(l.f)-1; i < j; i, j = i+1, j-1 {
		l.f[i], l.f[j] = l.f[j], l.f[i]
	}
	l.b = append(l.b, l.f...)
}

// Int 写入整数
func (l *Log) Int(v int) {
	if v < 0 {
		// 负数
		l.writeReversedInt(-v)
		l.f = append(l.f, '-')
	} else {
		// 正数
		l.writeReversedInt(v)
	}
	l.reversed()
}

// IntRightAlign 写入整数，左侧补齐 n 个 0
func (l *Log) IntRightAlign(v, n int) {
	if v < 0 {
		// 负数
		l.writeReversedInt(-v)
		for i := len(l.f); i < n; i++ {
			l.f = append(l.f, '0')
		}
		l.f = append(l.f, '-')
	} else {
		// 正数
		l.writeReversedInt(v)
		for i := len(l.f); i < n; i++ {
			l.f = append(l.f, '0')
		}
	}
	l.reversed()
}

// IntLeftAlign 写入整数，右侧补齐 0
func (l *Log) IntLeftAlign(v, n int) {
	if v < 0 {
		// 负数
		l.writeReversedInt(-v)
		n -= len(l.f)
		l.f = append(l.f, '-')
	} else {
		// 正数
		l.writeReversedInt(v)
		n -= len(l.f)
	}
	l.reversed()
	for i := 0; i < n; i++ {
		l.b = append(l.b, '0')
	}
}

// Text 写入字符串
func (l *Log) Text(s string) {
	l.b = append(l.b, s...)
}

// Bytes 写入字节数组
func (l *Log) Bytes(s []byte) {
	l.b = append(l.b, s...)
}

// Byte 写入字节
func (l *Log) Byte(s byte) {
	l.b = append(l.b, s)
}
