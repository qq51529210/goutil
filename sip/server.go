package sip

import (
	"context"
	"goutil/log"
	gosync "goutil/sync"
	"io"
	"net"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// logTraceUDP = "sip udp"
	// logTraceTCP = "sip tcp"
	udp = "udp"
	tcp = "tcp"
)

// Handler 是处理消息的接口
type Handler interface {
	// 返回 true 表示已经处理
	// 这样事务超时时间内就不会触发相同消息的回调
	// 而是直接回复第一次的响应消息数据
	HandleRequest(*Request) bool
	HandleResponse(*Response)
}

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

type udpServer struct {
	w  sync.WaitGroup
	c  *net.UDPConn
	at gosync.Map[string, *activeTx]
	pt gosync.Map[string, *passiveTx]
}

// tcpServer 包装 tcp 相关的数据
type tcpServer struct {
	l  *net.TCPListener
	c  gosync.Map[connKey, *tcpConn]
	at gosync.Map[string, *activeTx]
	pt gosync.Map[string, *passiveTx]
}

// Server 表示一个服务
type Server struct {
	// 回调函数
	Handler
	// 监听地址
	Addr string
	// udp 超时重发起始间隔
	MinRTO time.Duration
	// udp 超时重发最大间隔
	MaxRTO time.Duration
	// 事务超时时间
	TxTimeout time.Duration
	// 最大的消息字节数，防止内存爆掉哦
	MaxMessageLen int
	// 用户代理
	UserAgent string
	// 日志
	Logger *log.Logger
	// 用于同步等待协程退出
	w sync.WaitGroup
	// 状态
	ok int32
	// udp
	udp udpServer
	// tcp server
	tcp tcpServer
}

func (s *Server) isOK() bool {
	return atomic.LoadInt32(&s.ok) == 0
}

// Serve 开始服务，不会阻塞
func (s *Server) Serve() error {
	// udp 启动
	if err := s.serveUDP(); err != nil {
		return err
	}
	// tcp 启动
	if err := s.serveTCP(); err != nil {
		return err
	}
	// 日志
	if s.Logger == nil {
		s.Logger = log.DefaultLogger
	}
	//
	return nil
}

// Close 停止服务
func (s *Server) Close() error {
	if atomic.CompareAndSwapInt32(&s.ok, 0, 1) {
		// 关闭双服务
		s.closeUDP()
		s.closeTCP()
		// 等待所有协程退出
		s.w.Wait()
	}
	return nil
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
	s.Logger.Infof("listen udp %s", s.Addr)
	//
	s.udp.at.Init()
	s.udp.pt.Init()
	// 读取协程
	n := runtime.NumCPU()
	s.w.Add(n)
	for i := 0; i < n; i++ {
		go s.readUDPRoutine(i)
	}
	s.w.Add(3)
	// 检查
	go s.checkActiveTxTimeoutRoutine(udp, &s.udp.at)
	go s.checkPassiveTxTimeoutRoutine(udp, &s.udp.pt)
	// 消息重发
	go s.checkWriteUDPRoutine()
	//
	return nil
}

// readUDPRoutine 读取 udp 数据
func (s *Server) readUDPRoutine(i int) {
	// 清理
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.Warnf("udp read routine %d stop", i)
		// 结束
		s.w.Done()
	}()
	// 日志
	s.Logger.Debugf("udp read routine %d start", i)
	// 开始
	var err error
	r := newReader(nil, s.MaxMessageLen)
	d := &udpReadData{b: make([]byte, s.MaxMessageLen)}
	c := &udpConn{conn: s.udp.c}
	for s.isOK() {
		// 读取 udp 数据
		d.n, d.a, err = s.udp.c.ReadFromUDP(d.b)
		if err != nil {
			s.Logger.Errorf("udp read data %v", err)
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
					s.Logger.Errorf("udp dec message %v", err)
				}
				break
			}
			// 处理
			err = s.handleMsg(c, m, &s.udp.at, &s.udp.pt)
			if err != nil {
				s.Logger.Errorf("udp handle message %v", err)
				break
			}
		}
	}
}

// checkWriteUDPRoutine 检查超时重发协程
func (s *Server) checkWriteUDPRoutine() {
	// 计时器
	timer := time.NewTimer(s.MinRTO)
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.Warnf("udp check rto routine stop")
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
	}()
	// 日志
	s.Logger.Debug("udp check rto routine start")
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
			go s.writeUDPRoutine(ts[:m], now)
			ts = ts[m:]
		}
		if len(ts) > 0 {
			s.udp.w.Add(1)
			go s.writeUDPRoutine(ts, now)
		}
		// 等待并发结束
		s.udp.w.Wait()
		// 重置计时器
		timer.Reset(s.MinRTO)
	}
}

// writeUDPRoutine 发送 udp 数据
func (s *Server) writeUDPRoutine(ts []*activeTx, now time.Time) {
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 结束
		s.udp.w.Done()
	}()
	// 循环检查，然后发送，超时移除
	for _, t := range ts {
		if t.stopRT {
			continue
		}
		// 超时
		if now.Sub(t.writeTime) >= t.rto {
			err := t.conn.write(t.writeData.Bytes())
			if err != nil {
				s.Logger.Errorf("udp write %v", err)
				continue
			}
			// 保存发送时间
			t.writeTime = now
			// rto 倍增
			if t.rto < s.MaxRTO {
				t.rto *= 2
				if t.rto > s.MaxRTO {
					t.rto = s.MaxRTO
				}
			}
		}
	}
}

// closeUDP 关闭 udp 端口
func (s *Server) closeUDP() {
	if s.udp.c != nil {
		s.Logger.Info("udp close")
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

// serveTCP 开始 tcp 服务
func (s *Server) serveTCP() error {
	// 初始化
	a, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.tcp.l, err = net.ListenTCP(a.Network(), a)
	if err != nil {
		return err
	}
	s.Logger.Infof("tcp listen %s", s.Addr)
	//
	s.tcp.c.Init()
	s.tcp.at.Init()
	s.tcp.pt.Init()
	//
	s.w.Add(3)
	// 监听
	go s.listenTCPRoutine()
	// 检查
	go s.checkActiveTxTimeoutRoutine(tcp, &s.tcp.at)
	go s.checkPassiveTxTimeoutRoutine(tcp, &s.tcp.pt)
	// 返回
	return nil
}

// listenTCPRoutine 监听 tcp 连接，然后启动协程处理
func (s *Server) listenTCPRoutine() {
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 结束
		s.w.Done()
	}()
	for s.isOK() {
		// 监听
		conn, err := s.tcp.l.AcceptTCP()
		if err != nil {
			s.Logger.Errorf("tcp accept %v", err)
			continue
		}
		// 处理
		s.w.Add(1)
		go s.handleTCPConnRoutine(s.addTCPConn(conn))
	}
}

// addTCPConn 添加并返回
func (s *Server) addTCPConn(conn *net.TCPConn) *tcpConn {
	// 初始化
	c := new(tcpConn)
	c.init(conn)
	// 添加
	s.tcp.c.Set(c.key, c)
	//
	return c
}

// delTCPConn 移除并关闭
func (s *Server) delTCPConn(c *tcpConn) {
	s.tcp.c.Del(c.key)
	c.conn.Close()
}

// getTCPConn 获取
func (s *Server) getTCPConn(a *net.TCPAddr) *tcpConn {
	k := connKey{}
	k.Init(a.IP, a.Port)
	return s.tcp.c.Get(k)
}

// dialTCPConn 创建连接
func (s *Server) dialTCPConn(a *net.TCPAddr) (*tcpConn, error) {
	conn, err := net.DialTimeout(a.Network(), a.String(), s.TxTimeout)
	if err != nil {
		return nil, err
	}
	return s.addTCPConn(conn.(*net.TCPConn)), nil
}

// handleTCPConnRoutine 处理 tcp conn 消息
func (s *Server) handleTCPConnRoutine(c *tcpConn) {
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 移除
		s.delTCPConn(c)
		// 结束
		s.w.Done()
	}()
	//
	r := newReader(c.conn, s.MaxMessageLen)
	for s.isOK() {
		// 解析，错误直接返回关闭连接
		m := new(message)
		err := m.Dec(r, s.MaxMessageLen)
		if err != nil {
			s.Logger.Errorf("tcp dec message %v", err)
			return
		}
		// 处理
		err = s.handleMsg(c, m, &s.tcp.at, &s.tcp.pt)
		if err != nil {
			s.Logger.Errorf("tcp handle message %v", err)
			return
		}
	}
}

// closeTCP 停止监听，关闭所有的连接
func (s *Server) closeTCP() {
	if s.tcp.l != nil {
		s.Logger.Info("tcp close")
		s.tcp.l.Close()
		s.tcp.l = nil
		//
		cs := s.tcp.c.TakeAll()
		for _, c := range cs {
			c.conn.Close()
		}
	}
}

// handleMsg 处理消息
func (s *Server) handleMsg(c conn, m *message, at *gosync.Map[string, *activeTx], pt *gosync.Map[string, *passiveTx]) error {
	// 请求消息
	if m.isReq {
		// 日志
		s.Logger.DebugfTrace(m.txKey(), "request from %s %s\n%v", c.Network(), c.RemoteAddrString(), m)
		// 事务，返回一定不为 nil
		t := s.newPassiveTx(m, pt)
		if atomic.CompareAndSwapInt32(&t.handing, 0, 1) {
			// 在协程中处理
			s.w.Add(1)
			go s.handleRequestRoutine(c, t, m)
		} else {
			// 已经处理过
			if t.done {
				// 有响应数据，直接发送，无需回调
				d := t.writeData.Bytes()
				if len(d) > 0 {
					return c.write(d)
				}
			}
		}
		return nil
	}
	// 日志
	s.Logger.DebugfTrace(m.txKey(), "response from %s %s\n%v", c.Network(), c.RemoteAddrString(), m)
	// 响应消息
	if m.StartLine[1][0] == '1' {
		t := at.Get(m.txKey())
		if t != nil {
			t.stopRT = true
		}
		// 1xx 消息没什么卵用
		return nil
	}
	// 事务，不一定有
	t := s.delActiveTx(m.txKey(), at)
	if t != nil {
		// 在协程中处理
		s.w.Add(1)
		go s.handleResponseRoutine(t, m)
	}
	//
	return nil
}

// handleRequestRoutine 在协程中处理请求消息
func (s *Server) handleRequestRoutine(c conn, t *passiveTx, m *message) {
	old := time.Now()
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.DebugfTrace(t.TxKey(), "[%v] handle request", time.Since(old))
		// 结束
		s.w.Done()
	}()
	// 回调
	t.done = s.HandleRequest(&Request{
		Server:  s,
		message: m,
		tx:      t,
		conn:    c,
	})
	if !t.done {
		// 没有完成，回复标记，等下一次的消息再回调
		atomic.StoreInt32(&t.handing, 0)
	}
}

// handleResponseRoutine 在协程中处理响应消息
func (s *Server) handleResponseRoutine(t *activeTx, m *message) {
	old := time.Now()
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.DebugfTrace(t.TxKey(), "[%v] handle response", time.Since(old))
		// 无论回调有没有通知，这里都通知一下
		t.Finish(nil)
		// 结束
		s.w.Done()
	}()
	// 回调处理
	s.HandleResponse(&Response{
		Server:   s,
		activeTx: t,
		message:  m,
	})
}

// Request 发送请求并等待响应
func (s *Server) Request(ctx context.Context, r *Request, a net.Addr, d any) error {
	// tcp
	if _a, ok := a.(*net.TCPAddr); ok {
		var err error
		// 获取连接
		c := s.getTCPConn(_a)
		if c == nil {
			// 没有就创建
			c, err = s.dialTCPConn(_a)
			if err != nil {
				return err
			}
			// 启动处理协程
			s.w.Add(1)
			go s.handleTCPConnRoutine(c)
		}
		// 请求
		return s.doRequest(ctx, c, r, d, &s.tcp.at)
	}
	// udp
	if _a, ok := a.(*net.UDPAddr); ok {
		// 连接
		c := new(udpConn)
		c.conn = s.udp.c
		c.initAddr(_a)
		//
		return s.doRequest(ctx, c, r, d, &s.udp.at)
	}
	//
	return errUnknownAddress
}

// doRequest 封装 Request 的公共代码
func (s *Server) doRequest(ctx context.Context, c conn, r *Request, d any, at *gosync.Map[string, *activeTx]) error {
	// 事务
	t, err := s.newActiveTx(c, r.message, d, at)
	if err != nil {
		return err
	}
	r.tx = t
	// 日志
	s.Logger.DebugfTrace(t.TxKey(), "request to %s %s\n%s", c.Network(), c.RemoteAddrString(), t.writeData.String())
	// 立即发送
	err = c.write(t.writeData.Bytes())
	if err == nil {
		// 等待响应处理或底层超时
		select {
		case <-ctx.Done():
			err = ctx.Err()
			// 移除
			at.Del(t.key)
			// 通知
			t.Finish(err)
		case <-t.c:
			err = t.err
			// 要么是收到了响应的消息被调用
			// 要么是检查超时被被调用
			// 所以这里不需要调用
			// s.delActiveTx(t.key, at)
		}
	}
	// 日志
	s.Logger.DebugfTrace(t.TxKey(), "[%v] do request", time.Since(t.time))
	//
	return err
}
