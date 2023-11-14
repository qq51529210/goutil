package sdp

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// 一些常量
const (
	NetTypeIN   = "IN"
	AddrTypeIP4 = "IP4"
	AddrTypeIP6 = "IP6"
)

var (
	// ErrConnectionFormat 表示 Connection 格式错误
	ErrConnectionFormat = errors.New("error format c=")
)

// Connection 表示连接信息
type Connection struct {
	// 网络类型，一般为 IN ，表示 internet
	NetType string
	// 地址类型
	AddrType string
	// 连接地址
	Address string
}

// Parse 从 line 中解析
func (m *Connection) Parse(line string) error {
	p := strings.Fields(line)
	if len(p) != 3 {
		return ErrConnectionFormat
	}
	m.NetType = p[0]
	m.AddrType = p[1]
	m.Address = p[2]
	//
	return nil
}

// FormatTo 格式化
func (m *Connection) FormatTo(buf *bytes.Buffer) {
	fmt.Fprintf(buf, "c=%s %s %s\r\n",
		m.NetType,
		m.AddrType,
		m.Address)
}
