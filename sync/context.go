package sync

import (
	"context"
	"sync/atomic"
	"time"
)

type TimeoutContextPool struct {
	// 是否启动协程
	r int32
	// 数据
	d Map[string, *TimeoutContext]
}

// NewTimeoutContextPool 返回，然后做了初始化
func NewTimeoutContextPool() *TimeoutContextPool {
	p := new(TimeoutContextPool)
	p.d.Init()
	go p.routine()
	return p
}

func (txp *TimeoutContextPool) routine() {
	for {
		now := time.Now()
		// 所有
		txs := txp.d.Values()
		for _, tx := range txs {
			// 结束了，移除
			if atomic.LoadInt32(&tx.ok) == 1 {
				txp.d.Del(tx.id)
				continue
			}
			// 检查时间
			tt := tx.deadline
			if now.After(*tt) {
				// 超时通知
				tx.Finish(nil, context.DeadlineExceeded)
				// 移除
				txp.d.Del(tx.id)
				continue
			}
		}
		time.Sleep(time.Second)
	}
}

// New 获取/创建，返回 true 表示新的
// name 唯一标识，data 携带的数据，dur 超时时间
func (txp *TimeoutContextPool) New(id string, data any, dur time.Duration) (*TimeoutContext, bool) {
	// 协程懒启动
	if atomic.CompareAndSwapInt32(&txp.r, 0, 1) {
		go txp.routine()
	}
	// 锁
	txp.d.Lock()
	defer txp.d.Unlock()
	//
	tx := txp.d.D[id]
	if tx == nil {
		tx = new(TimeoutContext)
		tx.id = id
		tx.data = data
		tx.dur = dur
		tx.UpdateDeadlineTime()
		tx.done = make(chan struct{})
		txp.d.D[id] = tx
		return tx, true
	}
	return tx, false
}

// Get 返回
func (txp *TimeoutContextPool) Get(id string) *TimeoutContext {
	return txp.d.Get(id)
}

// TimeoutContext 超时上下文
type TimeoutContext struct {
	// 标识
	id string
	// GetTimeoutContext 传入的数据
	data any
	// 结果数据
	result any
	// 错误
	err error
	// 超时时长
	dur time.Duration
	// 超时时间点
	deadline *time.Time
	// 状态
	ok int32
	// 信号
	done chan struct{}
}

// Deadline 实现 context.Context
func (tc *TimeoutContext) Deadline() (time.Time, bool) {
	t := tc.deadline
	return *t, true
}

// Err 实现 context.Context
func (tc *TimeoutContext) Err() error {
	return tc.err
}

// Value 实现 context.Context
// 返回传入的 data
func (tc *TimeoutContext) Value(any) any {
	return tc.data
}

// Done 实现 context.Context
func (tc *TimeoutContext) Done() <-chan struct{} {
	return tc.done
}

// UpdateTime 更新超时时间
func (tc *TimeoutContext) UpdateDeadlineTime() {
	t := time.Now().Add(tc.dur)
	tc.deadline = &t
}

// Finish 结束通知
func (tc *TimeoutContext) Finish(result any, err error) {
	if atomic.CompareAndSwapInt32(&tc.ok, 0, 1) {
		tc.result = result
		tc.err = err
		close(tc.done)
	}
}

// Result 返回 Finish 传入的结果和错误
func (tc *TimeoutContext) Result() (any, error) {
	return tc.result, tc.err
}
