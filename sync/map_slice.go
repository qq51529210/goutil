package sync

import (
	"sync"
)

// MapSlice 封装同步的 map + slice
type MapSlice[K comparable, V any] struct {
	// 同步锁
	sync.RWMutex
	// map 数据
	D map[K]V
	// slice 数据
	S []V
}

// Init 初始化
func (p *MapSlice[K, V]) Init() {
	p.D = make(map[K]V)
}

// Len 返回数量
func (p *MapSlice[K, V]) Len() int {
	p.RLock()
	n := len(p.D)
	p.RUnlock()
	return n
}

// All 返回所有
func (p *MapSlice[K, V]) All() []V {
	p.RLock()
	a := p.S
	p.RUnlock()
	return a
}

// Set 添加
func (p *MapSlice[K, V]) Set(k K, v V) {
	p.Lock()
	p.D[k] = v
	p.ResetSlice()
	p.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (p *MapSlice[K, V]) TrySet(k K, v V) bool {
	p.Lock()
	_, ok := p.D[k]
	if !ok {
		p.D[k] = v
		p.ResetSlice()
	}
	p.Unlock()
	return !ok
}

// Has 是否存在
func (p *MapSlice[K, V]) Has(k K) bool {
	p.RLock()
	_, ok := p.D[k]
	p.RUnlock()
	return ok
}

// Get 返回
func (p *MapSlice[K, V]) Get(k K) (v V) {
	p.RLock()
	v = p.D[k]
	p.RUnlock()
	return v
}

// Del 删除
func (p *MapSlice[K, V]) Del(k K) {
	p.Lock()
	n := len(p.D)
	delete(p.D, k)
	if n != len(p.D) {
		p.ResetSlice()
	}
	p.Unlock()
}

// BatchDel 批量删除
func (p *MapSlice[K, V]) BatchDel(k []K) {
	p.Lock()
	n := len(p.D)
	for i := 0; i < len(k); i++ {
		delete(p.D, k[i])
	}
	if n != len(p.D) {
		p.ResetSlice()
	}
	p.Unlock()
}

// Values 返回所有值
func (p *MapSlice[K, V]) Values() (v []V) {
	p.RLock()
	for _, d := range p.D {
		v = append(v, d)
	}
	p.RUnlock()
	return
}

// Keys 返回所有键
func (p *MapSlice[K, V]) Keys() (k []K) {
	p.RLock()
	for d := range p.D {
		k = append(k, d)
	}
	p.RUnlock()
	return
}

// Take 移除后返回
func (p *MapSlice[K, V]) Take(k K) V {
	p.Lock()
	n := len(p.D)
	v := p.D[k]
	delete(p.D, k)
	if n != len(p.D) {
		p.ResetSlice()
	}
	p.Unlock()
	return v
}

// TakeAll 返回所有值，清除列表
func (p *MapSlice[K, V]) TakeAll() (v []V) {
	p.Lock()
	for _, d := range p.D {
		v = append(v, d)
	}
	p.D = make(map[K]V)
	p.S = p.S[:0]
	p.Unlock()
	return
}

// ResetSlice 重置切片，没有锁
func (p *MapSlice[K, V]) ResetSlice() {
	p.S = p.S[:0]
	for _, v := range p.D {
		p.S = append(p.S, v)
	}
}
