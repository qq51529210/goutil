package util

import (
	"sync"
)

// Slice 封装同步的 map
type Slice[K comparable] struct {
	sync.RWMutex
	D []K
}

// Init 初始化
func (p *Slice[K]) Init() {
	p.D = make([]K, 0)
}

// Len 返回数量
func (p *Slice[K]) Len() int {
	p.RLock()
	n := len(p.D)
	p.RUnlock()
	return n
}

// Set 设置，存在返回 false
func (p *Slice[K]) Set(k K) bool {
	ok := false
	p.Lock()
	for i := 0; i < len(p.D); i++ {
		if p.D[i] == k {
			ok = true
			break
		}
	}
	p.Unlock()
	//
	if !ok {
		p.D = append(p.D, k)
	}
	return !ok
}

// Has 是否存在
func (p *Slice[K]) Has(k K) bool {
	ok := false
	p.RLock()
	for i := 0; i < len(p.D); i++ {
		if p.D[i] == k {
			ok = true
			break
		}
	}
	p.RUnlock()
	return ok
}

// Del 移除
func (p *Slice[K]) Del(k K) {
	p.Lock()
	for i := 0; i < len(p.D); i++ {
		if p.D[i] == k {
			p.D = append(p.D[:i], p.D[i+1:]...)
			break
		}
	}
	p.Unlock()
}

// Copy 返回拷贝
func (p *Slice[K]) Copy() []K {
	var k []K
	p.RLock()
	for i := 0; i < len(p.D); i++ {
		k = append(k, p.D[i])
	}
	p.RUnlock()
	return k
}

// Set 封装同步的 map
type Set[K comparable] struct {
	sync.RWMutex
	D map[K]struct{}
}

// Init 初始化
func (p *Set[K]) Init() {
	p.D = make(map[K]struct{})
}

// Len 返回数量
func (p *Set[K]) Len() int {
	p.RLock()
	n := len(p.D)
	p.RUnlock()
	return n
}

// Set 设置
func (p *Set[K]) Set(k K) {
	p.Lock()
	p.D[k] = struct{}{}
	p.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (p *Set[K]) TrySet(k K) bool {
	p.Lock()
	_, ok := p.D[k]
	if !ok {
		p.D[k] = struct{}{}
	}
	p.Unlock()
	return !ok
}

// Has 是否存在
func (p *Set[K]) Has(k K) bool {
	p.RLock()
	_, ok := p.D[k]
	p.RUnlock()
	return ok
}

// Del 移除
func (p *Set[K]) Del(k K) {
	p.Lock()
	delete(p.D, k)
	p.Unlock()
}

// Keys 返回所有键
func (p *Set[K]) Keys() (k []K) {
	p.RLock()
	for d := range p.D {
		k = append(k, d)
	}
	p.RUnlock()
	return
}

// Map 封装同步的 map
type Map[K comparable, V any] struct {
	// 同步锁
	sync.RWMutex
	// 数据
	D map[K]V
}

// Init 初始化
func (p *Map[K, V]) Init() {
	p.D = make(map[K]V)
}

// Len 返回数量
func (p *Map[K, V]) Len() int {
	p.RLock()
	n := len(p.D)
	p.RUnlock()
	return n
}

// Set 设置
func (p *Map[K, V]) Set(k K, v V) {
	p.Lock()
	p.D[k] = v
	p.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (p *Map[K, V]) TrySet(k K, v V) bool {
	p.Lock()
	_, ok := p.D[k]
	if !ok {
		p.D[k] = v
	}
	p.Unlock()
	return !ok
}

// Get 返回
func (p *Map[K, V]) Get(k K) (v V) {
	p.RLock()
	v = p.D[k]
	p.RUnlock()
	return v
}

// Del 移除
func (p *Map[K, V]) Del(k K) {
	p.Lock()
	delete(p.D, k)
	p.Unlock()
}

// Values 返回所有值
func (p *Map[K, V]) Values() (v []V) {
	p.RLock()
	for _, d := range p.D {
		v = append(v, d)
	}
	p.RUnlock()
	return
}

// Keys 返回所有键
func (p *Map[K, V]) Keys() (k []K) {
	p.RLock()
	for d := range p.D {
		k = append(k, d)
	}
	p.RUnlock()
	return
}

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

// ResetSlice 重置切片，没有锁
func (p *MapSlice[K, V]) ResetSlice() {
	p.S = p.S[:0]
	for _, v := range p.D {
		p.S = append(p.S, v)
	}
}
