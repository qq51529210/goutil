package sdp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrOriginFormat 表示 Origin 格式错误
	ErrOriginFormat = errors.New("error format o=")
)

// Origin 表示会话创建者信息
type Origin struct {
	// 用户名
	Username string
	// 会话 ID
	SessionID string
	// 会话版本
	SessionVersion string
	Connection
}

// Parse 从 line 中解析
func (m *Origin) Parse(line string) error {
	p := strings.Fields(line)
	n := len(p)
	if n < 6 {
		return ErrOriginFormat
	}
	n--
	m.Address = p[n]
	n--
	m.AddrType = p[n]
	n--
	m.NetType = p[n]
	n--
	m.SessionVersion = p[n]
	n--
	m.SessionID = p[n]
	p = p[:n]
	if len(p) > 1 {
		m.Username = strings.Join(p[:n], " ")
	} else {
		m.Username = p[0]
	}
	//
	return nil
}

// FormatTo 格式化
func (m *Origin) FormatTo(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "o=%s %s %s %s %s %s\r\n",
		m.Username,
		m.SessionID,
		m.SessionVersion,
		m.NetType,
		m.AddrType,
		m.Address)
}
