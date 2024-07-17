package sip

import (
	"context"
	"fmt"
	gs "goutil/sync"
	"net"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type tcpServer struct {
	s *Server
	// 监听
	listener *net.TCPListener
	// 连接池
	conn gs.Map[connKey, *tcpConn]
	// 同步等待
	w sync.WaitGroup
	// 主动事务
	activeTx gs.Map[string, *tcpActiveTx]
	// 被动事务
	passiveTx gs.Map[string, *tcpPassiveTx]
	// 状态
	ok int32
}

func (s *tcpServer) isOK() bool {
	return atomic.LoadInt32(&s.ok) == 1
}

// Serve 监听 address 开始服务
func (s *tcpServer) Serve(address string) error {
	s.conn.Init()
	s.activeTx.Init()
	s.passiveTx.Init()
	// 地址
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	// 监听
	listener, err := net.ListenTCP(addr.Network(), addr)
	if err != nil {
		return err
	}
	s.listener = listener
	// 监听
	s.w.Add(1)
	go s.listenRoutine()
	// 检查
	s.w.Add(2)
	go s.checkActiveTxRoutine()
	go s.checkPassiveTxRoutine()
	// 日志
	s.s.logger.Infof("listen tcp %s", address)
	// 状态
	atomic.StoreInt32(&s.ok, 1)
	// 返回
	return nil
}

// listenRoutine 监听 tcp 连接，然后启动协程处理
func (s *tcpServer) listenRoutine() {
	defer func() {
		// 结束
		s.w.Done()
		// 异常
		if s.s.logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	for s.isOK() {
		// 接受
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			s.s.logger.Errorf("tcp accept %v", err)
			continue
		}
		// 开协程处理处理
		c := s.addConn(conn)
		s.w.Add(1)
		go s.handleConnRoutine(c)
	}
}

// addTCPConn 添加并返回
func (s *tcpServer) addConn(conn *net.TCPConn) *tcpConn {
	a := conn.RemoteAddr().(*net.TCPAddr)
	// 初始化
	c := new(tcpConn)
	c.conn = conn
	c.remoteIP = a.IP.String()
	c.remotePort = a.Port
	c.remoteAddr = fmt.Sprintf("%s:%d", c.remoteIP, c.remotePort)
	c.key.Init(a.IP, a.Port)
	// 添加
	s.conn.Set(c.key, c)
	//
	return c
}

// delConn 移除并关闭
func (s *tcpServer) delConn(conn *tcpConn) {
	s.conn.Del(conn.key)
	conn.conn.Close()
}

// getTCPConn 获取
func (s *tcpServer) getConn(addr *net.TCPAddr) *tcpConn {
	k := connKey{}
	k.Init(addr.IP, addr.Port)
	return s.conn.Get(k)
}

// dialConn 创建连接
func (s *tcpServer) dialConn(addr *net.TCPAddr) (*tcpConn, error) {
	conn, err := net.DialTimeout(addr.Network(), addr.String(), s.s.msgTimeout)
	if err != nil {
		return nil, err
	}
	return s.addConn(conn.(*net.TCPConn)), nil
}

// handleConnRoutine 处理 tcp conn 消息
func (s *tcpServer) handleConnRoutine(c *tcpConn) {
	defer func() {
		// 结束
		s.w.Done()
		// 移除
		s.delConn(c)
		// 异常
		s.s.logger.Recover(recover())
	}()
	r := newReader(c.conn, s.s.maxMessageLen)
	for s.isOK() {
		// 解析，错误直接返回关闭连接
		m := new(Message)
		if err := m.Dec(r, s.s.maxMessageLen); err != nil {
			s.s.logger.Errorf("tcp parse message %v", err)
			return
		}
		// 处理
		s.handleMsg(c, m)
	}
}

// handleMsg 处理 msg
func (s *tcpServer) handleMsg(conn *tcpConn, msg *Message) {
	method := strings.ToUpper(msg.Header.CSeq.Method)
	if msg.isReq {
		// 回调，没有注册不处理
		hf := s.s.handleFunc.reqFunc[method]
		if len(hf) > 0 {
			// 事务
			t := s.newPassiveTx(msg.TxKey())
			// 已经完成处理
			if atomic.LoadInt32(&t.ok) == 1 {
				return
			}
			// 没有完成，在协程中处理
			if atomic.CompareAndSwapInt32(&t.handing, 0, 1) {
				s.w.Add(1)
				go s.handleRequestRoutine(conn, t, msg, &reqFuncChain{f: hf})
			}
		}
		return
	}
	// 回调，没有注册不处理
	hf := s.s.handleFunc.resFunc[method]
	if len(hf) > 0 {
		// 响应消息
		if msg.StartLine[1][0] == '1' {
			// 1xx 消息没什么卵用，就不回调了
			return
		}
		// 事务，不一定有
		if t := s.deleteAndGetActiveTx(msg.TxKey()); t != nil {
			// 在协程中处理
			s.w.Add(1)
			go s.handleResponseRoutine(conn, t, msg, &resFuncChain{f: hf})
		}
	}
}

// handleRequestRoutine 在协程中处理请求消息
func (s *tcpServer) handleRequestRoutine(c *tcpConn, t *tcpPassiveTx, m *Message, f *reqFuncChain) {
	cost := time.Now()
	defer func() {
		// 结束
		s.w.Done()
		// 日志
		s.s.logger.DebugfTrace(t.id, "cost %v", time.Since(cost))
		// 异常
		s.s.logger.Recover(recover())
	}()
	// 日志
	s.s.logger.DebugfTrace(t.id, "request from tcp %s \n%v", c.remoteAddr, m)
	// 上下文
	var ctx Request
	ctx.tx = t
	ctx.Ser = s.s
	ctx.conn = c
	ctx.Message = m
	ctx.RemoteNetwork = networkTCP
	ctx.RemoteIP = c.remoteIP
	ctx.RemotePort = c.remotePort
	ctx.RemoteAddr = c.remoteAddr
	// 回调
	ctx.f = f
	f.handle(&ctx)
	// 没有完成，回复标记，等下一次的消息再回调
	if atomic.LoadInt32(&t.ok) == 0 {
		atomic.StoreInt32(&t.handing, 0)
	}
}

// handleResponseRoutine 在协程中处理响应消息
func (s *tcpServer) handleResponseRoutine(c *tcpConn, t *tcpActiveTx, m *Message, f *resFuncChain) {
	defer func() {
		// 结束
		s.w.Done()
		// 无论回调有没有通知，这里都通知一下
		t.finish(nil)
		// 异常
		s.s.logger.Recover(recover())
	}()
	// 日志
	s.s.logger.DebugfTrace(t.id, "response from udp %s \n%v", c.remoteAddr, m)
	// 上下文
	var ctx Response
	ctx.tx = t
	ctx.Ser = s.s
	ctx.conn = c
	ctx.Message = m
	ctx.ReqData = t.data
	ctx.RemoteNetwork = networkTCP
	ctx.RemoteIP = c.remoteIP
	ctx.RemotePort = c.remotePort
	ctx.RemoteAddr = c.remoteAddr
	// 回调
	ctx.f = f
	f.handle(&ctx)
}

// checkActiveTxRoutine 检查主动事务的超时
func (s *tcpServer) checkActiveTxRoutine() {
	// 计时器
	dur := s.s.checkTxDuration()
	timer := time.NewTimer(dur)
	defer func() {
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
		// 异常
		if s.s.logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	for s.isOK() {
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

// newActiveTx 添加并返回，用于主动发送请求
func (s *tcpServer) newActiveTx(id string, data any) (*tcpActiveTx, bool) {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	// 添加
	t, ok := s.activeTx.D[id]
	if t != nil {
		return t, ok
	}
	// 新的
	t = new(tcpActiveTx)
	t.id = id
	t.deadline = time.Now().Add(s.s.msgTimeout)
	t.done = make(chan struct{})
	t.data = data
	//
	s.activeTx.D[t.id] = t
	//
	return t, ok
}

// deleteAndGetActiveTx 看名称
func (s *tcpServer) deleteAndGetActiveTx(id string) *tcpActiveTx {
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

// deleteActiveTx 看名称
func (s *tcpServer) deleteActiveTx(t *tcpActiveTx, err error) {
	t.finish(err)
	s.activeTx.Del(t.id)
}

// checkPassiveTxRoutine 检查被动事务的超时
func (s *tcpServer) checkPassiveTxRoutine() {
	// 计时器
	dur := s.s.checkTxDuration()
	timer := time.NewTimer(dur)
	defer func() {
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
		// 异常
		if s.s.logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	// 开始
	for s.isOK() {
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

// newPassiveTx 添加并返回，用于被动接收请求
func (s *tcpServer) newPassiveTx(id string) *tcpPassiveTx {
	// 锁
	s.passiveTx.Lock()
	defer s.passiveTx.Unlock()
	//
	t := s.passiveTx.D[id]
	if t == nil {
		t = new(tcpPassiveTx)
		t.id = id
		t.deadline = time.Now().Add(s.s.msgTimeout)
		t.done = make(chan struct{})
		//
		s.passiveTx.D[id] = t
	}
	//
	return t
}

func (s *tcpServer) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.ok, 1, -1) {
		// 关闭 conn
		s.listener.Close()
		// 事务通知
		s.shutdownActiveTx()
		s.shutdownPassiveTx()
		// 关闭连接
		s.shutdownConn()
		// 等待所有协程退出
		s.w.Wait()
	}
}

// shutdownConn 关闭所有连接
func (s *tcpServer) shutdownConn() {
	// 锁
	s.conn.Lock()
	defer s.conn.Unlock()
	//
	for _, d := range s.conn.D {
		d.conn.Close()
	}
	s.conn.D = make(map[connKey]*tcpConn)
}

// shutdownPassiveTx 通知所有的主动事务，服务关闭了
func (s *tcpServer) shutdownActiveTx() {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	//
	for _, d := range s.activeTx.D {
		d.finish(ErrServerShutdown)
	}
	s.activeTx.D = make(map[string]*tcpActiveTx)
}

// shutdownPassiveTx 通知所有的被动事务，服务关闭了
func (s *tcpServer) shutdownPassiveTx() {
	// 锁
	s.passiveTx.Lock()
	defer s.passiveTx.Unlock()
	//
	for _, d := range s.passiveTx.D {
		d.finish(ErrServerShutdown)
	}
	s.passiveTx.D = make(map[string]*tcpPassiveTx)
}

// Request 发送请求
func (s *tcpServer) Request(ctx context.Context, msg *Message, addr *net.TCPAddr, data any) error {
	cost := time.Now()
	// 连接
	conn := s.getConn(addr)
	if conn == nil {
		// 没有就创建
		c, err := s.dialConn(addr)
		if err != nil {
			return err
		}
		// 启动处理协程
		s.w.Add(1)
		go s.handleConnRoutine(c)
	}
	// 事务
	t, ok := s.newActiveTx(msg.TxKey(), data)
	// 第一次
	if !ok {
		if err := t.writeMsg(conn, msg); err != nil {
			s.deleteActiveTx(t, err)
			return err
		}
	}
	// 日志
	s.s.logger.DebugfTrace(t.id, "request to tcp %s \n%v", conn.remoteAddr, msg)
	// 等待
	var err error
	select {
	case <-ctx.Done():
		// 传入的上下文
		err = ctx.Err()
	case <-t.Done():
		// 底层超时，或者 RequestAbort
		err = t.Err()
	}
	// 日志
	s.s.logger.DebugfTrace(t.id, "cost %v", time.Since(cost))
	// 移除
	s.deleteActiveTx(t, err)
	if err == ErrFinish {
		return nil
	}
	return err
}
