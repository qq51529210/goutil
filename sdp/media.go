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
	A map[string][]string
	// 国标视频格式
	// f=
	F string
	// 国标 ssrc
	// y=
	Y string
	// 简化查询
	// a=sendonly / recvonly / sendrecv
	SendRecv string
}

// AddAttr 添加
func (m *Media) AddAttr(k, v string) {
	if m.A == nil {
		m.A = make(map[string][]string)
	}
	a, ok := m.A[k]
	if !ok {
		a = make([]string, 0)
	}
	m.A[k] = append(a, v)
}

// SetAttr 设置
func (m *Media) SetAttr(k string, v []string) {
	if m.A == nil {
		m.A = make(map[string][]string)
	}
	if v == nil || len(v) < 1 {
		m.A[k] = []string{}
	} else {
		m.A[k] = v
	}
}

// GetAttr 返回第一个值
func (m *Media) GetAttr(k string) string {
	a, ok := m.A[k]
	if !ok || len(a) < 1 {
		return ""
	}
	return a[0]
}

func (m *Media) parseA(s string) {
	// 简化
	if s == SendRecv || s == SendOnly || s == RecvOnly {
		m.SendRecv = s
	}
	//
	var k, v string
	i := strings.IndexByte(s, ':')
	if i < 0 {
		k = s
	} else {
		k = s[:i]
		v = s[i+1:]
	}
	as, ok := m.A[k]
	if !ok {
		as = make([]string, 0)
	}
	m.A[k] = append(as, v)
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
	for k, as := range m.A {
		if len(as) < 1 {
			fmt.Fprintf(buf, "a=%s\r\n", k)
			continue
		}
		for _, a := range as {
			if a == "" {
				fmt.Fprintf(buf, "a=%s\r\n", k)
				continue
			}
			fmt.Fprintf(buf, "a=%s:%s\r\n", k, a)
		}
	}
	// f=
	if m.F != "" {
		fmt.Fprintf(buf, "f=%s\r\n", m.F)
	}
	// y=
	if m.Y != "" {
		fmt.Fprintf(buf, "y=%s\r\n", m.Y)
	}
}
