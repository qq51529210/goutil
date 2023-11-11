package sip

import (
	"bytes"
	"context"
	"fmt"
	"goutil/log"
	gosync "goutil/sync"
	"time"
)

// baseTx 实现一个 context.Context
type baseTx struct {
	// 池的 key
	key string
	// 用于判断超时清理
	deadline time.Time
	// 信号
	signal *gosync.Signal
	// 错误
	err error
	// 用于保存发起请求时传入的数据
	data any
	// 创建时间
	time time.Time
}

func (m *baseTx) Deadline() (time.Time, bool) {
	return m.deadline, true
}

func (m *baseTx) Err() error {
	return m.err
}

func (m *baseTx) Value(any) any {
	return m.data
}

func (m *baseTx) Done() <-chan struct{} {
	return m.signal.C
}

func (m *baseTx) TxKey() string {
	return m.key
}

func (m *baseTx) Time() time.Time {
	return m.time
}

// Finish 异步通知，用于在处理响应的时候，通知发送请求的那个协程
// 底层的超时通知是 context.DeadlineExceeded
func (m *baseTx) Finish(err error) {
	if m.signal.Close() {
		m.err = err
	}
}

// activeTx 用于主动发起请求
type activeTx struct {
	baseTx
	// 使用的连接
	conn conn
	// 用于 udp 消息重发间隔，每发送一次叠加一倍，但是有最大值
	rto time.Duration
	// 发送时间，用于 udp 消息重发计算
	writeTime time.Time
	// 用于发送数据，用于 udp 消息重发
	writeData bytes.Buffer
}

// newActiveTx 添加并返回，用于主动发送请求
func (s *Server) newActiveTx(c conn, m *message, d any, at *gosync.Map[string, *activeTx]) (*activeTx, error) {
	//
	t := new(activeTx)
	t.key = m.txKey()
	t.time = time.Now()
	t.deadline = t.time.Add(s.TxTimeout)
	t.signal = gosync.NewSignal()
	t.data = d
	t.conn = c
	t.rto = s.RTO
	t.writeTime = t.time
	m.Enc(&t.writeData)
	// 添加
	at.Lock()
	_, ok := at.D[t.key]
	if ok {
		t.signal.Close()
		return nil, errTransactionExists
	}
	at.Unlock()
	//
	log.Debugf("%s new active tx", t.key)
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
func (s *Server) checkActiveTxTimeoutRoutine(name string, at *gosync.Map[string, *activeTx]) {
	logTrace := fmt.Sprintf("%s check active tx routine", name)
	// 计时器
	dur := s.TxTimeout / 2
	timer := time.NewTimer(dur)
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
				//
				log.DebugfTrace(t.key, "active tx timeout cost %v", time.Since(t.time))
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

// newPassiveTx 添加并返回，用于被动接收请求
func (s *Server) newPassiveTx(m *message, pt *gosync.Map[string, *passiveTx]) *passiveTx {
	k := m.txKey()
	// 添加
	pt.Lock()
	t := pt.D[k]
	if t == nil {
		t = new(passiveTx)
		t.key = k
		t.time = time.Now()
		t.deadline = t.time.Add(s.TxTimeout)
		t.signal = gosync.NewSignal()
		pt.D[k] = t
		//
		log.DebugfTrace(k, "new passive tx")
	}
	pt.Unlock()
	//
	return t
}

// checkPassiveTxTimeoutRoutine 检查被动事务的超时
func (s *Server) checkPassiveTxTimeoutRoutine(name string, pt *gosync.Map[string, *passiveTx]) {
	logTrace := fmt.Sprintf("%s check passive tx routine", name)
	// 计时器
	dur := s.TxTimeout / 2
	timer := time.NewTimer(dur)
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
				//
				log.DebugfTrace(t.key, "passive tx timeout cost %v", time.Since(t.time))
			}
		}
		// 重置计时器
		timer.Reset(dur)
	}
}
