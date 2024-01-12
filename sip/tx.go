package sip

import (
	"bytes"
	"context"
	gosync "goutil/sync"
	"sync/atomic"
	"time"
)

type tx interface {
	context.Context
	TxKey() string
	Finish(err error)
	dataBuffer() *bytes.Buffer
	conn() conn
}

// baseTx 实现一个 context.Context
type baseTx struct {
	// 池的 key
	key string
	// 状态
	ok int32
	// 信号
	exit chan struct{}
	// 错误
	err error
	// 用于判断超时清理
	deadline time.Time
	// 创建时间
	time time.Time
	// 使用的连接
	c conn
}

func (m *baseTx) Deadline() (time.Time, bool) {
	return m.deadline, true
}

func (m *baseTx) Err() error {
	return m.err
}

func (m *baseTx) Done() <-chan struct{} {
	return m.exit
}

func (m *baseTx) TxKey() string {
	return m.key
}

func (m *baseTx) conn() conn {
	return m.c
}

// Finish 异步通知，用于在处理响应的时候，通知发送请求的那个协程
// 底层的超时通知是 context.DeadlineExceeded
// 不要保存在其他协程作为 context.Context
// 因为 Err() 可能返回 nil
func (m *baseTx) Finish(err error) {
	if atomic.CompareAndSwapInt32(&m.ok, 0, 1) {
		m.err = err
		close(m.exit)
	}
}

// activeTx 用于主动发起请求
type activeTx struct {
	baseTx
	// 用于保存发起请求时传入的数据
	data any
	// 用于 udp 消息重发间隔，每发送一次叠加一倍，但是有最大值
	rto time.Duration
	// 发送时间，用于 udp 消息重发计算
	writeTime time.Time
	// 用于发送数据，用于 udp 消息重发
	writeData bytes.Buffer
	// 停止发送，一般是遇到 1xx 类响应
	stopRT bool
}

func (m *activeTx) dataBuffer() *bytes.Buffer {
	return &m.writeData
}

func (m *activeTx) Value(any) any {
	return m.data
}

// newActiveTx 添加并返回，用于主动发送请求
func (s *Server) newActiveTx(c conn, m *message, d any, at *gosync.Map[string, *activeTx]) (*activeTx, error) {
	//
	t := new(activeTx)
	t.key = m.txKey()
	t.time = time.Now()
	t.deadline = t.time.Add(s.TxTimeout)
	t.data = d
	t.c = c
	t.rto = s.MinRTO
	t.writeTime = t.time
	m.Enc(&t.writeData)
	// 添加
	at.Lock()
	_, ok := at.D[t.key]
	if ok {
		// 已存在
		return nil, errTransactionExists
	}
	t.exit = make(chan struct{})
	at.D[t.key] = t
	at.Unlock()
	//
	return t, nil
}

// delActiveTx 移除
func (s *Server) delActiveTx(k string, at *gosync.Map[string, *activeTx]) *activeTx {
	at.Lock()
	t := at.D[k]
	if t != nil {
		delete(at.D, k)
	}
	at.Unlock()
	//
	return t
}

// checkActiveTxTimeoutRoutine 检查主动事务的超时
func (s *Server) checkActiveTxTimeoutRoutine(network string, at *gosync.Map[string, *activeTx]) {
	// 计时器
	dur := s.TxTimeout / 2
	timer := time.NewTimer(dur)
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.Warnf("%s check active tx routine stop", network)
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
	}()
	// 日志
	s.Logger.Debugf("%s check active tx routine start", network)
	// 开始
	var ts []*activeTx
	for s.isOK() {
		// 时间
		now := <-timer.C
		// 组装
		ts = ts[:0]
		at.RLock()
		for _, d := range at.D {
			ts = append(ts, d)
		}
		at.RUnlock()
		// 检查
		for _, t := range ts {
			// 超时
			if now.After(t.deadline) {
				// 移除
				at.Del(t.key)
				// 通知
				t.Finish(context.DeadlineExceeded)
			}
		}
		// 重置计时器
		timer.Reset(dur)
	}
}

// passiveTx 用于被动接收请求
type passiveTx struct {
	baseTx
	// 用于发送数据
	writeData bytes.Buffer
	// 用于控制多消息并发时的单一处理
	handing int32
	// 用于判断是否处理完毕
	done bool
}

func (m *passiveTx) dataBuffer() *bytes.Buffer {
	return &m.writeData
}

func (m *passiveTx) Value(any) any {
	return nil
}

// newPassiveTx 添加并返回，用于被动接收请求
func (s *Server) newPassiveTx(c conn, m *message, pt *gosync.Map[string, *passiveTx]) *passiveTx {
	k := m.txKey()
	// 添加
	pt.Lock()
	t := pt.D[k]
	if t == nil {
		t = new(passiveTx)
		t.key = k
		t.time = time.Now()
		t.c = c
		t.deadline = t.time.Add(s.TxTimeout)
		t.exit = make(chan struct{})
		pt.D[k] = t
	}
	pt.Unlock()
	//
	return t
}

// checkPassiveTxTimeoutRoutine 检查被动事务的超时
func (s *Server) checkPassiveTxTimeoutRoutine(network string, pt *gosync.Map[string, *passiveTx]) {
	// 计时器
	dur := s.TxTimeout / 2
	timer := time.NewTimer(dur)
	defer func() {
		// 异常
		s.Logger.Recover(recover())
		// 日志
		s.Logger.Warnf("%s check passive tx routine stop", network)
		// 计时器
		timer.Stop()
		// 结束
		s.w.Done()
	}()
	// 日志
	s.Logger.Debugf("%s check passive tx routine start", network)
	// 开始
	var ts []*passiveTx
	for s.isOK() {
		// 时间
		now := <-timer.C
		// 组装
		ts = ts[:0]
		pt.RLock()
		for _, d := range pt.D {
			ts = append(ts, d)
		}
		pt.RUnlock()
		// 检查
		for _, t := range ts {
			// 超时
			if now.After(t.deadline) {
				// 移除
				pt.Del(t.key)
				// 通知
				t.Finish(context.DeadlineExceeded)
			}
		}
		// 重置计时器
		timer.Reset(dur)
	}
}
