package sip

import (
	"io"
	"net"
)

// udpData 实现 io.Reader ，用于读取 udp 数据包
type udpData struct {
	// udp 数据
	b []byte
	// 数据的大小
	n int
	// 用于保存 read 的下标
	i int
	// 地址
	a *net.UDPAddr
}

// Len 返回剩余的数据
func (p *udpData) Len() int {
	return p.n - p.i
}

// Read 实现 io.Reader
func (p *udpData) Read(buf []byte) (int, error) {
	// 没有数据
	if p.i == p.n {
		return 0, io.EOF
	}
	// 还有数据，copy
	n := copy(buf, p.b[p.i:p.n])
	// 增加下标
	p.i += n
	// 返回
	return n, nil
}
