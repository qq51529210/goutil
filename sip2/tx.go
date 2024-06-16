package sip

import (
	"sync/atomic"
	"time"
)

type tx interface {
	// context.Context 接口
	Done() <-chan struct{}
	// context.Context 接口
	Err() error
	// context.Context 接口
	Deadline() (time.Time, bool)
	// 返回事务的标识
	ID() string
	// 完成，一般在处理响应中调用
	finish(err error)
}

type baseTx struct {
	id string
	// 状态
	ok int32
	// 信号
	done chan struct{}
	// 错误
	err error
	// 用于判断超时清理
	deadline time.Time
	// 创建时间
	createTime time.Time
}

// 实现 context.Context 接口
func (t *baseTx) Done() <-chan struct{} {
	return t.done
}

// 实现 context.Context 接口
func (t *baseTx) Err() error {
	return t.err
}

// 实现 context.Context 接口
func (t *baseTx) Deadline() (time.Time, bool) {
	return t.deadline, true
}

func (t *baseTx) ID() string {
	return t.id
}

func (t *baseTx) finish(err error) {
	if atomic.CompareAndSwapInt32(&t.ok, 0, 1) {
		t.err = err
		close(t.done)
	}
}
