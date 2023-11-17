package sync

import (
	"goutil/log"
	"sync/atomic"
	"time"
)

// Locker 用于控制在部署多个服务时，确保只有一个服务在执行定时任务
type Locker struct {
	// 日志
	Trace string
	// 是否抢到
	locked bool
	// 并发控制
	time    time.Time
	handing int32
	// 抢锁
	LockFunc func() (bool, error)
	// 时间到了，是否启动协程回调 HandleFunc
	TimeupFunc func(d time.Duration) bool
	// 启用/禁用，协程还是在跑，只是不会回调 HandleFunc
	EnableFunc func() bool
	// 在协程中回调
	HandleFunc func()
}

// Run 开始
func (l *Locker) Run() {
	go l.routine()
}

// routine 抢锁协程
func (l *Locker) routine() {
	defer func() {
		// 异常
		log.Recover(recover())
		// 日志
		log.InfoTrace(l.Trace, "stop")
	}()
	// 执行
	log.InfoTrace(l.Trace, "start")
	timer := time.NewTimer(0)
	for {
		now := <-timer.C
		if l.EnableFunc() {
			// 抢锁
			ok, err := l.LockFunc()
			if err != nil {
				log.ErrorTrace(l.Trace, err)
			} else {
				l.locked = ok
			}
			// 抢到
			if l.locked && l.TimeupFunc(now.Sub(l.time)) && atomic.CompareAndSwapInt32(&l.handing, 0, 1) {
				l.time = now
				go l.handleRoutine()
			}
		}
		// 重置计时器
		timer.Reset(time.Second)
	}
}

// lockRoutine 处理
func (l *Locker) handleRoutine() {
	old := time.Now()
	defer func() {
		// 异常
		log.Recover(recover())
		// 日志
		log.DebugfTrace(l.Trace, "handle cost %v", time.Since(old))
		// 标记
		atomic.StoreInt32(&l.handing, 0)
	}()
	l.HandleFunc()
}
