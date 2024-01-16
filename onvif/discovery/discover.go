package discovery

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"goutil/log"
	"net"
	"sync/atomic"
	"time"

	"goutil/uid"
)

const (
	// 发送的消息，就差 uuid 了
	msgFmt = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope" xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
  <Header xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
    <a:MessageID>uuid:%s</a:MessageID>
    <a:To>urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
    <a:Action>http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
  </Header>
  <Body>
    <Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<Types>tds:Device</Types>
    </Probe>
  </Body>
</Envelope>`
	// 用于读取消息的大小，10k 足够了
	readBufLen = 1024 * 10
)

// Discover 用于探测
type Discover struct {
	// 是否正常
	running int32
	// 退出探测协程信号
	quit chan struct{}
}

// Run 启动协程后台探测，通过 handle 回调
func (d *Discover) Run(iface, addr string, dur time.Duration, handle func(addr string)) error {
	var conns []*net.UDPConn
	defer func() {
		for _, c := range conns {
			c.Close()
		}
	}()
	// 多播地址
	mAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}
	// 监听地址
	lAddr := &net.UDPAddr{
		IP:   mAddr.IP,
		Port: 0,
	}
	if iface == "" {
		// 所有网络接口
		ifis, err := net.Interfaces()
		if err != nil {
			return err
		}
		for i := 0; i < len(ifis); i++ {
			ifi := &ifis[i]
			// 过滤接口
			if ifi.Flags&net.FlagLoopback != 0 ||
				ifi.Flags&net.FlagUp != 0 ||
				ifi.Flags&net.FlagMulticast != 0 {
				continue
			}
			// 底层连接
			conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
			if err != nil {
				return err
			}
			conns = append(conns, conn)
		}
	} else {
		ifi, err := net.InterfaceByName(iface)
		if err != nil {
			return err
		}
		// 底层连接
		conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
		if err != nil {
			return err
		}
		conns = append(conns, conn)
	}
	atomic.StoreInt32(&d.running, 1)
	d.quit = make(chan struct{})
	// 启动读写
	for _, conn := range conns {
		go d.readRoutine(conn, handle)
		go d.writeRoutine(conn, mAddr, dur)
	}
	// 等待结束
	<-d.quit
	//
	return nil
}

// Stop 停止探测
func (d *Discover) Stop() {
	if atomic.CompareAndSwapInt32(&d.running, 1, 0) {
		close(d.quit)
	}
}

// writeRoutine 协程中发送消息
func (d *Discover) writeRoutine(conn *net.UDPConn, addr *net.UDPAddr, dur time.Duration) {
	// 计时器
	timer := time.NewTimer(0)
	defer func() {
		// 异常
		log.Recover(recover())
		// 关闭
		conn.Close()
		// 计时器
		timer.Stop()
	}()
	buf := bytes.NewBuffer(nil)
	for {
		select {
		case <-d.quit:
			return
		case <-timer.C:
			buf.Reset()
			fmt.Fprintf(buf, msgFmt, uid.UUID1())
			_, err := conn.WriteTo(buf.Bytes(), addr)
			if err != nil {
				log.Error(err)
			}
			timer.Reset(dur)
		}
	}
}

// envelope 用于解析地址
type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProbeMatches struct {
			ProbeMatch struct {
				XAddrs string
			}
		}
	}
}

// readRoutine 协程中读取消息
func (d *Discover) readRoutine(conn *net.UDPConn, fn func(addr string)) {
	defer func() {
		// 异常
		log.Recover(recover())
	}()
	buf := make([]byte, readBufLen)
	for atomic.LoadInt32(&d.running) == 1 {
		// 读取
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Error(err)
			continue
		}
		// 解析
		var res envelope
		err = xml.Unmarshal(buf[:n], &res)
		if err != nil {
			log.Error(err)
			continue
		}
		// 回调通知
		fn(res.Body.ProbeMatches.ProbeMatch.XAddrs)
	}
}
