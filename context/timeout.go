package context

import (
	"context"
	"sync/atomic"
	"time"
)

// Timeout 超时上下文
type Timeout struct {
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
func (c *Timeout) Run(data any, timeout time.Duration) {
	t := time.Now().Add(timeout)
	c.timeout = timeout
	c.data = data
	c.deadline = &t
	c.timer = time.NewTimer(timeout)
	c.c = make(chan struct{})
	//
	go c.routine()
}

// Deadline 实现 context.Context
func (c *Timeout) Deadline() (time.Time, bool) {
	t := c.deadline
	return *t, true
}

// Err 实现 context.Context
func (c *Timeout) Err() error {
	return c.err
}

// Value 实现 context.Context
func (c *Timeout) Value(any) any {
	return c.data
}

// Done 实现 context.Context
func (c *Timeout) Done() <-chan struct{} {
	return c.c
}

// UpdateTime 更新超时时间
func (c *Timeout) UpdateTime() {
	t := time.Now().Add(c.timeout)
	c.deadline = &t
}

// Finish 结束通知
func (c *Timeout) Finish(err error) {
	if atomic.CompareAndSwapInt32(&c.ok, 0, 1) {
		c.err = err
		close(c.c)
	}
}

// routine 判断超时
func (c *Timeout) routine() {
	defer func() {
		// 计时器
		c.timer.Stop()
		// 回调
		if c.OnFinish != nil {
			c.OnFinish()
		}
	}()
	for {
		select {
		case <-c.c:
			return
		case now := <-c.timer.C:
			t := c.deadline
			// 超时了
			dur := t.Sub(now)
			if dur < 0 {
				c.Finish(context.DeadlineExceeded)
				return
			}
			c.timer.Reset(dur)
		}
	}
}
