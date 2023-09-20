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
func (s *Slice[K]) Init() {
	s.D = make([]K, 0)
}

// Set 设置，存在返回 false
func (s *Slice[K]) Set(k K) bool {
	ok := false
	s.Lock()
	for i := 0; i < len(s.D); i++ {
		if s.D[i] == k {
			ok = true
			break
		}
	}
	s.Unlock()
	//
	if !ok {
		s.D = append(s.D, k)
	}
	return !ok
}

// Has 是否存在
func (s *Slice[K]) Has(k K) bool {
	ok := false
	s.RLock()
	for i := 0; i < len(s.D); i++ {
		if s.D[i] == k {
			ok = true
			break
		}
	}
	s.RUnlock()
	return ok
}

// Del 移除
func (s *Slice[K]) Del(k K) {
	s.Lock()
	for i := 0; i < len(s.D); i++ {
		if s.D[i] == k {
			s.D = append(s.D[:i], s.D[i+1:]...)
			break
		}
	}
	s.Unlock()
}

// Copy 返回拷贝
func (s *Slice[K]) Copy() []K {
	var k []K
	s.RLock()
	for i := 0; i < len(s.D); i++ {
		k = append(k, s.D[i])
	}
	s.RUnlock()
	return k
}

// Set 封装同步的 map
type Set[K comparable] struct {
	sync.RWMutex
	D map[K]struct{}
}

// Init 初始化
func (s *Set[K]) Init() {
	s.D = make(map[K]struct{})
}

// Set 设置
func (s *Set[K]) Set(k K) {
	s.Lock()
	s.D[k] = struct{}{}
	s.Unlock()
}

// TrySet 尝试设置，存在返回 false
func (s *Set[K]) TrySet(k K) bool {
	s.Lock()
	_, ok := s.D[k]
	if !ok {
		s.D[k] = struct{}{}
	}
	s.Unlock()
	return !ok
}

// Has 是否存在
func (s *Set[K]) Has(k K) bool {
	s.RLock()
	_, ok := s.D[k]
	s.RUnlock()
	return ok
}

// Del 移除
func (s *Set[K]) Del(k K) {
	s.Lock()
	delete(s.D, k)
	s.Unlock()
}

// Keys 返回所有键
func (s *Set[K]) Keys() (k []K) {
	s.RLock()
	for d := range s.D {
		k = append(k, d)
	}
	s.RUnlock()
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
func (m *Map[K, V]) Init() {
	m.D = make(map[K]V)
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
