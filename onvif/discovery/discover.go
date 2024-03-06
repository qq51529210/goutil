package discovery

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"goutil/sync"
	"goutil/uid"
	"net"
	"time"
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

// // Run 启动协程后台探测，通过 handle 回调
// func Run(iface, addr string, dur time.Duration, logger *log.Logger, handle func(addr string)) (*Discover, error) {
// 	d := new(Discover)
// 	d.quit = make(chan struct{})
// 	atomic.StoreInt32(&d.running, 1)
// 	// 初始化
// 	mAddr, err := d.initConn(iface, addr)
// 	if err != nil {
// 		d.Stop()
// 		return nil, err
// 	}
// 	// 启动读写
// 	for _, conn := range d.conns {
// 		go d.readRoutine(conn, handle)
// 		go d.writeRoutine(conn, mAddr, dur)
// 	}
// 	//
// 	return d, nil
// }

// // Discover 用于探测
// type Discover struct {
// 	// 打开的端口
// 	conns []*net.UDPConn
// 	// 并发控制
// 	running int32
// 	// 退出信号
// 	quit chan struct{}
// 	// 日志
// 	logger *log.Logger
// }

// // initConn 初始化 conn
// func (d *Discover) initConn(iface, addr string) (*net.UDPAddr, error) {
// 	// 多播地址
// 	mAddr, err := net.ResolveUDPAddr("udp", addr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// 监听地址
// 	lAddr := &net.UDPAddr{
// 		IP:   mAddr.IP,
// 		Port: 0,
// 	}
// 	if iface == "" {
// 		// 所有网络接口
// 		ifis, err := net.Interfaces()
// 		if err != nil {
// 			return nil, err
// 		}
// 		for i := 0; i < len(ifis); i++ {
// 			ifi := &ifis[i]
// 			// 过滤接口
// 			if ifi.Flags&net.FlagLoopback != 0 {
// 				continue
// 			}
// 			// 底层连接
// 			conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
// 			if err != nil {
// 				return nil, err
// 			}
// 			d.conns = append(d.conns, conn)
// 		}
// 	} else {
// 		ifi, err := net.InterfaceByName(iface)
// 		if err != nil {
// 			return nil, err
// 		}
// 		// 底层连接
// 		conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
// 		if err != nil {
// 			return nil, err
// 		}
// 		d.conns = append(d.conns, conn)
// 	}
// 	return mAddr, nil
// }

// // IsRunning 是否运行中
// func (d *Discover) IsRunning() bool {
// 	return atomic.LoadInt32(&d.running) == 1
// }

// // Stop 停止探测
// func (d *Discover) Stop() {
// 	if atomic.CompareAndSwapInt32(&d.running, 1, 0) {
// 		// 通知
// 		close(d.quit)
// 		// 关闭连接
// 		for _, c := range d.conns {
// 			c.Close()
// 		}
// 	}
// }

// // writeRoutine 协程中发送消息
// func (d *Discover) writeRoutine(conn *net.UDPConn, addr *net.UDPAddr, dur time.Duration) {
// 	// 计时器
// 	timer := time.NewTimer(0)
// 	defer func() {
// 		// 异常
// 		d.logger.Recover(recover())
// 		// 关闭
// 		conn.Close()
// 		// 计时器
// 		timer.Stop()
// 	}()
// 	buf := bytes.NewBuffer(nil)
// 	for {
// 		select {
// 		case <-d.quit:
// 			return
// 		case <-timer.C:
// 			buf.Reset()
// 			fmt.Fprintf(buf, msgFmt, uid.UUID1())
// 			_, err := conn.WriteTo(buf.Bytes(), addr)
// 			if err != nil {
// 				d.logger.Error(err)
// 			}
// 			timer.Reset(dur)
// 		}
// 	}
// }

// // envelope 用于解析地址
// type envelope struct {
// 	XMLName xml.Name `xml:"Envelope"`
// 	Body    struct {
// 		ProbeMatches struct {
// 			ProbeMatch struct {
// 				XAddrs string
// 			}
// 		}
// 	}
// }

// // readRoutine 协程中读取消息
// func (d *Discover) readRoutine(conn *net.UDPConn, fn func(addr string)) {
// 	defer func() {
// 		// 异常
// 		d.logger.Recover(recover())
// 	}()
// 	buf := make([]byte, readBufLen)
// 	for atomic.LoadInt32(&d.running) == 1 {
// 		// 读取
// 		n, _, err := conn.ReadFromUDP(buf)
// 		if err != nil {
// 			d.logger.Error(err)
// 			continue
// 		}
// 		// 解析
// 		var res envelope
// 		err = xml.Unmarshal(buf[:n], &res)
// 		if err != nil {
// 			d.logger.Error(err)
// 			continue
// 		}
// 		// 回调通知
// 		fn(res.Body.ProbeMatches.ProbeMatch.XAddrs)
// 	}
// }

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

// Discover 返回发现的地址
// iface 是网卡接口名称，空则使用所有网卡
// addr 是监听的多播地址
// dur 是发送多播消息的间隔
func Discover(ctx context.Context, iface, addr string, dur time.Duration) (string, error) {
	// 初始化
	mAddr, conn, err := initDiscoverConns(iface, addr)
	if err != nil {
		return "", err
	}
	// 结束信号
	resCh := sync.NewChan[string](1)
	errCh := sync.NewChan[error](1)
	// 关闭 conn
	defer func() {
		conn.Close()
		resCh.Close()
		errCh.Close()
	}()
	// 启动读写
	go discoverReadRoutine(conn, resCh, errCh)
	go discoverWriteRoutine(ctx, conn, resCh, errCh, mAddr, dur)
	// 等待结果
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case a := <-resCh.C:
		return a, nil
	case err := <-errCh.C:
		return "", err
	}
}

// initDiscoverConns 初始化地址和监听的conn
func initDiscoverConns(iface, addr string) (*net.UDPAddr, *net.UDPConn, error) {
	// 多播地址
	mAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, nil, err
	}
	// 监听地址
	lAddr := &net.UDPAddr{
		IP:   mAddr.IP,
		Port: 0,
	}
	// var conns []*net.UDPConn
	// if iface == "" {
	// 	// 所有网络接口
	// 	ifis, err := net.Interfaces()
	// 	if err != nil {
	// 		return nil, nil, err
	// 	}
	// 	for i := 0; i < len(ifis); i++ {
	// 		ifi := &ifis[i]
	// 		// 过滤接口
	// 		if ifi.Flags&net.FlagLoopback != 0 {
	// 			continue
	// 		}
	// 		// 底层连接
	// 		conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
	// 		if err != nil {
	// 			return nil, conns, err
	// 		}
	// 		conns = append(conns, conn)
	// 	}
	// } else {
	ifi, err := net.InterfaceByName(iface)
	if err != nil {
		return nil, nil, err
	}
	// 底层连接
	conn, err := net.ListenMulticastUDP(mAddr.Network(), ifi, lAddr)
	if err != nil {
		return nil, nil, err
	}
	// 	conns = append(conns, conn)
	// }
	// return mAddr, conns, nil
	return mAddr, conn, nil
}

func discoverWriteRoutine(ctx context.Context, conn *net.UDPConn, resCh *sync.Chan[string], errCh *sync.Chan[error], addr *net.UDPAddr, dur time.Duration) {
	// 计时器
	timer := time.NewTimer(0)
	defer func() {
		// 异常
		recover()
		// 计时器
		timer.Stop()
	}()
	buf := bytes.NewBuffer(nil)
	for {
		select {
		case <-ctx.Done():
			return
		case <-resCh.C:
			return
		case <-errCh.C:
			return
		case <-timer.C:
			buf.Reset()
			fmt.Fprintf(buf, msgFmt, uid.UUID1())
			_, err := conn.WriteTo(buf.Bytes(), addr)
			if err != nil {
				errCh.Send(err)
				return
			}
			timer.Reset(dur)
		}
	}
}

func discoverReadRoutine(conn *net.UDPConn, resCh *sync.Chan[string], errCh *sync.Chan[error]) {
	defer func() {
		// 异常
		recover()
	}()
	buf := make([]byte, readBufLen)
	// 读取
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		errCh.Send(err)
		return
	}
	// 解析
	var data envelope
	err = xml.Unmarshal(buf[:n], &data)
	if err != nil {
		errCh.Send(err)
		return
	}
	// 成功通知
	resCh.Send(data.Body.ProbeMatches.ProbeMatch.XAddrs)
}
