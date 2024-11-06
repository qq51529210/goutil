package discovery

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	gs "goutil/sync"
	"goutil/uid"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	// 用于读取消息的大小，10k 足够了
	readBufLen = 1024 * 100
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

type discoverer struct {
	// 监听的连接
	conn *net.UDPConn
	// 监听的地址
	addr *net.UDPAddr
	// 保存第一个错误
	err *gs.Signal[error]
	// 发现的地址
	data map[string]int
}

func (d *discoverer) init(iface, addr string) error {
	// 网卡
	ifa, err := net.InterfaceByName(iface)
	if err != nil {
		return err
	}
	// 多播地址
	a, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return err
	}
	d.addr = a
	// 底层连接
	c, err := net.ListenMulticastUDP(d.addr.Network(), ifa, &net.UDPAddr{
		IP:   d.addr.IP,
		Port: 0,
	})
	if err != nil {
		return err
	}
	d.conn = c
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
		if _, err := d.conn.WriteToUDP(buf.Bytes(), d.addr); err != nil {
			// 看看是不是网络超时
			if e, ok := err.(net.Error); ok {
				if e.Timeout() {
					return
				}
			}
			// 错误通知
			d.err.Close(err)
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
			// 错误通知
			d.err.Close(err)
			return
		}
		// 解析，错误继续
		var data result
		if err := xml.Unmarshal(buf[:n], &data); err != nil {
			continue
		}
		// 有这样的
		// http://192.168.31.9/onvif/device_service http://[fe80::869a:40ff:fec0:5c28]/onvif/device_service
		part := strings.Fields(data.Body.ProbeMatches.ProbeMatch.XAddrs)
		for _, p := range part {
			if p != "" {
				if _, ok := d.data[p]; ok {
					continue
				}
				d.data[p] = 1
			}
		}
	}
}

// Discover 阻塞发现 timeout 超时后返回
func (d *discoverer) discover(ctx context.Context, duration time.Duration) ([]string, error) {
	d.data = make(map[string]int)
	d.err = gs.NewSignal[error]()
	// 启动读写
	var wg sync.WaitGroup
	wg.Add(2)
	go d.readRoutine(&wg)
	go d.writeRoutine(&wg)
	// 等待
	select {
	case <-ctx.Done():
		// 调用取消
		d.err.Close(ctx.Err())
	case <-time.After(duration):
		// 超时
		d.err.Close(nil)
	case <-d.err.C:
		// 读写错误
	}
	// 关闭
	d.conn.Close()
	//
	wg.Wait()
	// 错误
	if err := d.err.Result(); err != nil {
		return nil, err
	}
	// 返回数组
	ms := make([]string, 0, len(d.data))
	for k := range d.data {
		ms = append(ms, k)
	}
	return ms, nil
}

// Discover 发现一次，然后关闭 conn
func Discover(ctx context.Context, iface, addr string, duration time.Duration) ([]string, error) {
	d := new(discoverer)
	if err := d.init(iface, addr); err != nil {
		return nil, err
	}
	return d.discover(ctx, duration)
}
