package discovery

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"goutil/uid"
	"net"
	"sync"
	"sync/atomic"
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

type discoverer struct {
	// 监听的连接
	conn *net.UDPConn
	// 监听的地址
	addr *net.UDPAddr
	// 回调
	callback Callback
	// 取消的标记
	cancel int32
}

func (d *discoverer) init(iface, addr string) error {
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

func (d *discoverer) writeRoutine(wg *sync.WaitGroup) {
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
			// 看看是不是网络超时
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					return
				}
			}
			// 是取消的标记
			if atomic.LoadInt32(&d.cancel) == 1 {
				return
			}
			// 其他错误
			d.callback("", err)
			return
		}
		// 休息
		time.Sleep(time.Second)
	}
}

func (d *discoverer) readRoutine(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		// 异常
		recover()
	}()
	buf := make([]byte, readBufLen)
	host := make(map[string]int)
	// 读取
	for {
		n, _, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			// 看看是不是网络超时
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					return
				}
			}
			// 是取消的标记
			if atomic.LoadInt32(&d.cancel) == 1 {
				return
			}
			// 其他错误
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
		addr := data.Body.ProbeMatches.ProbeMatch.XAddrs
		if _, ok := host[addr]; ok {
			continue
		} else {
			host[addr] = 1
		}
		d.callback(addr, err)
	}
}

// Discover 阻塞发现 timeout 超时后返回
func (d *discoverer) discover(ctx context.Context, duration time.Duration) error {
	// 超时
	if duration > 0 {
		dealine := time.Now().Add(duration)
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
	// 等待
	var err error
	select {
	case <-ctx.Done():
		atomic.StoreInt32(&d.cancel, 1)
		err = ctx.Err()
	case <-time.After(duration):
		atomic.StoreInt32(&d.cancel, 1)
	}
	// 关闭
	d.conn.Close()
	//
	wg.Wait()
	//
	return err
}

// Discover 发现一次，然后关闭 conn
func Discover(ctx context.Context, iface, addr string, callback Callback, duration time.Duration) error {
	d := new(discoverer)
	d.callback = callback
	if err := d.init(iface, addr); err != nil {
		return err
	}
	return d.discover(ctx, duration)
}
