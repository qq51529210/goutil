package zlm

// import (
// 	"mms/cfg"
// 	"sync"
// 	"time"

// 	"github.com/qq51529210/log"
// 	"github.com/qq51529210/uuid"
// )

// const (
// 	// QueryNameToken 查询参数 token=xx
// 	QueryNameToken = "token"
// )

// var (
// 	_playToken playToken
// )

// type playToken struct {
// 	sync.RWMutex
// 	token map[string]*time.Time
// }

// func (t *playToken) init() {
// 	t.token = make(map[string]*time.Time)
// 	// 启动协程
// 	go t.routine()
// }

// // routine 协程中循环检查过期的 token
// func (t *playToken) routine() {
// 	// 计时器
// 	timer := time.NewTimer(0)
// 	defer func() {
// 		// 异常
// 		log.Recover(recover())
// 		// 计时器
// 		timer.Stop()
// 	}()
// 	for {
// 		now := <-timer.C
// 		//
// 		dur := time.Duration(cfg.Cfg.Media.PlayTokenTimeout) * time.Second / 3
// 		if dur < time.Second {
// 			dur = time.Second
// 		}
// 		// 检查
// 		t.Lock()
// 		for k, v := range t.token {
// 			if v != nil && now.Sub(*v) >= dur {
// 				delete(t.token, k)
// 			}
// 		}
// 		t.Unlock()
// 		// 重置计时器
// 		timer.Reset(dur)
// 	}
// }

// // new 返回新的 token ，不开启返回空字符串
// func (t *playToken) new() string {
// 	if cfg.Cfg.Media.PlayTokenTimeout < 1 {
// 		return ""
// 	}
// 	// 生成
// 	now := time.Now()
// 	token := uuid.SnowflakeIDString()
// 	t.Lock()
// 	t.token[token] = &now
// 	t.Unlock()
// 	// 返回
// 	return token
// }

// // has 返回 token 是否存在，不开启返回 true
// func (t *playToken) has(token string) bool {
// 	if cfg.Cfg.Media.PlayTokenTimeout < 1 {
// 		return true
// 	}
// 	if token == "" {
// 		return false
// 	}
// 	// 时间
// 	dur := time.Duration(cfg.Cfg.Media.PlayTokenTimeout) * time.Second
// 	now := time.Now()
// 	// 上锁
// 	t.Lock()
// 	defer t.Unlock()
// 	tt := t.token[token]
// 	if tt != nil {
// 		// 过期直接删除
// 		if now.Sub(*tt) >= dur {
// 			delete(t.token, token)
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }
