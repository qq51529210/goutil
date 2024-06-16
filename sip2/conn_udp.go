package sip

import (
	"bytes"
	"net"
)

type udpConn struct {
	conn *net.UDPConn
	addr *net.UDPAddr
	baseConn
}

func (c *udpConn) write(b []byte) error {
	_, err := c.conn.WriteToUDP(b, c.addr)
	return err
}

func (c *udpConn) writeMsg(msg *Message) error {
	var buf bytes.Buffer
	msg.Enc(&buf)
	_, err := c.conn.WriteToUDP(buf.Bytes(), c.addr)
	return err
}
