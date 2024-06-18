package sip

import (
	"encoding/binary"
	"fmt"
	"net"
)

type conn interface {
	RemoteAddr() net.Addr
	Network() string
	RemoteIP() string
	RemotePort() int
	RemoteAddrString() string
	write([]byte) error
	isUDP() bool
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

// tcpConn 表示 tcp 连接
type tcpConn struct {
	key connKey
	// 底层连接
	conn *net.TCPConn
	// ip
	ip string
	// 端口
	port int
	// ip:port
	ipport string
}

func (c *tcpConn) init(conn *net.TCPConn) {
	c.conn = conn
	a := c.conn.RemoteAddr().(*net.TCPAddr)
	c.key.Init(a.IP, a.Port)
	c.ip = a.IP.String()
	c.port = a.Port
	c.ipport = fmt.Sprintf("%s:%d", c.ip, c.port)
}

func (c *tcpConn) write(b []byte) (err error) {
	_, err = c.conn.Write(b)
	return
}

func (c *tcpConn) isUDP() bool {
	return false
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *tcpConn) Network() string {
	return "tcp"
}

func (c *tcpConn) RemoteIP() string {
	return c.ip
}

func (c *tcpConn) RemotePort() int {
	return c.port
}

func (c *tcpConn) RemoteAddrString() string {
	return c.ipport
}

// udpConn 表示 udp 连接
type udpConn struct {
	// 底层连接
	conn *net.UDPConn
	// 地址
	addr *net.UDPAddr
	// ip
	ip string
	// 端口
	port int
	// ip:port
	ipport string
}

func (c *udpConn) initAddr(a *net.UDPAddr) {
	c.addr = a
	c.ip = a.IP.String()
	c.port = a.Port
	c.ipport = fmt.Sprintf("%s:%d", c.ip, c.port)
}

func (c *udpConn) isUDP() bool {
	return true
}

func (c *udpConn) write(b []byte) (err error) {
	_, err = c.conn.WriteToUDP(b, c.addr)
	return
}

func (c *udpConn) RemoteAddr() net.Addr {
	return c.addr
}

func (c *udpConn) Network() string {
	return "udp"
}

func (c *udpConn) RemoteIP() string {
	return c.ip
}

func (c *udpConn) RemotePort() int {
	return c.port
}

func (c *udpConn) RemoteAddrString() string {
	return c.ipport
}
