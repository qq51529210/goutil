package uid

import (
	"sync"
	"time"
)

const (
	snowflakeMechineIDMask = uint16(0b00000011)<<8 | 0b11111111
)

var (
	_snowflake snowflake
)

func init() {
	_snowflake.init()
}

// snowflake 用于生成雪花 id
type snowflake struct {
	sync.Mutex
	// 时间戳，毫秒, 41 bit
	ts int64
	// 序列号，用于同一时间戳递增, 12 bit
	sn uint16
	// 机器编号, 10 bit
	mid uint64
	// 最大序列号
	maxSN uint16
	// 用于压缩 uin64 成字符串
	hexTb []byte
}

// init 初始化
func (m *snowflake) init() {
	m.ts = time.Now().UnixMilli()
	m.maxSN = uint16(0b00001111)<<8 | 0b11111111
	m.hexTb = make([]byte, 10+('z'-'a'+1)+('Z'-'A'+1))
	n := 0
	for i := byte('0'); i <= '9'; i++ {
		m.hexTb[n] = i
		n++
	}
	for i := byte('a'); i <= 'z'; i++ {
		m.hexTb[n] = i
		n++
	}
	for i := byte('A'); i <= 'Z'; i++ {
		m.hexTb[n] = i
		n++
	}
}

// new 生成
func (m *snowflake) new() uint64 {
	// 当前时间
	ts := time.Now().UnixMilli()
	var sn uint16
	var mid uint64
	m.Lock()
	// 相同
	if ts == m.ts {
		// 序列号自增
		m.sn++
		// 序列号溢出
		if m.sn >= m.maxSN {
			// 归零
			m.sn = 0
			// 直接递增
			ts++
		}
	} else {
		// 序列号归零
		m.sn = 0
	}
	// 这里假设时间不会回退
	m.ts = ts
	sn = m.sn
	mid = m.mid
	m.Unlock()
	// 41 timestamp | 10 bit mechine id | 12 serial number
	return uint64(ts)<<22 | mid | uint64(sn)
}

// setMID 设置机器编号
func (m *snowflake) setMID(id uint16) {
	m.Lock()
	m.mid = uint64(id&snowflakeMechineIDMask) << 12
	m.Unlock()
}

// SetSnowflakeMechineID 设置机器
func SetSnowflakeMechineID(id uint16) {
	_snowflake.setMID(id)
}

// SnowflakeID 返回 id
func SnowflakeID() uint64 {
	return _snowflake.new()
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
		m = id % uint64(len(_snowflake.hexTb))
		b[i] = _snowflake.hexTb[m]
		id = id / uint64(len(_snowflake.hexTb))
		if id == 0 {
			break
		}
		i--
	}
	return string(b[i:])
}
