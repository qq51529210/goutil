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

// // uuid 版本
// type UUIDVersion byte

// // uuid 版本 1-5
// const (
// 	UUIDVersionV1 UUIDVersion = iota + 1
// 	UUIDVersionV2
// 	UUIDVersionV3
// 	UUIDVersionV4
// 	UUIDVersionV5
// )

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
	} else {
		if hyphen {
			m.hex(m.ltb, buf)
			return string(buf[16:])
		}
		m.hexWithoutHyphen(m.ltb, buf)
		return string(buf[16:48])
	}
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

// UUID 自己实现的，没有版本号的，不标准的
func UUID(upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.new(buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUID1 版本 1
func UUID1(upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v1(buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUID1 版本 2
func UUID2(upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v2(buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUID1 版本 3
func UUID3(namespace, name []byte, upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v3(namespace, name, buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUID1 版本 4
func UUID4(upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v4(buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUID1 版本 5
func UUID5(namespace, name []byte, upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	_uuid.v5(namespace, name, buf)
	return _uuid.hexString(buf, upper, hyphen)
}

// UUIDFrom 从两个 64 位的整数生成
func UUIDFrom(n1, n2 uint64, upper, hyphen bool) string {
	buf := make([]byte, uuidBufLen)
	binary.LittleEndian.PutUint64(buf[0:], n1)
	binary.LittleEndian.PutUint64(buf[8:], n2)
	return _uuid.hexString(buf, upper, hyphen)
}
