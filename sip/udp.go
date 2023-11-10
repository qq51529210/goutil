package sip

import (
	"bytes"
	"fmt"
	"gbgw/util"
	"gbgw/util/log"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

// udpReadData 实现 io.Reader ，用于读取 udp 数据包
type udpReadData struct {
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
func (p *udpReadData) Len() int {
	return p.n - p.i
}

// Read 实现 io.Reader
func (p *udpReadData) Read(buf []byte) (int, error) {
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

type udpWriteData struct {
	// key
	k string
	// 数据
	b bytes.Buffer
	// 地址
	a *net.UDPAddr
	// 发送时间戳
	t int64
	// 过期时间戳
	e int64
}

type udpServer struct {
	w  sync.WaitGroup
	c  *net.UDPConn
	p  sync.Pool
	at util.Map[string, *activeTx]
	pt util.Map[string, *passiveTx]
}

// serveUDP 启动 udp 服务
func (s *Server) serveUDP() error {
	// 初始化地址
	a, err := net.ResolveUDPAddr("udp", s.Addr)
	if err != nil {
		return err
	}
	// 初始化底层连接
	s.udp.c, err = net.ListenUDP(a.Network(), a)
	if err != nil {
		return err
	}
	log.InfofTrace(logTraceUDP, "listen %s", s.Addr)
	// 池
	s.udp.p.New = func() any {
		return new(udpWriteData)
	}
	s.udp.at.Init()
	s.udp.pt.Init()
	// 读取协程
	// n := runtime.NumCPU()
	n := 1
	s.w.Add(n)
	for i := 0; i < n; i++ {
		go s.readUDPRoutine(i)
	}
	s.w.Add(3)
	// 检查
	go s.checkActiveTxTimeoutRoutine(logTraceUDP, &s.udp.at)
	go s.checkPassiveTxTimeoutRoutine(logTraceUDP, &s.udp.pt)
	// 消息重发
	go s.checkWriteUDPRoutine()
	//
	return nil
}

// readUDPRoutine 读取 udp 数据
func (s *Server) readUDPRoutine(i int) {
	logTrace := fmt.Sprintf("%s read routine %d", logTraceUDP, i)
	// 清理
	defer func() {
		// 异常
		log.Recover(recover())
		// 日志
		log.InfoTrace(logTrace, "stop")
		// 结束
		s.w.Done()
	}()
	// 日志
	log.InfoTrace(logTrace, "start")
	// 开始
	var err error
	r := newReader(nil, s.MaxMessageLen)
	d := &udpReadData{b: make([]byte, s.MaxMessageLen)}
	c := &udpConn{conn: s.udp.c}
	for s.isOK() {
		// 读取 udp 数据
		d.n, d.a, err = s.udp.c.ReadFromUDP(d.b)
		if err != nil {
			log.ErrorfTrace(logTrace, "read %v", err)
			continue
		}
		d.i = 0
		r.Reset(d)
		// 地址
		c.initAddr(d.a)
		// 一个数据包可能有多个消息，这里需要循环解析处理
		for s.isOK() {
			// 解析
			m := new(message)
			err = m.Dec(r, s.MaxMessageLen)
			if err != nil {
				if err != io.EOF {
					log.ErrorfTrace(logTrace, "dec message %v", err)
					break
				}
				break
			}
			// 处理
			err = s.handleMsg(c, m, &s.udp.at, &s.udp.pt)
			if err != nil {
				log.ErrorfTrace(logTrace, "handle message %v", err)
				break
			}
		}
	}
}

// checkWriteUDPRoutine 检查超时重发协程
func (s *Server) checkWriteUDPRoutine() {
	logTrace := fmt.Sprintf("%s check rto routine", logTraceUDP)
	// 计时器
	timer := time.NewTimer(s.RTO)
	defer func() {
		// 异常
		log.Recover(recover())
		// 计时器
		timer.Stop()
		// 日志
		log.InfoTrace(logTrace, "stop")
		// 结束
		s.w.Done()
	}()
	// 日志
	log.InfoTrace(logTrace, "start")
	// 开始
	var ts []*activeTx
	at := &s.udp.at
	for s.isOK() {
		// 时间到
		now := <-timer.C
		// 组装
		ts = ts[:0]
		at.RLock()
		for _, t := range at.D {
			ts = append(ts, t)
		}
		at.RUnlock()
		// 并发计算
		n := runtime.NumCPU()
		for len(ts) > n {
			m := len(ts) / n
			s.udp.w.Add(1)
			s.writeUDPRoutine(ts[:m], now)
			ts = ts[m:]
		}
		if len(ts) > 0 {
			s.udp.w.Add(1)
			s.writeUDPRoutine(ts, now)
		}
		// 等待并发结束
		s.udp.w.Wait()
		// 重置计时器
		timer.Reset(s.RTO)
	}
}

// writeUDPRoutine 发送 udp 数据
func (s *Server) writeUDPRoutine(ts []*activeTx, now time.Time) {
	defer func() {
		// 异常
		log.Recover(recover())
		// 结束
		s.udp.w.Done()
	}()
	// 循环检查，然后发送，超时移除
	for _, t := range ts {
		// 超时
		if now.Sub(t.writeTime) >= t.rto {
			err := t.conn.write(t.writeData.Bytes())
			if err != nil {
				log.ErrorTrace(logTraceUDP, err)
			} else {
				// 保存发送时间
				t.writeTime = now
				// rto 倍增
				if t.rto < s.MaxRTO {
					t.rto *= 2
					if t.rto > s.MaxRTO {
						t.rto = s.MaxRTO
					}
				}
				//
				log.DebugfTrace(t.key, "retransmission rto %v", t.rto)
			}
		}
	}
}

// closeUDP 关闭 udp 端口
func (s *Server) closeUDP() {
	if s.udp.c != nil {
		log.InfoTrace(logTraceUDP, "close")
		s.udp.c.Close()
		s.udp.c = nil
		//
		for _, t := range s.udp.at.TakeAll() {
			t.Finish(errServerClosed)
		}
		for _, t := range s.udp.pt.TakeAll() {
			t.Finish(errServerClosed)
		}
	}
}
