package sip

import (
	"context"
	"fmt"
	"goutil/log"
	gosync "goutil/sync"
	"goutil/uid"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	logTraceUDP = "sip udp"
	logTraceTCP = "sip tcp"
)

// Handler 是处理消息的接口
type Handler interface {
	// 返回 true 表示已经处理
	// 这样事务超时时间内就不会触发相同消息的回调
	// 而是直接回复第一次的响应消息数据
	HandleRequest(*Request) bool
	HandleResponse(*Response)
}

// Server 表示一个服务
type Server struct {
	// 回调函数
	Handler
	// 监听地址
	Addr string
	// udp 超时重发时间
	RTO time.Duration
	// udp 最大超时重发时间
	MaxRTO time.Duration
	// 事务超时时间
	TxTimeout time.Duration
	// 最大的消息字节数，防止内存爆掉哦
	MaxMessageLen int
	// 用于同步等待协程退出
	w sync.WaitGroup
	// 状态
	ok int32
	// udp
	udp udpServer
	// tcp server
	tcp tcpServer
}

func (s *Server) isOK() bool {
	return atomic.LoadInt32(&s.ok) == 0
}

// Serve 开始服务，不会阻塞
func (s *Server) Serve() error {
	// udp 启动
	if err := s.serveUDP(); err != nil {
		return err
	}
	// tcp 启动
	if err := s.serveTCP(); err != nil {
		return err
	}
	// 日志
	log.InfoTrace(SIP, "ok")
	//
	return nil
}

// Close 停止服务
func (s *Server) Close() error {
	if atomic.CompareAndSwapInt32(&s.ok, 0, 1) {
		// 关闭双服务
		s.closeUDP()
		s.closeTCP()
		// 等待所有协程退出
		s.w.Wait()
		// 日志
		log.InfoTrace(SIP, "closed")
	}
	return nil
}

// handleMsg 处理消息
func (s *Server) handleMsg(c conn, m *message, at *gosync.Map[string, *activeTx], pt *gosync.Map[string, *passiveTx]) error {
	// 请求消息
	if m.isReq {
		// 日志
		log.DebugfTrace(SIP, "read request from %s %s\n%v", c.Network(), c.RemoteAddrString(), m)
		// 事务，返回一定不为 nil
		t := s.newPassiveTx(m, pt)
		if atomic.CompareAndSwapInt32(&t.handing, 0, 1) {
			// 在协程中处理
			s.w.Add(1)
			go s.handleRequestRoutine(c, t, m)
		} else {
			// 已经处理过
			if t.done {
				// 有响应数据，直接发送，无需回调
				d := t.writeData.Bytes()
				if len(d) > 0 {
					return c.write(d)
				}
			}
		}
		return nil
	}
	// 日志
	log.DebugfTrace(SIP, "read response from %s %s\n%v", c.Network(), c.RemoteAddrString(), m)
	// 响应消息
	if m.StartLine[1][0] == '1' {
		// 1xx 消息没什么卵用
		return nil
	}
	// 事务，不一定有
	t := s.delActiveTx(m.txKey(), at)
	if t != nil {
		// 在协程中处理
		s.w.Add(1)
		go s.handleResponseRoutine(t, m)
	}
	//
	return nil
}

// handleRequestRoutine 在协程中处理请求消息
func (s *Server) handleRequestRoutine(c conn, t *passiveTx, m *message) {
	defer func() {
		// 异常
		log.Recover(recover())
		// 日志
		log.DebugfTrace(t.TxKey(), "handle request cost %v", time.Since(t.time))
		// 结束
		s.w.Done()
	}()
	// 回调
	t.done = s.HandleRequest(&Request{
		Server:    s,
		message:   m,
		passiveTx: t,
		conn:      c,
	})
	if !t.done {
		// 没有完成，回复标记，等下一次的消息再回调
		atomic.StoreInt32(&t.handing, 0)
	}
}

// handleResponseRoutine 在协程中处理响应消息
func (s *Server) handleResponseRoutine(t *activeTx, m *message) {
	defer func() {
		// 异常
		log.Recover(recover())
		// 无论回调有没有通知，这里都通知一下
		t.Finish(nil)
		// 日志
		log.DebugfTrace(t.TxKey(), "handle response cost %v", time.Since(t.time))
		// 结束
		s.w.Done()
	}()
	// 回调处理
	s.HandleResponse(&Response{
		activeTx: t,
		message:  m,
	})
}

// NewRequest 创建请求
func (s *Server) NewRequest(proto, method, localName, remoteName, remoteAddr, maxForwards, contentType, contact string) *Request {
	m := &Request{message: new(message)}
	// start line
	m.StartLine[0] = method
	m.StartLine[1] = fmt.Sprintf("SIP:%s@%s", remoteName, remoteAddr)
	m.StartLine[2] = SIPVersion
	// via
	m.Header.Via = append(m.Header.Via, &Via{
		Proto:   proto,
		Address: s.Addr,
		Branch:  fmt.Sprintf("%s%d", BranchPrefix, uid.SnowflakeID()),
	})
	// From
	m.Header.From.URI.Scheme = SIP
	m.Header.From.URI.Name = localName
	m.Header.From.URI.Domain = s.Addr
	m.Header.From.Tag = fmt.Sprintf("%d", uid.SnowflakeID())
	// To
	m.Header.To.URI.Scheme = SIP
	m.Header.To.URI.Name = remoteName
	m.Header.To.URI.Domain = remoteAddr
	// m.Header.To.Tag = ""
	// Call-ID
	m.Header.CallID = fmt.Sprintf("%d", uid.SnowflakeID())
	// CSeq
	m.Header.CSeq.SN = GetSNString()
	m.Header.CSeq.Method = method
	// Max-Forwards
	m.Header.MaxForwards = maxForwards
	// Content-Type
	m.Header.ContentType = contentType
	// Contact
	m.Header.Contact.Scheme = SIP
	m.Header.Contact.Name = localName
	m.Header.Contact.Domain = contact
	//
	return m
}

// Request 发送请求并等待响应
func (s *Server) Request(ctx context.Context, r *Request, a net.Addr, d any) error {
	// tcp
	if _a, ok := a.(*net.TCPAddr); ok {
		var err error
		// 获取连接
		c := s.getTCPConn(_a)
		if c == nil {
			// 没有就创建
			c, err = s.dialTCPConn(_a)
			if err != nil {
				return err
			}
		}
		// 请求
		return s.doRequest(ctx, c, r.message, d, &s.tcp.at)
	}
	// udp
	if _a, ok := a.(*net.UDPAddr); ok {
		// 连接
		c := new(udpConn)
		c.conn = s.udp.c
		c.initAddr(_a)
		//
		return s.doRequest(ctx, c, r.message, d, &s.udp.at)
	}
	//
	return errUnknownAddress
}

// doRequest 封装 Request 的公共代码
func (s *Server) doRequest(ctx context.Context, c conn, m *message, d any, at *gosync.Map[string, *activeTx]) error {
	// 事务
	t, err := s.newActiveTx(c, m, d, at)
	if err != nil {
		return err
	}
	// 日志
	log.DebugfTrace(t.TxKey(), "write request to %s %s\n%s", c.Network(), c.RemoteAddrString(), t.writeData.String())
	// 立即发送
	err = c.write(t.writeData.Bytes())
	if err == nil {
		// 等待响应处理或底层超时
		select {
		case <-ctx.Done():
			err = ctx.Err()
			// 移除
			at.Del(t.key)
			// 通知
			t.Finish(err)
		case <-t.signal.C:
			err = t.err
			// 要么是收到了响应的消息被移除
			// 要么是检查超时被移除
			// 所以这里不需要显示调用，提高性能
		}
	}
	// 日志
	log.DebugfTrace(t.TxKey(), "do request cost %v", time.Since(t.time))
	//
	return err
}
