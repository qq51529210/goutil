package zlm

import (
	"fmt"
	"strconv"
	"sync"
	"util/log"
)

const (
	// 最大的 ssrc
	maxSSRC = 9999
)

// ssrc 表示 ssrc 池
type ssrcPool struct {
	sync.Mutex
	d [maxSSRC]int
	n int
}

func (p *ssrcPool) Get() (n int) {
	n = -1
	// 上锁
	p.Lock()
	defer p.Unlock()
	// 查询
	if p.n >= maxSSRC {
		p.n = 0
	}
	for ; p.n < maxSSRC; p.n++ {
		if p.d[p.n] == 0 {
			n = p.n
			p.d[p.n] = 1
			break
		}
	}
	return
}

func (p *ssrcPool) Put(ssrc string) {
	// 解析整数
	if len(ssrc) != 10 {
		return
	}
	// 后四位
	n, err := strconv.ParseInt(ssrc[6:], 10, 32)
	if err != nil {
		log.Error(err)
		return
	}
	if n < 0 || n >= maxSSRC {
		return
	}
	// 设置
	p.Lock()
	p.d[n] = 0
	p.Unlock()
}

// GetRealTimeSSRC 获取国标直播流 ssrc
func (s *Server) GetRealTimeSSRC(gbID string) string {
	n := s.realTimeSSRC.Get()
	if n != -1 {
		return fmt.Sprintf("0%s%.4d", gbID[3:8], n)
	}
	return ""
}

// GetHistorySSRC 获取国标历史流 ssrc
func (s *Server) GetHistorySSRC(gbID string) string {
	n := s.historySSRC.Get()
	if n != -1 {
		return fmt.Sprintf("1%s%.4d", gbID[3:8], n)
	}
	return ""
}

// PutSSRC 回收 ssrc
func (s *Server) PutSSRC(ssrc string) {
	if ssrc[0] == '1' {
		s.historySSRC.Put(ssrc)
	} else {
		s.realTimeSSRC.Put(ssrc)
	}
}
