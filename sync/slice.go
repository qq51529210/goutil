package sync

import "sync"

// NewSlice 返回初始化的 Slice
func NewSlice[K comparable]() *Slice[K] {
	m := new(Slice[K])
	m.Init()
	return m
}

// Slice 封装同步的 slice
type Slice[K comparable] struct {
	sync.RWMutex
	D []K
}

// Init 初始化
func (m *Slice[K]) Init() {
	m.D = make([]K, 0)
}

// Len 返回数量
func (m *Slice[K]) Len() int {
	m.RLock()
	n := len(m.D)
	m.RUnlock()
	return n
}

// Slice 设置，存在返回 false
func (m *Slice[K]) Slice(k K) bool {
	ok := false
	m.Lock()
	for i := 0; i < len(m.D); i++ {
		if m.D[i] == k {
			ok = true
			break
		}
	}
	m.Unlock()
	//
	if !ok {
		m.D = append(m.D, k)
	}
	return !ok
}

// Has 是否存在
func (m *Slice[K]) Has(k K) bool {
	ok := false
	m.RLock()
	for i := 0; i < len(m.D); i++ {
		if m.D[i] == k {
			ok = true
			break
		}
	}
	m.RUnlock()
	return ok
}

// Del 移除
func (m *Slice[K]) Del(k K) {
	m.Lock()
	for i := 0; i < len(m.D); i++ {
		if m.D[i] == k {
			m.D = append(m.D[:i], m.D[i+1:]...)
			break
		}
	}
	m.Unlock()
}

// Copy 返回拷贝
func (m *Slice[K]) Copy() []K {
	var k []K
	m.RLock()
	for i := 0; i < len(m.D); i++ {
		k = append(k, m.D[i])
	}
	m.RUnlock()
	return k
}
