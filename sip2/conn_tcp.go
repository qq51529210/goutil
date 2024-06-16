package sip

import (
	"bytes"
	"fmt"
	"net"
)

type tcpConn struct {
	key connKey
	// 底层连接
	conn *net.TCPConn
	baseConn
}

func (c *tcpConn) init(conn *net.TCPConn) {
	c.conn = conn
	a := c.conn.RemoteAddr().(*net.TCPAddr)
	c.remoteIP = a.IP.String()
	c.remotePort = a.Port
	c.remoteAddr = fmt.Sprintf("%s:%d", c.remoteIP, c.remotePort)
}

func (c *tcpConn) Network() string {
	return "tcp"
}

func (c *tcpConn) write(b []byte) error {
	_, err := c.conn.Write(b)
	return err
}

func (c *tcpConn) writeMsg(msg *Message) error {
	var buf bytes.Buffer
	msg.Enc(&buf)
	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *tcpConn) isUDP() bool {
	return false
}
