package sync

import "sync"

// Set 封装同步的 map 表示集合
type Set[K comparable] struct {
	sync.RWMutex
	D map[K]struct{}
}

// Init 初始化
func (m *Set[K]) Init() {
	m.D = make(map[K]struct{})
}

// Len 返回数量
func (m *Set[K]) Len() int {
	m.RLock()
	n := len(m.D)
	m.RUnlock()
	return n
}

// Set 设置
func (m *Set[K]) Set(k K) {
	m.Lock()
	m.D[k] = struct{}{}
	m.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (m *Set[K]) TrySet(k K) bool {
	m.Lock()
	_, ok := m.D[k]
	if !ok {
		m.D[k] = struct{}{}
	}
	m.Unlock()
	return !ok
}

// Has 是否存在
func (m *Set[K]) Has(k K) bool {
	m.RLock()
	_, ok := m.D[k]
	m.RUnlock()
	return ok
}

// Del 移除
func (m *Set[K]) Del(k K) {
	m.Lock()
	delete(m.D, k)
	m.Unlock()
}

// Keys 返回所有键
func (m *Set[K]) Keys() (k []K) {
	m.RLock()
	for d := range m.D {
		k = append(k, d)
	}
	m.RUnlock()
	return
}

// TakeAll 返回所有值，清除列表
func (m *Set[K]) TakeAll() (k []K) {
	m.Lock()
	for d := range m.D {
		k = append(k, d)
	}
	m.D = make(map[K]struct{})
	m.Unlock()
	return
}
