package zlm

import (
	"strconv"
	"sync"
)

const (
	// 最大的 ssrc
	maxSSRC = 9999
)

// SSRC 用于国标的 ssrc
type SSRC struct {
	sync.Mutex
	d [maxSSRC]byte
	i int
}

// Put 回收
func (s *SSRC) Put(str string) bool {
	// 解析整数后四位
	n, err := strconv.ParseInt(str[6:], 10, 32)
	if err != nil {
		return false
	}
	if n < 0 || n >= maxSSRC {
		return false
	}
	s.Lock()
	s.d[n] = 0
	s.Unlock()
	//
	return true
}

// Get 获取
func (s *SSRC) Get() int {
	i := -1
	s.Lock()
	for ; s.i < len(s.d); s.i++ {
		if s.d[s.i] == 0 {
			s.d[s.i] = 1
			i = s.i
			break
		}
	}
	if s.i >= len(s.d) {
		s.i = 0
	}
	s.Unlock()
	//
	return i
}
