package uid

import (
	"sync"
	"time"
)

var (
	sid snowflake
)

func init() {
	sid.init()
}

// snowflake 用于生成雪花 id
type snowflake struct {
	sync.Mutex
	// timestamp, 41 bit
	ts int64
	// serial number, 12 bit
	sn uint16
	// max serial number
	msn uint16
	// mechine id, 10 bit
	mid uint64
	// mechine id mask, 0b00000011 11111111
	midm uint16
	// 用于压缩 uin64
	tb []byte
}

// init 初始化
func (s *snowflake) init() {
	s.ts = time.Now().UTC().Unix()
	s.msn = uint16(0b00001111)<<8 | 0b11111111
	s.midm = uint16(0b00000011)<<8 | 0b11111111
	s.tb = make([]byte, 10+('z'-'a'+1)+('Z'-'A'+1))
	n := 0
	for i := byte('0'); i <= '9'; i++ {
		s.tb[n] = i
		n++
	}
	for i := byte('a'); i <= 'z'; i++ {
		s.tb[n] = i
		n++
	}
	for i := byte('A'); i <= 'Z'; i++ {
		s.tb[n] = i
		n++
	}
}

// SetSnowflakeGroupID 设置机器
func SetSnowflakeMechineID(id uint16) {
	sid.Lock()
	sid.mid = uint64(id&sid.midm) << 12
	sid.Unlock()
}

// SnowflakeID 返回 id
func SnowflakeID() uint64 {
	// 当前时间
	ts := time.Now().UnixNano()
	var sn uint16
	var mid uint64
	sid.Lock()
	// 相同
	if ts == sid.ts {
		// 序列号自增
		sid.sn++
		// 序列号溢出
		if sid.sn > sid.msn {
			// 归零
			sid.sn = 0
			// 直接递增
			ts++
		}
	} else {
		// 序列号归零
		sid.sn = 0
		// 这里假设时间不会回退
		sid.ts = ts
	}
	// 保存当前，因为要解锁了
	sn = sid.sn
	mid = sid.mid
	sid.Unlock()
	// 41 timestamp | 10 bit mechine id | 12 serial number
	return uint64(ts)<<22 | mid | uint64(sn&sid.msn)
}

// SnowflakeIDString 返回压缩成 alpha+number 的字符串
func SnowflakeIDString() string {
	return SnowflakeIDFrom(SnowflakeID())
}

// SnowflakeIDFrom 返回压缩成 alpha+number 的字符串
func SnowflakeIDFrom(id uint64) string {
	b := make([]byte, 20)
	i := 19
	m := uint64(0)
	for {
		m = id % uint64(len(sid.tb))
		b[i] = sid.tb[m]
		i--
		id = id / uint64(len(sid.tb))
		if id == 0 {
			break
		}
	}
	return string(b[i:])
}
