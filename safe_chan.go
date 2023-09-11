package util

import (
	"sync"
	"sync/atomic"
)

// SafeChan 用于安全的并发写入和关闭
// 不至于 panic
type SafeChan[T any] struct {
	C  chan T
	m  sync.Mutex
	ok bool
}

// NewSafeChan 返回新的 SafeChan
func NewSafeChan[T any](len int) *SafeChan[T] {
	c := new(SafeChan[T])
	c.C = make(chan T, len)
	c.ok = true
	return c
}

// Close 关闭
func (s *SafeChan[T]) Close() {
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
func (s *SafeChan[T]) Send(v T) bool {
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
type Signal struct {
	// 信号
	C chan struct{}
	o int32
}

// NewSignal 返回新的 Signal
func NewSignal() *Signal {
	s := new(Signal)
	s.C = make(chan struct{})
	return s
}

// Close 关闭
func (s *Signal) Close() {
	if atomic.CompareAndSwapInt32(&s.o, 0, 1) {
		close(s.C)
	}
}
