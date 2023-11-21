package gorm

import (
	"goutil/log"
	"time"
)

// MustLoad 一直加载知道成功
func MustLoad[K comparable](k K, f func(k K) bool) {
	go mustLoadRoutine(k, f)
}

// mustLoadRoutine 是 MustLoad 的协程
func mustLoadRoutine[K comparable](k K, f func(k K) bool) {
	defer func() {
		log.Recover(recover())
	}()
	for !f(k) {
		time.Sleep(time.Second)
	}
}
