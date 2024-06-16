package sip

import (
	"context"
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
	go s.checkActiveTxRoutine()
	go s.checkPassiveTxRoutine()
	// 日志
	s.s.Logger.Infof("listen tcp %s", address)
	// 返回
	return nil
}

// listenRoutine 监听 tcp 连接，然后启动协程处理
func (s *tcpServer) listenRoutine() {
	defer func() {
		// 结束
		s.w.Done()
		// 异常
		if s.s.Logger.Recover(recover()) {
			os.Exit(1)
		}
	}()
	for atomic.LoadInt32(&s.ok) == 0 {
		// 接受
		conn, err := s.listener.AcceptTCP()
		if err != nil {
			s.s.Logger.Errorf("tcp accept %v", err)
			continue
		}
		// 处理
		c := s.addConn(conn)
		s.w.Add(1)
		go s.handleConnRoutine(c)
	}
}

// addTCPConn 添加并返回
func (s *tcpServer) addConn(conn *net.TCPConn) *tcpConn {
	// 初始化
	c := new(tcpConn)
	c.init(conn)
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
	conn, err := net.DialTimeout(addr.Network(), addr.String(), s.s.MsgTimeout)
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
		s.s.Logger.Recover(recover())
	}()
	r := newReader(c.conn, s.s.MaxMessageLen)
	for atomic.LoadInt32(&s.ok) == 0 {
		// 解析，错误直接返回关闭连接
		m := new(Message)
		if err := m.Dec(r, s.s.MaxMessageLen); err != nil {
			s.s.Logger.Errorf("tcp parse message %v", err)
			return
		}
		// 处理
		s.handleMsg(c, m)
	}
}

func (s *tcpServer) handleMsg(conn *tcpConn, msg *Message) {
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
		// 没有完成，在协程中处理
		if atomic.LoadInt32(&t.ok) != 1 && atomic.CompareAndSwapInt32(&t.handing, 0, 1) {
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
func (s *tcpServer) handleRequestRoutine(t *tcpPassiveTx, m *Message, f []HandleRequestFunc) {
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
func (s *tcpServer) handleResponseRoutine(t *tcpActiveTx, m *Message, f []HandleResponseFunc) {
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

// tcpActiveTx 主动发起请求的事务
type tcpActiveTx struct {
	baseTx
	// 请求的数据
	data any
}

// newActiveTx 添加并返回，用于主动发送请求
func (s *tcpServer) newActiveTx(id string, data any) (*tcpActiveTx, bool) {
	// 锁
	s.activeTx.Lock()
	defer s.activeTx.Unlock()
	// 添加
	s.activeTx.Lock()
	t, ok := s.activeTx.D[id]
	if t != nil {
		return t, ok
	}
	t = new(tcpActiveTx)
	t.id = id
	t.createTime = time.Now()
	t.deadline = t.createTime.Add(s.s.MsgTimeout)
	t.done = make(chan struct{})
	t.data = data
	//
	s.activeTx.D[t.id] = t
	//
	return t, ok
}

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

// tcpPassiveTx 被动接受请求的事务
type tcpPassiveTx struct {
	baseTx
	// 连接
	conn *tcpConn
	// 用于控制多消息并发时的单一处理
	handing int32
}

// newPassiveTx 添加并返回，用于被动接收请求
func (s *tcpServer) newPassiveTx(conn *tcpConn, id string) *tcpPassiveTx {
	// 锁
	s.passiveTx.Lock()
	defer s.passiveTx.Unlock()
	//
	t := s.passiveTx.D[id]
	if t == nil {
		t = new(tcpPassiveTx)
		t.id = id
		t.createTime = time.Now()
		t.deadline = t.createTime.Add(s.s.MsgTimeout)
		t.done = make(chan struct{})
		t.conn = conn
		//
		s.passiveTx.D[id] = t
	}
	//
	return t
}

func (s *tcpServer) Shutdown() {
	if atomic.CompareAndSwapInt32(&s.ok, 0, 1) {
		// 关闭 conn
		s.listener.Close()
		// 事务通知
		s.shutdownActiveTx()
		s.shutdownPassiveTx()
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
	s.passiveTx.D = make(map[string]*tcpPassiveTx)
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
	t, ok := s.newActiveTx(msg.txKey(), data)
	// 第一次
	if !ok {
		// 发送
		if err := conn.writeMsg(msg); err != nil {
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
