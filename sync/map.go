package sync

import "sync"

// NewMap 返回初始化的 Map
func NewMap[K comparable, V any]() *Map[K, V] {
	m := new(Map[K, V])
	m.Init()
	return m
}

// Map 封装同步的 map
type Map[K comparable, V any] struct {
	// 同步锁
	sync.RWMutex
	// 数据
	D map[K]V
}

// Init 初始化
func (m *Map[K, V]) Init() {
	m.D = make(map[K]V)
}

// Len 返回数量
func (m *Map[K, V]) Len() int {
	m.RLock()
	n := len(m.D)
	m.RUnlock()
	return n
}

// Set 设置
func (m *Map[K, V]) Set(k K, v V) {
	m.Lock()
	m.D[k] = v
	m.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (m *Map[K, V]) TrySet(k K, v V) bool {
	m.Lock()
	_, ok := m.D[k]
	if !ok {
		m.D[k] = v
	}
	m.Unlock()
	return !ok
}

// Has 是否存在
func (m *Map[K, V]) Has(k K) bool {
	m.RLock()
	_, ok := m.D[k]
	m.RUnlock()
	return ok
}

// Get 返回
func (m *Map[K, V]) Get(k K) (v V) {
	m.RLock()
	v = m.D[k]
	m.RUnlock()
	return v
}

// Del 移除
func (m *Map[K, V]) Del(k K) {
	m.Lock()
	delete(m.D, k)
	m.Unlock()
}

// Values 返回所有值
func (m *Map[K, V]) Values() (v []V) {
	m.RLock()
	for _, d := range m.D {
		v = append(v, d)
	}
	m.RUnlock()
	return
}

// Keys 返回所有键
func (m *Map[K, V]) Keys() (k []K) {
	m.RLock()
	for d := range m.D {
		k = append(k, d)
	}
	m.RUnlock()
	return
}

// Take 移除后返回
func (m *Map[K, V]) Take(k K) V {
	m.Lock()
	v := m.D[k]
	delete(m.D, k)
	m.Unlock()
	return v
}

// TakeAll 返回所有值，清除列表
func (m *Map[K, V]) TakeAll() (v []V) {
	m.Lock()
	for _, d := range m.D {
		v = append(v, d)
	}
	m.D = make(map[K]V)
	m.Unlock()
	return
}
