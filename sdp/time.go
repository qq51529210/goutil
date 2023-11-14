package sdp

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrSDPTimeFormat 表示 Time 格式错误
	ErrSDPTimeFormat = errors.New("error sdp time format")
)

// Time 表示会话时间
type Time struct {
	// 开始时间，单位秒
	Start int64
	// 结束时间，单位秒
	Stop int64
}

// Parse 从 line 中解析
func (m *Time) Parse(line string) error {
	p := strings.Fields(line)
	if len(p) != 2 {
		return ErrSDPTimeFormat
	}
	var err error
	m.Start, err = strconv.ParseInt(p[0], 10, 64)
	if err != nil {
		return ErrSDPTimeFormat
	}
	m.Stop, err = strconv.ParseInt(p[1], 10, 64)
	if err != nil {
		return ErrSDPTimeFormat
	}
	//
	return nil
}

// FormatTo 格式化
func (m *Time) FormatTo(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "t=%d %d\r\n",
		m.Start,
		m.Stop)
}
