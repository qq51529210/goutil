package sync

import (
	"sync"
	"sync/atomic"
)

// Chan 用于安全的并发写入和关闭
// 不至于 panic
type Chan[T any] struct {
	C  chan T
	m  sync.Mutex
	ok bool
}

// NewChan 返回新的 Chan
func NewChan[T any](len int) *Chan[T] {
	c := new(Chan[T])
	c.C = make(chan T, len)
	c.ok = true
	return c
}

// Close 关闭
func (s *Chan[T]) Close() {
	s.m.Lock()
	if !s.ok {
		s.m.Unlock()
		return
	}
	close(s.C)
	s.ok = false
	s.m.Unlock()
}

// Send 写入
func (s *Chan[T]) Send(v T) bool {
	s.m.Lock()
	// 已经关闭
	if !s.ok {
		s.m.Unlock()
		return false
	}
	// 写入
	select {
	case s.C <- v:
		s.m.Unlock()
		return true
	default:
		s.m.Unlock()
		return false
	}
}

// Signal 用于信号退出之类的
type Signal[T any] struct {
	// 信号
	C chan struct{}
	// 并发
	o int32
	// 通知的数据
	d T
}

// NewSignal 返回新的 Signal
func NewSignal[T any]() *Signal[T] {
	s := new(Signal[T])
	s.C = make(chan struct{})
	return s
}

// Close 关闭，第一次关闭返回 true
func (s *Signal[T]) Close(d T) bool {
	if atomic.CompareAndSwapInt32(&s.o, 0, 1) {
		s.d = d
		close(s.C)
		return true
	}
	return false
}

func (s *Signal[T]) Result() T {
	return s.d
}
