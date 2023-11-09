package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const (
	uuidBufLen = 16 + 36
)

var (
	_uuid uuid
)

func init() {
	_uuid.init()
}

// uuid 用于生成 uuid
type uuid struct {
	sync.Mutex
	// clock seq
	sn uint16
	// max clock seq
	msn uint16
	// timestamp
	ts int64
	// 节点，初始化使用第一个网卡
	node []byte
	// 十六进制小写字符表
	ltb []byte
	// 十六进制大写字符表
	utb []byte
	// v2 gid
	gid int
	// v2 uid
	uid int
	// v4 随机数
	rand *rand.Rand
}

func (m *uuid) init() {
	m.ts = time.Now().UnixNano()
	m.msn = uint16(0b00001111)<<8 | 0b11111111
	m.node = make([]byte, 6)
	m.ltb = []byte("0123456789abcdef")
	m.utb = []byte("0123456789ABCDEF")
	// v2
	m.gid = os.Getgid()
	m.uid = os.Getuid()
	// v4
	m.rand = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	// 使用网卡初始化 node
	ifs, err := net.Interfaces()
	if nil != err {
		// 读取不了就随机
		m.rand.Read(m.node[:])
	} else {
		ok := false
		for i := 0; i < len(ifs); i++ {
			if len(ifs[i].HardwareAddr) >= 6 {
				copy(m.node[0:], ifs[i].HardwareAddr)
				ok = true
				break
			}
		}
		// 没有就随机
		if !ok {
			m.rand.Read(m.node[:])
		}
	}
}

func (m *uuid) time() (ts int64, sn uint16) {
	// 时间
	ts = time.Now().UnixNano()
	// 序列号
	m.Lock()
	if m.ts == ts {
		// 序列号自增
		m.sn++
		// 序列号溢出
		if m.sn > m.msn {
			// 归零
			sid.sn = 0
			// 直接递增
			ts++
		}
	} else {
		// 序列号归零
		m.sn = 0
	}
	sn = m.sn
	m.ts = ts
	m.Unlock()
	//
	return
}

func (m *uuid) new(data []byte) {
	// 节点
	copy(data, m.node)
	// 时间
	ts, sn := m.time()
	binary.BigEndian.PutUint64(data[6:], uint64(ts))
	// 时钟序列
	binary.BigEndian.PutUint16(data[14:], sn)
}

func (m *uuid) v1(data []byte) {
	// 时间
	ts, sn := m.time()
	binary.LittleEndian.PutUint64(data, uint64(ts))
	// 版本
	data[6] = (data[6] & 0x0f) | 0x10
	// 变种和序列号高位
	data[8] = (byte(sn>>8) & 0x3f) | 0x80
	// 序列号低位
	data[9] = byte(sn)
	// 节点
	copy(data[10:], m.node)
}

func (m *uuid) v2(data []byte) {
	// 时间
	ts, sn := m.time()
	binary.LittleEndian.PutUint64(data, uint64(ts))
	// 版本
	data[6] = (data[6] & 0x0f) | 0x20
	// 变种和序列号高位
	data[8] = (byte(sn>>8) & 0x3f) | 0x80
	// 序列号低位
	data[9] = byte(sn)
	// 节点
	copy(data[10:], m.node)
}

func (m *uuid) v3(namespace, name, data []byte) {
	h := md5.New()
	h.Reset()
	h.Write(namespace)
	h.Write(name)
	b := h.Sum(nil)
	b[6] = (b[6] & 0x0f) | 0x30
	b[8] = (b[8] & 0x3f) | 0x80
	copy(data, b)
}

func (m *uuid) v4(data []byte) {
	m.rand.Read(data)
	data[6] = (data[6] & 0x0f) | 0x40
	data[8] = (data[8] & 0x3f) | 0x80
}

func (m *uuid) v5(namespace, name, data []byte) {
	h := sha1.New()
	h.Reset()
	h.Write(namespace)
	h.Write(name)
	b := h.Sum(nil)
	b[6] = (b[6] & 0x0f) | 0x50
	b[8] = (b[8] & 0x3f) | 0x80
	copy(data, b)
}

func (m *uuid) hexString(buf []byte, upper, hyphen bool) string {
	if upper {
		if hyphen {
			m.hex(m.utb, buf)
			return string(buf[16:])
		}
		m.hexWithoutHyphen(m.utb, buf)
		return string(buf[16:48])
	}
	if hyphen {
		m.hex(m.ltb, buf)
		return string(buf[16:])
	}
	m.hexWithoutHyphen(m.ltb, buf)
	return string(buf[16:48])
}

func (m *uuid) hex(table, buf []byte) {
	i, j := 0, 16
	for i < 4 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
	}
	buf[j] = '-'
	j++
	for i < 10 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
		buf[j] = '-'
		j++
	}
	for i < 16 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
	}
}

func (m *uuid) hexWithoutHyphen(table, buf []byte) {
	i, j := 0, 16
	for i < 4 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
	}
	for i < 10 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
	}
	for i < 16 {
		buf[j] = table[buf[i]>>4]
		j++
		buf[j] = table[buf[i]&0x0f]
		j++
		i++
	}
}

// SetUUIDNode 设置 node ，取 6 个字节
func SetUUIDNode(node string) {
	copy(_uuid.node[:], []byte(node))
}

// UUID 返回小写带连字符的
// 自己实现的，没有版本号的，不标准的
func UUID() string {
	buf := make([]byte, uuidBufLen)
	_uuid.new(buf)
	return _uuid.hexString(buf, false, true)
}

// UUIDNoHyphen 返回小写不带连字符的
// 自己实现的，没有版本号的，不标准的
func UUIDNoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.new(buf)
	return _uuid.hexString(buf, false, false)
}

// UUIDUpper 返回小写带连字符的
// 自己实现的，没有版本号的，不标准的
func UUIDUpper() string {
	buf := make([]byte, uuidBufLen)
	_uuid.new(buf)
	return _uuid.hexString(buf, true, true)
}

// UUIDUpperNoHyphen 返回小写不带连字符的
// 自己实现的，没有版本号的，不标准的
func UUIDUpperNoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.new(buf)
	return _uuid.hexString(buf, true, false)
}

// UUID1 版本 1 ，返回小写带连字符的
func UUID1() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v1(buf)
	return _uuid.hexString(buf, false, true)
}

// UUID1NoHyphen 版本 1 ，返回小写不带连字符的
func UUID1NoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v1(buf)
	return _uuid.hexString(buf, false, false)
}

// UUID1Upper 版本 1 ，返回大写带连字符的
func UUID1Upper() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v1(buf)
	return _uuid.hexString(buf, true, true)
}

// UUID1UpperNoHyphen 版本 1 ，返回大写不带连字符的
func UUID1UpperNoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v1(buf)
	return _uuid.hexString(buf, true, false)
}

// UUID2 版本 2 ，返回小写带连字符的
func UUID2() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v2(buf)
	return _uuid.hexString(buf, false, true)
}

// UUID2NoHyphen 版本 2 ，返回小写不带连字符的
func UUID2NoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v2(buf)
	return _uuid.hexString(buf, false, false)
}

// UUID2Upper 版本 2 ，返回大写带连字符的
func UUID2Upper() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v2(buf)
	return _uuid.hexString(buf, true, true)
}

// UUID2UpperNoHyphen 版本 2 ，返回大写不带连字符的
func UUID2UpperNoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v2(buf)
	return _uuid.hexString(buf, true, false)
}

// UUID3 版本 3 ，返回小写带连字符的
func UUID3(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v3(namespace, name, buf)
	return _uuid.hexString(buf, false, true)
}

// UUID3NoHyphen 版本 3 ，返回小写不带连字符的
func UUID3NoHyphen(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v3(buf, namespace, name)
	return _uuid.hexString(buf, false, false)
}

// UUID3Upper 版本 3 ，返回大写带连字符的
func UUID3Upper(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v3(buf, namespace, name)
	return _uuid.hexString(buf, true, true)
}

// UUID3UpperNoHyphen 版本 3 ，返回大写不带连字符的
func UUID3UpperNoHyphen(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v3(buf, namespace, name)
	return _uuid.hexString(buf, true, false)
}

// UUID4 版本 4 ，返回小写带连字符的
func UUID4() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v4(buf)
	return _uuid.hexString(buf, false, true)
}

// UUID4NoHyphen 版本 4 ，返回小写不带连字符的
func UUID4NoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v4(buf)
	return _uuid.hexString(buf, false, false)
}

// UUID4Upper 版本 4 ，返回大写带连字符的
func UUID4Upper() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v4(buf)
	return _uuid.hexString(buf, true, true)
}

// UUID4UpperNoHyphen 版本 4 ，返回大写不带连字符的
func UUID4UpperNoHyphen() string {
	buf := make([]byte, uuidBufLen)
	_uuid.v4(buf)
	return _uuid.hexString(buf, true, false)
}

// UUID5 版本 5 ，返回小写带连字符的
func UUID5(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v5(namespace, name, buf)
	return _uuid.hexString(buf, false, true)
}

// UUID5NoHyphen 版本 5 ，返回小写不带连字符的
func UUID5NoHyphen(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v5(buf, namespace, name)
	return _uuid.hexString(buf, false, false)
}

// UUID5Upper 版本 5 ，返回大写带连字符的
func UUID5Upper(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v5(buf, namespace, name)
	return _uuid.hexString(buf, true, true)
}

// UUID5UpperNoHyphen 版本 5 ，返回大写不带连字符的
func UUID5UpperNoHyphen(namespace, name []byte) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v5(buf, namespace, name)
	return _uuid.hexString(buf, true, false)
}

// UUIDFrom 从两个 64 位的整数生成
func UUIDFrom(n1, n2 uint64, upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	binary.LittleEndian.PutUint64(buf[0:], n1)
	binary.LittleEndian.PutUint64(buf[8:], n2)
	return _uuid.hexString(buf, upper, hyphen)
}

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

// SetSnowflakeMechineID 设置机器
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
		id = id / uint64(len(sid.tb))
		if id == 0 {
			break
		}
		i--
	}
	return string(b[i:])
}
