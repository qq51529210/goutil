package sync

import "sync"

// NewSlice 返回初始化的 Slice
func NewSlice[V any]() *Slice[V] {
	m := new(Slice[V])
	m.Init()
	return m
}

// Slice 封装同步的 slice
type Slice[V any] struct {
	sync.RWMutex
	D []V
}

// Init 初始化
func (m *Slice[V]) Init() {
	m.D = make([]V, 0)
}

// Len 返回数量
func (m *Slice[V]) Len() int {
	m.RLock()
	n := len(m.D)
	m.RUnlock()
	return n
}

// Add 添加
func (m *Slice[V]) Add(v V) {
	m.Lock()
	m.D = append(m.D, v)
	m.Unlock()
}

// Add 添加
func (m *Slice[V]) Get(i int) (v V) {
	m.RLock()
	if i < len(m.D) {
		v = m.D[i]
	}
	m.RUnlock()
	return
}

// Del 移除
func (m *Slice[V]) Del(i int) {
	m.Lock()
	if i < len(m.D) {
		m.D = append(m.D[:i], m.D[i+1:]...)
	}
	m.Unlock()
}

// Copy 返回拷贝
func (m *Slice[V]) Copy() []V {
	var k []V
	m.RLock()
	for i := 0; i < len(m.D); i++ {
		k = append(k, m.D[i])
	}
	m.RUnlock()
	return k
}

// TakeAll 返回所有值，清除列表
func (m *Slice[V]) TakeAll() (d []V) {
	m.Lock()
	d = m.D
	m.D = make([]V, 0)
	m.Unlock()
	return
}
