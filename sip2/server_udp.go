package sip

import (
	"bytes"
	"context"
	"fmt"
	gs "goutil/sync"
	"io"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// udpServer udp 服务
type udpServer struct {
	// 引用
	s *Server
	// 消息超时重传
	minRTO, maxRTO time.Duration
	// 底层连接
	conn *net.UDPConn
	// 同步等待
	w sync.WaitGroup
	// 主动事务
	activeTx gs.Map[string, *udpActiveTx]
	// 被动事务
	passiveTx gs.Map[string, *udpPassiveTx]
	// 状态
	ok int32
}

// Serve 监听 address 开始服务
func (s *udpServer) Serve(address string) error {
	s.activeTx.Init()
	s.passiveTx.Init()
	// 地址
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}
	// 连接
	conn, err := net.ListenUDP(addr.Network(), addr)
	if err != nil {
		return err
	}
	s.conn = conn
	// 读数据
	n := runtime.NumCPU()
	s.w.Add(n)
	for i := 0; i < n; i++ {
		go s.readRoutine()
	}
	// 事务检查
	s.w.Add(2)
	go s.checkActiveTxRoutine()
	go s.checkPassiveTxRoutine()
	// 超时重发检查
	s.w.Add(1)
	go s.checkRTORoutine()
	// 日志
	s.s.Logger.Infof("listen udp %s, read routine %d", address, n)
	// 状态
	atomic.StoreInt32(&s.ok, 1)
	//
	return nil
}

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

// initConn 初始化 c 的字段
func (s *udpServer) initConn(c *udpConn, a *net.UDPAddr) {
	c.conn = s.conn
	c.addr = a
	c.remoteIP = a.IP.String()
	c.remotePort = a.Port
	c.remoteAddr = fmt.Sprintf("%s:%d", c.remoteIP, c.remotePort)
}

// readUDPRoutine 读取 udp 数据，解析成 Message ，然后处理
func (s *udpServer) readRoutine() {
	// 清理
	defer func() {
		// 结束
		s.w.Done()
		// 异常
		if s.s.Logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	var err error
	d := &udpData{b: make([]byte, s.s.MaxMessageLen)}
	r := &reader{r: d}
	for atomic.LoadInt32(&s.ok) == 0 {
		// 读取 udp 数据
		d.n, d.a, err = s.conn.ReadFromUDP(d.b)
		if err != nil {
			s.s.Logger.Errorf("udp read data %v", err)
			continue
		}
		// 初始化准备解析
		d.i = 0
		r.begin = 0
		r.end = 0
		r.parsed = 0
		r.buf = d.b
		// 连接
		c := &udpConn{}
		s.initConn(c, d.a)
		// 一个数据包可能有多个消息，这里需要循环解析处理
		for atomic.LoadInt32(&s.ok) == 0 {
			// 解析
			m := new(Message)
			if err = m.Dec(r, s.s.MaxMessageLen); err != nil {
				if err != io.EOF {
					s.s.Logger.Errorf("udp parse message %v", err)
				}
				break
			}
			// 处理
			s.handleMsg(c, m)
		}
	}
}

func (s *udpServer) handleMsg(conn *udpConn, msg *Message) {
	method := strings.ToUpper(msg.Header.CSeq.Method)
	// 请求消息
	if msg.isReq {
		// 回调
		hf, ok := s.s.handleFunc.reqFunc[method]
		if !ok {
			// 不支持的方法，而且没有注册 RequestNotFoundFunc ，这里直接回复
			if s.s.handleFunc.reqNotFoundFunc == nil {
				s.s.handleRequestNotFound(conn, msg)
				return
			}
			// 回调
			if res := s.s.handleFunc.reqNotFoundFunc(msg); res != nil {
				if err := conn.writeMsg(res); err != nil {
					s.s.Logger.Errorf("write udp data %v", err)
				}
			}
			return
		}
		// 事务
		t := s.newPassiveTx(conn, msg.txKey())
		// 已经完成处理
		if atomic.LoadInt32(&t.ok) == 1 {
			// 直接发送响应数据缓存，无需回调
			if err := t.rewrite(); err != nil {
				s.s.Logger.ErrorfTrace(t.id, "rewrite response %v", err)
			}
			return
		}
		// 没有完成，在协程中处理
		if atomic.CompareAndSwapInt32(&t.handing, 0, 1) {
			s.w.Add(1)
			go s.handleRequestRoutine(t, msg, hf)
		}
		return
	}
	// 回调
	hf, ok := s.s.handleFunc.resFunc[method]
	if !ok {
		// 没有注册响应处理
		return
	}
	// 响应消息
	if msg.StartLine[1][0] == '1' {
		// 停止超时重传
		if t := s.activeTx.Get(msg.txKey()); t != nil {
			t.rtoStop = true
		}
		// 1xx 消息没什么卵用，就不回调了
		return
	}
	// 事务，不一定有
	if t := s.deleteAndGetActiveTx(msg.txKey()); t != nil {
		// 在协程中处理
		s.w.Add(1)
		go s.handleResponseRoutine(t, msg, hf)
	}
}

// handleRequestRoutine 在协程中处理请求消息
func (s *udpServer) handleRequestRoutine(t *udpPassiveTx, m *Message, f []HandleRequestFunc) {
	defer func() {
		// 异常
		s.s.Logger.Recover(recover())
		// 结束
		s.w.Done()
	}()
	// 上下文
	var req Request
	req.tx = t
	req.Message = m
	req.Server = s.s
	req.RemoteNetwork = t.conn.Network()
	req.RemoteIP = t.conn.remoteIP
	req.RemotePort = t.conn.remotePort
	req.RemoteAddr = t.conn.remoteAddr
	// 回调
	req.handleFunc = f
	req.callback()
	// 没有完成，回复标记，等下一次的消息再回调
	if atomic.LoadInt32(&t.ok) == 0 {
		atomic.StoreInt32(&t.handing, 0)
	}
}

// handleResponseRoutine 在协程中处理响应消息
func (s *udpServer) handleResponseRoutine(t *udpActiveTx, m *Message, f []HandleResponseFunc) {
	defer func() {
		// 异常
		s.s.Logger.Recover(recover())
		// 无论回调有没有通知，这里都通知一下
		t.finish(nil)
		// 结束
		s.w.Done()
	}()
	// 上下文
	var res Response
	res._Context.tx = t
	res._Context.Message = m
	res._Context.Server = s.s
	// 回调
	res.handleFunc = f
	res.callback()
}

// checkRTORoutine 检测消息超时重传
func (s *udpServer) checkRTORoutine() {
	// 计时器
	timer := time.NewTimer(s.minRTO)
	defer func() {
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
		// 异常
		if s.s.Logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	wg := new(sync.WaitGroup)
	for atomic.LoadInt32(&s.ok) == 0 {
		// 时间到
		now := <-timer.C
		// 副本
		ts := s.activeTx.Values()
		// 并发计算
		n := runtime.NumCPU()
		for len(ts) > n {
			m := len(ts) / n
			wg.Add(1)
			go s.rtoRoutine(wg, ts[:m], now)
			ts = ts[m:]
		}
		if len(ts) > 0 {
			wg.Add(1)
			go s.rtoRoutine(wg, ts, now)
		}
		// 等待并发结束
		wg.Wait()
		// 重置计时器
		timer.Reset(s.minRTO)
	}
}

// rtoRoutine 发送 udp 数据
func (s *udpServer) rtoRoutine(wg *sync.WaitGroup, ts []*udpActiveTx, now time.Time) {
	defer func() {
		// 结束
		wg.Done()
		// 异常
		s.s.Logger.Recover(recover())
	}()
	// 循环检查，然后发送，超时移除
	for _, t := range ts {
		// 停止发送，在收到 1xx 后设置
		if t.rtoStop {
			continue
		}
		// 是否应该超时重传
		if now.Sub(t.rtoTime) >= t.rto {
			d := t.rtoData
			if err := t.conn.write(d.Bytes()); err != nil {
				s.s.Logger.Errorf("rto %v", err)
				continue
			}
			// 保存发送时间
			t.rtoTime = now
			// 如果没有到达最大值
			if t.rto < s.maxRTO {
				// 翻倍
				t.rto *= 2
				// 但是不能比最大值
				if t.rto > s.maxRTO {
					t.rto = s.maxRTO
				}
			}
		}
	}
}

// checkActiveTxRoutine 检查主动事务的超时
func (s *udpServer) checkActiveTxRoutine() {
	// 计时器
	dur := s.s.checkTxDuration()
	timer := time.NewTimer(dur)
	defer func() {
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
		// 异常
		if s.s.Logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	for atomic.LoadInt32(&s.ok) == 0 {
		// 时间
		now := <-timer.C
		// 组装
		ts := s.activeTx.Values()
		// 检查
		for _, t := range ts {
			// 超时
			if now.After(t.deadline) {
				// 移除
				s.activeTx.Del(t.id)
				// 通知
				t.finish(context.DeadlineExceeded)
			}
		}
		// 重置计时器
		timer.Reset(dur)
	}
}

// udpActiveTx 主动发起请求的事务
type udpActiveTx struct {
	baseTx
	// 连接
	conn *udpConn
	// 请求的数据
	data any
	// 消息重发的间隔，发送一次叠加一倍
	rto time.Duration
	// 发送的数据
	rtoData *bytes.Buffer
	// 发送数据的时间
	rtoTime time.Time
	// 停止 rto
	rtoStop bool
}

// newActiveTx 添加并返回，用于主动发送请求
func (s *udpServer) newActiveTx(id string, conn *udpConn, data any) (*udpActiveTx, bool) {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	// 添加
	s.activeTx.Lock()
	t, ok := s.activeTx.D[id]
	if t != nil {
		return t, ok
	}
	t = new(udpActiveTx)
	t.id = id
	t.createTime = time.Now()
	t.deadline = t.createTime.Add(s.s.MsgTimeout)
	t.done = make(chan struct{})
	t.data = data
	t.conn = conn
	t.rto = s.minRTO
	t.rtoTime = t.createTime
	t.rtoData = bytes.NewBuffer(nil)
	//
	s.activeTx.D[t.id] = t
	//
	return t, ok
}

func (s *udpServer) deleteAndGetActiveTx(id string) *udpActiveTx {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	//
	t := s.activeTx.D[id]
	if t != nil {
		delete(s.activeTx.D, id)
	}
	//
	return t
}

func (s *udpServer) deleteActiveTx(t *udpActiveTx, err error) {
	t.finish(err)
	s.activeTx.Del(t.id)
}

// checkPassiveTxRoutine 检查被动事务的超时
func (s *udpServer) checkPassiveTxRoutine() {
	// 计时器
	dur := s.s.checkTxDuration()
	timer := time.NewTimer(dur)
	defer func() {
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
		// 异常
		if s.s.Logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	for atomic.LoadInt32(&s.ok) == 0 {
		// 时间
		now := <-timer.C
		// 组装
		ts := s.passiveTx.Values()
		// 检查
		for _, t := range ts {
			// 超时
			if now.After(t.deadline) {
				// 移除
				s.passiveTx.Del(t.id)
				// 通知
				t.finish(context.DeadlineExceeded)
			}
		}
		// 重置计时器
		timer.Reset(dur)
	}
}

// udpPassiveTx 被动接受请求的事务
type udpPassiveTx struct {
	baseTx
	// 连接
	conn *udpConn
	// 用于控制多消息并发时的单一处理
	handing int32
	// 响应的数据缓存
	dataBuff *bytes.Buffer
}

func (t *udpPassiveTx) rewrite() error {
	td := t.dataBuff
	d := td.Bytes()
	if len(d) > 0 {
		return t.conn.write(td.Bytes())
	}
	return nil
}

// newPassiveTx 添加并返回，用于被动接收请求
func (s *udpServer) newPassiveTx(conn *udpConn, id string) *udpPassiveTx {
	// 锁
	s.passiveTx.Lock()
	defer s.passiveTx.Unlock()
	//
	t := s.passiveTx.D[id]
	if t == nil {
		t = new(udpPassiveTx)
		t.id = id
		t.createTime = time.Now()
		t.deadline = t.createTime.Add(s.s.MsgTimeout)
		t.done = make(chan struct{})
		t.conn = conn
		t.dataBuff = bytes.NewBuffer(nil)
		//
		s.passiveTx.D[id] = t
	}
	//
	return t
}

// Shutdown 停止服务
func (s *udpServer) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.ok, 0, 1) {
		// 关闭 conn
		s.conn.Close()
		// 事务通知
		s.shutdownActiveTx()
		s.shutdownPassiveTx()
		// 等待所有协程退出
		s.w.Wait()
	}
}

// shutdownPassiveTx 通知所有的主动事务，服务关闭了
func (s *udpServer) shutdownActiveTx() {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	//
	for _, d := range s.activeTx.D {
		d.finish(ErrServerShutdown)
	}
	s.passiveTx.D = make(map[string]*udpPassiveTx)
}

// shutdownPassiveTx 通知所有的被动事务，服务关闭了
func (s *udpServer) shutdownPassiveTx() {
	// 锁
	s.passiveTx.Lock()
	defer s.passiveTx.Unlock()
	//
	for _, d := range s.passiveTx.D {
		d.finish(ErrServerShutdown)
	}
	s.passiveTx.D = make(map[string]*udpPassiveTx)
}

// Request 发送请求
func (s *udpServer) Request(ctx context.Context, msg *Message, addr *net.UDPAddr, data any) error {
	// 连接
	conn := &udpConn{}
	s.initConn(conn, addr)
	// 事务
	t, ok := s.newActiveTx(msg.txKey(), conn, data)
	// 第一次
	if !ok {
		// 格式化消息
		msg.Enc(t.rtoData)
		// 立即发送一次
		d := t.rtoData
		if err := conn.write(d.Bytes()); err != nil {
			s.deleteActiveTx(t, err)
			return err
		}
	}
	// 等待
	var err error
	select {
	case <-ctx.Done():
		// 传入的上下文
		err = ctx.Err()
	case <-t.Done():
		// 底层超时
		err = t.Err()
	}
	// 移除
	s.deleteActiveTx(t, err)
	if err == ErrFinish {
		return nil
	}
	return err
}
