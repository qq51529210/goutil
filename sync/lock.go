package sync

import (
	"sync/atomic"
	"time"
)

// 接口
type LockHandler interface {
	// 抢锁
	Lock() (bool, error)
	// 是否到时间启用回调
	LockInterval() time.Duration
	// 是否到时间启用回调
	HandleInterval() time.Duration
	// 是否启用
	IsEnable() bool
	// 回调
	Handle()
	// 回调
	OnError(error)
}

// Locker 用于控制在部署多个服务时，确保只有一个服务在执行定时任务
type Locker struct {
	h LockHandler
	// 是否抢到
	locked bool
	// 通知
	exit *Signal[struct{}]
	// 上一次调用的时间
	time time.Time
	// 是否正在调用
	handing int32
}

// 开始抢锁
func RunLocker(handler LockHandler) *Locker {
	l := new(Locker)
	l.h = handler
	l.exit = NewSignal[struct{}]()
	go l.routine()
	return l
}

// 是否获得锁
func (l *Locker) Locked() bool {
	return l.locked
}

// 停止抢锁
func (l *Locker) Stop() {
	l.exit.Close(struct{}{})
}

// 并发抢锁
func (l *Locker) routine() {
	// 计时器
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-l.exit.C:
			return
		case now := <-timer.C:
			// 处理
			l.work(now)
		}
		// 休息
		timer.Reset(l.h.LockInterval())
	}
}

func (l *Locker) work(now time.Time) {
	// 禁用
	if !l.h.IsEnable() {
		return
	}
	// 抢锁
	var err error
	l.locked, err = l.h.Lock()
	if err != nil {
		l.h.OnError(err)
		return
	}
	// 抢到
	if l.locked && atomic.CompareAndSwapInt32(&l.handing, 0, 1) {
		if now.Sub(l.time) >= l.h.HandleInterval() {
			l.time = now
			go l.handleRoutine()
		} else {
			atomic.StoreInt32(&l.handing, 0)
		}
	}
}

// lockRoutine 处理
func (l *Locker) handleRoutine() {
	defer func() {
		atomic.StoreInt32(&l.handing, 0)
	}()
	l.h.Handle()
}
