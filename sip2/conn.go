package sip

import (
	"encoding/binary"
	"net"
)

// conn 抽象一些 conn 的接口，这样调用者好操作
type conn interface {
	Network() string
	RemoteIP() string
	RemotePort() int
	RemoteAddr() string
	write([]byte) error
	writeMsg(*Message) error
	isUDP() bool
}

// baseConn 封装一些公共方法
type baseConn struct {
	remoteIP   string
	remotePort int
	remoteAddr string
}

func (c *baseConn) RemoteIP() string {
	return c.remoteIP
}

func (c *baseConn) RemotePort() int {
	return c.remotePort
}

func (c *baseConn) RemoteAddr() string {
	return c.remoteAddr
}

// connKey 用于连接表，ip + 端口 标识一个连接
type connKey struct {
	// IPV6地址字符数组前64位
	ip1 uint64
	// IPV6地址字符数组后64位
	ip2 uint64
	// 端口
	port uint16
}

// 将128位的ip地址（v4的转成v6）的字节分成两个64位整数，加上端口，作为key
func (k *connKey) Init(ip net.IP, port int) {
	if len(ip) == net.IPv4len {
		k.ip1 = 0
		k.ip2 = uint64(0xff)<<40 | uint64(0xff)<<32 |
			uint64(ip[0])<<24 | uint64(ip[1])<<16 |
			uint64(ip[2])<<8 | uint64(ip[3])
	} else {
		k.ip1 = binary.BigEndian.Uint64(ip[0:])
		k.ip2 = binary.BigEndian.Uint64(ip[8:])
	}
	k.port = uint16(port)
}
