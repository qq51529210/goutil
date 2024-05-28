package discovery

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"goutil/uid"
	"net"
	"sync"
	"time"
)

const (
	// 用于读取消息的大小，10k 足够了
	readBufLen = 1024 * 10
	// MulticastAddr 默认的使用的多播地址
	MulticastAddr = "239.255.255.250:3702"
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
)

// result 用于解析地址
type result struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProbeMatches struct {
			ProbeMatch struct {
				XAddrs string
			}
		}
	}
}

// Callback 回调，err 是 net.Conn 返回的错误，超时是正常
type Callback func(addr string, err error)

type Discoverer struct {
	// 监听的连接
	conn *net.UDPConn
	// 监听的地址
	addr *net.UDPAddr
	// 回调
	callback Callback
}

func (d *Discoverer) init(iface, addr string) error {
	// 网卡
	ifa, err := net.InterfaceByName(iface)
	if err != nil {
		return err
	}
	// 多播地址
	if d.addr, err = net.ResolveUDPAddr("udp", addr); err != nil {
		return err
	}
	// 底层连接
	if d.conn, err = net.ListenMulticastUDP(d.addr.Network(), ifa, &net.UDPAddr{
		IP:   d.addr.IP,
		Port: 0,
	}); err != nil {
		return err
	}
	//
	return nil
}

func (d *Discoverer) writeRoutine(wg *sync.WaitGroup) {
	// 计时器
	timer := time.NewTimer(0)
	defer func() {
		wg.Done()
		// 异常
		recover()
		// 计时器
		timer.Stop()
	}()
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, msgFmt, uid.UUID1())
	for {
		if _, err := d.conn.WriteTo(buf.Bytes(), d.addr); err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					return
				}
			}
			d.callback("", err)
			return
		}
		// 休息
		time.Sleep(time.Second)
	}
}

func (d *Discoverer) readRoutine(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		// 异常
		recover()
	}()
	buf := make([]byte, readBufLen)
	// 读取
	for {
		n, _, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					return
				}
			}
			d.callback("", err)
			return
		}
		// 解析
		var data result
		if err := xml.Unmarshal(buf[:n], &data); err != nil {
			d.callback("", err)
			return
		}
		// 成功通知
		d.callback(data.Body.ProbeMatches.ProbeMatch.XAddrs, err)
	}
}

// Discover 阻塞发现 timeout 超时后返回
func (d *Discoverer) Discover(timeout time.Duration) error {
	// 超时
	if timeout > 0 {
		dealine := time.Now().Add(timeout)
		if err := d.conn.SetReadDeadline(dealine); err != nil {
			d.conn.Close()
			return err
		}
		if err := d.conn.SetWriteDeadline(dealine); err != nil {
			d.conn.Close()
			return err
		}
	}
	// 启动读写
	var wg sync.WaitGroup
	wg.Add(2)
	go d.readRoutine(&wg)
	go d.writeRoutine(&wg)
	wg.Wait()
	//
	return nil
}

// Close 用于关闭底层的 conn
func (d *Discoverer) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}

func NewDiscoverer(iface, addr string, callback Callback) (*Discoverer, error) {
	d := new(Discoverer)
	d.callback = callback
	return d, d.init(iface, addr)
}
