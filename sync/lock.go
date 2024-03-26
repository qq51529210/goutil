package sync

import (
	"goutil/log"
	"sync/atomic"
	"time"
)

// Locker 用于控制在部署多个服务时，确保只有一个服务在执行定时任务
type Locker struct {
	// 日志
	Trace  string
	Logger *log.Logger
	// 抢锁周期
	Interval time.Duration
	// 抢锁
	Lock func() (bool, error)
	// 时间到了，是否启动协程回调 Handle
	IsTimeup func(d time.Duration) bool
	// 启用/禁用，协程还是在跑，只是不会回调 Handle
	IsEnabled func() bool
	// 在协程中回调
	Handle func()
	// 是否抢到
	locked bool
	// 并发控制
	time    time.Time
	handing int32
}

// IsLocked 是否获得锁
func (l *Locker) IsLocked() bool {
	return l.locked
}

// Run 开始
func (l *Locker) Run() {
	go l.routine()
}

// routine 抢锁协程
func (l *Locker) routine() {
	defer func() {
		// 异常
		l.Logger.Recover(recover())
		// 日志
		l.Logger.InfoTrace(l.Trace, "stop")
	}()
	// 执行
	l.Logger.InfoTrace(l.Trace, "start")
	timer := time.NewTimer(0)
	for {
		now := <-timer.C
		// 处理
		l.handle(&now)
		// 休息
		if l.Interval > 0 {
			timer.Reset(l.Interval)
		}
	}
}

func (l *Locker) handle(now *time.Time) {
	// 禁用
	if !l.IsEnabled() {
		return
	}
	// 抢锁
	var err error
	l.locked, err = l.Lock()
	if err != nil {
		l.Logger.ErrorTrace(l.Trace, err)
		return
	}
	// 抢到
	if l.locked && l.IsTimeup(now.Sub(l.time)) && atomic.CompareAndSwapInt32(&l.handing, 0, 1) {
		l.time = *now
		go l.handleRoutine()
	}
}

// lockRoutine 处理
func (l *Locker) handleRoutine() {
	defer func() {
		// 异常
		l.Logger.Recover(recover())
		// 标记
		atomic.StoreInt32(&l.handing, 0)
	}()
	l.Handle()
}
