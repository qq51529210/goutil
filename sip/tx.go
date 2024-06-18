package sip

import (
	"bytes"
	"sync/atomic"
	"time"
)

type tx interface {
	// context.Context 接口
	Done() <-chan struct{}
	// context.Context 接口
	Err() error
	// context.Context 接口
	Deadline() (time.Time, bool)
	// 返回事务的标识
	ID() string
	// 完成通知
	finish(error)
	// 为了在处理 Request 中进行抽象调用
	// udp 发送后会加入到响应缓存
	// tcp 就直接发送了
	writeMsg(conn, *Message) error
}

type baseTx struct {
	id string
	// 状态
	ok int32
	// 信号
	done chan struct{}
	// 错误
	err error
	// 用于判断超时清理
	deadline time.Time
}

// 实现 context.Context 接口
func (t *baseTx) Done() <-chan struct{} {
	return t.done
}

// 实现 context.Context 接口
func (t *baseTx) Err() error {
	return t.err
}

// 实现 context.Context 接口
func (t *baseTx) Deadline() (time.Time, bool) {
	return t.deadline, true
}

func (t *baseTx) ID() string {
	return t.id
}

func (t *baseTx) finish(err error) {
	if atomic.CompareAndSwapInt32(&t.ok, 0, 1) {
		t.err = err
		close(t.done)
	}
}

// udpActiveTx 主动发起请求的事务
type udpActiveTx struct {
	baseTx
	// 连接
	conn *udpConn
	// 请求的数据
	data any
	// 消息重发的间隔，发送一次叠加一倍
	rto time.Duration
	// 发送的数据
	rtoData *bytes.Buffer
	// 发送数据的时间
	rtoTime *time.Time
	// 停止 rto
	rtoStop bool
}

func (t *udpActiveTx) writeMsg(c conn, m *Message) error {
	var b bytes.Buffer
	m.Enc(&b)
	// 保留，在 rto 检查时使用
	t.rtoData = &b
	tt := time.Now()
	t.rtoTime = &tt
	// 先发一次
	return c.write(b.Bytes())
}

// udpPassiveTx 被动接受请求的事务
type udpPassiveTx struct {
	baseTx
	// 用于控制多消息并发时的单一处理
	handing int32
	// 响应的数据缓存
	dataBuff *bytes.Buffer
}

func (t *udpPassiveTx) writeMsg(c conn, m *Message) error {
	var b bytes.Buffer
	m.Enc(&b)
	// 保留，在下一次重复请求时使用
	t.dataBuff = &b
	return c.write(b.Bytes())
}

// tcpActiveTx 主动发起请求的事务
type tcpActiveTx struct {
	baseTx
	// 请求的数据
	data any
}

func (t *tcpActiveTx) writeMsg(c conn, m *Message) error {
	var b bytes.Buffer
	m.Enc(&b)
	return c.write(b.Bytes())
}

// tcpPassiveTx 被动接受请求的事务
type tcpPassiveTx struct {
	baseTx
	// 用于控制多消息并发时的单一处理
	handing int32
}

func (t *tcpPassiveTx) writeMsg(c conn, m *Message) error {
	var b bytes.Buffer
	m.Enc(&b)
	return c.write(b.Bytes())
}
