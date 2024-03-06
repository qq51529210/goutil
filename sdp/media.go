package sdp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// 一些常量
const (
	ProtoUDP = "RTP/AVP"
	ProtoTCP = "TCP/RTP/AVP"
)

// sdp 中的 a=
const (
	SendRecv = "sendrecv"
	SendOnly = "sendonly"
	RecvOnly = "recvonly"
)

var (
	// ErrMediaFormat 表示 Media 格式错误
	ErrMediaFormat = errors.New("error format m=")
)

// Media 表示会话时间
type Media struct {
	// 媒体类型，video / audio
	Type string
	// 端口号
	Port string
	// 传输协议
	Proto string
	// 格式列表
	FMT string
	// 连接信息
	// c=
	C *Connection
	// 属性
	// a=
	A []string
	// 其他
	Other map[string][]string
}

// SearchA 查找第一个 a=prefix ，返回剩下的部分
func (m *Media) SearchA(prefix string) string {
	for i := 0; i < len(m.A); i++ {
		s := strings.TrimPrefix(m.A[i], prefix)
		if s != m.A[i] {
			return s
		}
	}
	return ""
}

// SearchAllA 查找所有 a=prefix ，返回剩下的部分
func (m *Media) SearchAllA(prefix string) []string {
	var ss []string
	for i := 0; i < len(m.A); i++ {
		s := strings.TrimPrefix(m.A[i], prefix)
		if s != m.A[i] {
			ss = append(ss, s)
		}
	}
	return ss
}

// AddOther 添加 k=v
func (m *Media) AddOther(k, v string) {
	if m.Other == nil {
		m.Other = make(map[string][]string)
	}
	a, ok := m.Other[k]
	if !ok {
		a = make([]string, 0)
	}
	m.Other[k] = append(a, v)
}

// Parse 从 line 中解析
func (m *Media) Parse(line string) error {
	i := strings.IndexByte(line, ' ')
	if i < 0 {
		return ErrMediaFormat
	}
	m.Type = strings.TrimSpace(line[:i])
	line = strings.TrimSpace(line[i+1:])
	//
	i = strings.IndexByte(line, ' ')
	if i < 0 {
		return ErrMediaFormat
	}
	m.Port = strings.TrimSpace(line[:i])
	line = strings.TrimSpace(line[i+1:])
	//
	i = strings.IndexByte(line, ' ')
	if i < 0 {
		return ErrMediaFormat
	}
	m.Proto = strings.TrimSpace(line[:i])
	line = strings.TrimSpace(line[i+1:])
	//
	m.FMT = strings.TrimSpace(line)
	//
	return nil
}

// FormatTo 格式化
func (m *Media) FormatTo(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "m=%s %s %s %s\r\n",
		m.Type,
		m.Port,
		m.Proto,
		m.FMT)
	// c=
	if m.C != nil {
		m.C.FormatTo(buf)
	}
	// a=
	for _, v := range m.A {
		fmt.Fprintf(buf, "a=%s\r\n", v)
	}
	// 其他
	for k, as := range m.Other {
		for _, v := range as {
			fmt.Fprintf(buf, "%s=%s\r\n", k, v)
		}
	}
}
