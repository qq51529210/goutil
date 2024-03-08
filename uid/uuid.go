package uid

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
	// 时钟序列
	sn uint16
	// 最大时钟序列
	maxSN uint16
	// timestamp
	ts int64
	// 节点，初始化使用第一个网卡
	node []byte
	// 十六进制小写字符表
	lhexTb []byte
	// 十六进制大写字符表
	uhexTb []byte
	// v2 gid
	gid int
	// v2 uid
	uid int
	// v4 随机数
	rand *rand.Rand
}

func (m *uuid) init() {
	m.ts = time.Now().UnixNano()
	m.maxSN = uint16(0b00001111)<<8 | 0b11111111
	m.node = make([]byte, 6)
	m.lhexTb = []byte("0123456789abcdef")
	m.uhexTb = []byte("0123456789ABCDEF")
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
	m.Lock()
	// 时间
	ts = time.Now().UnixNano() / 100
	// 与上一次相同
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
	}
	m.ts = ts
	sn = m.sn
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
			m.hex(m.uhexTb, buf)
			return string(buf[16:])
		}
		m.hexWithoutHyphen(m.uhexTb, buf)
		return string(buf[16:48])
	}
	if hyphen {
		m.hex(m.lhexTb, buf)
		return string(buf[16:])
	}
	m.hexWithoutHyphen(m.lhexTb, buf)
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
