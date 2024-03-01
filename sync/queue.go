package sync

import (
	"container/list"
	"sync"
)

// NewQueue 返回初始化的 Queue
func NewQueue[V any](max int) *Queue[V] {
	m := new(Queue[V])
	m.Init(max)
	return m
}

// Queue 封装同步的 slice
type Queue[V any] struct {
	sync.RWMutex
	D   list.List
	max int
}

// Init 初始化
func (m *Queue[V]) Init(max int) {
	m.D.Init()
	m.max = max
}

// Len 返回数量
func (m *Queue[V]) Len() int {
	m.RLock()
	n := m.D.Len()
	m.RUnlock()
	return n
}

// Push 添加
func (m *Queue[V]) Push(v V) {
	m.Lock()
	// 如果超出最大，移除第一个
	if m.max > 0 && m.max <= m.D.Len() {
		m.D.Remove(m.D.Front())
	}
	m.D.PushBack(v)
	m.Unlock()
}

// Pop 获取第一个
func (m *Queue[V]) Pop(i int) (v V) {
	m.RLock()
	if e := m.D.Front(); e != nil {
		v = e.Value.(V)
	}
	m.RUnlock()
	return
}

// Copy 返回拷贝
func (m *Queue[V]) Copy() (vs []V) {
	m.RLock()
	vs = m.copy()
	m.RUnlock()
	return
}

func (m *Queue[V]) copy() (vs []V) {
	for ele := m.D.Front(); ele != nil; ele = ele.Next() {
		vs = append(vs, ele.Value.(V))
	}
	return
}

// TakeAll 返回所有值，清除列表
func (m *Queue[V]) TakeAll() (vs []V) {
	m.Lock()
	vs = m.copy()
	// 重置
	m.D.Init()
	m.Unlock()
	return
}
