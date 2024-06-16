package sip

import (
	"bytes"
	"net"
)

type tcpConn struct {
	key  connKey
	conn *net.TCPConn
	baseConn
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
