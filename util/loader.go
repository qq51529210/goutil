package util

import (
	"time"
	"util/log"
)

// DataLoader 后台加载数据，一直到成功
type DataLoader[K comparable, M any] struct {
	// 加载数据回调，成功返回 true
	LoadCB func(k K) bool
}

// Load 在后台协程中加载
func (loader *DataLoader[K, M]) Load(k K) {
	go loader.routine(k)
}

// routine 协程中加载
func (loader *DataLoader[K, M]) routine(k K) {
	// 计时器
	timer := time.NewTimer(0)
	defer func() {
		// 异常
		log.Recover(recover())
		// 计时器
		timer.Stop()
	}()
	// 加载一直到成功
	loader.mustLoad(timer, k)
}

// mustLoad 循环加载一直到成功
func (loader *DataLoader[K, M]) mustLoad(t *time.Timer, k K) {
	for {
		<-t.C
		// 加载
		if loader.LoadCB(k) {
			return
		}
		t.Reset(time.Second)
	}
}
