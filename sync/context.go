package sync

import (
	"context"
	"sync/atomic"
	"time"
)

// Context 超时上下文
type Context struct {
	// 超时
	timeout time.Duration
	// 状态
	ok int32
	// 信号
	c chan struct{}
	// 错误
	err error
	// 用于保存发起请求时传入的数据
	data any
	// 用于判断超时清理
	deadline *time.Time
	// 超时计时器
	timer *time.Timer
	// 回调
	OnFinish func()
}

// Run 初始化
func (m *Context) Run(data any, timeout time.Duration) {
	t := time.Now().Add(timeout)
	m.timeout = timeout
	m.data = data
	m.deadline = &t
	m.timer = time.NewTimer(timeout)
	m.c = make(chan struct{})
	//
	go m.routine()
}

// Deadline 实现 context.Context
func (m *Context) Deadline() (time.Time, bool) {
	t := m.deadline
	return *t, true
}

// Err 实现 context.Context
func (m *Context) Err() error {
	return m.err
}

// Value 实现 context.Context
func (m *Context) Value(any) any {
	return m.data
}

// Done 实现 context.Context
func (m *Context) Done() <-chan struct{} {
	return m.c
}

// UpdateTime 更新超时时间
func (m *Context) UpdateTime() {
	t := time.Now().Add(m.timeout)
	m.deadline = &t
}

// Finish 结束通知
func (m *Context) Finish(err error) {
	if atomic.CompareAndSwapInt32(&m.ok, 0, 1) {
		m.err = err
		close(m.c)
	}
}

// routine 判断超时
func (m *Context) routine() {
	defer func() {
		// 计时器
		m.timer.Stop()
		// 回调
		if m.OnFinish != nil {
			m.OnFinish()
		}
	}()
	for {
		select {
		case <-m.c:
			return
		case now := <-m.timer.C:
			t := m.deadline
			// 超时了
			dur := now.Sub(*t)
			if dur < 0 {
				m.Finish(context.DeadlineExceeded)
				return
			}
			m.timer.Reset(dur)
		}
	}
}
