package sip

import (
	"bytes"
	"context"
	"fmt"
	"goutil/log"
	"goutil/uid"
	"net"
	"strconv"
	"time"
)

type Server struct {
	// 用户代理
	UserAgent string
	// 日志
	Logger *log.Logger
	// 消息最大的字节数，接收到的消息如果大于这个数会被丢弃
	MaxMessageLen int
	// 发起请求的超时时间，或者响应消息缓存的超时时间
	MsgTimeout time.Duration
	// 回调函数
	handleFunc
	// udp 服务
	udp udpServer
	// tcp 服务
	tcp tcpServer
}

// ServeUDP 启动 udp 服务
// address 是监听地址
// minRTO 消息的超时重发最小间隔
// maxRTO 消息的超时重发最小间隔，从 minRTO 开始，重发一次增加一倍，直到 maxRTO
func (s *Server) ServeUDP(address string, minRTO, maxRTO time.Duration) error {
	s.udp.s = s
	s.udp.minRTO = minRTO
	s.udp.maxRTO = maxRTO
	return s.udp.Serve(address)
}

// ServeTCP 启动 tcp 服务
func (s *Server) ServeTCP(address string) error {
	s.tcp.s = s
	return s.tcp.Serve(address)
}

// Shutdown 停止所有服务，阻塞等待全部退出
func (s *Server) Shutdown() {
	s.udp.Shutdown()
	s.tcp.Shutdown()
}

// Request 使用 context.Background() 调用 RequestWithContext
func (s *Server) Request(msg *Message, addr net.Addr, data any) error {
	return s.RequestWithContext(context.Background(), msg, addr, data)
}

// RequestWithContext 向发送请求，阻塞等待异步响应的通知
// ctx 是上下文，用于控制底层结束
// msg 是消息
// addr 是对方的地址，如果是 tcp 类型，而且不存在连接池中，则主动创建连接
// data 是需要传递的上下文数据，可以在异步响应的回调函数 Context.Value(nil) 拿到
// 返回的错误是 Context.Err() 或者是 ctx.Err()
func (s *Server) RequestWithContext(ctx context.Context, msg *Message, addr net.Addr, data any) error {
	// tcp
	if a, ok := addr.(*net.TCPAddr); ok {
		return s.tcp.Request(ctx, msg, a, data)
	}
	// udp
	if a, ok := addr.(*net.UDPAddr); ok {
		return s.udp.Request(ctx, msg, a, data)
	}
	// 其他
	return ErrUnknownAddress
}

func (s *Server) handleRequestNotFound(conn conn, msg *Message) {
	msg.Header.KeepBasic()
	msg.Header.Set("Allow", s.handleFunc.reqMethods)
	msg.Body.Reset()
	if err := s.response(conn, msg, StatusMethodNotAllowed, ""); err != nil {
		s.Logger.ErrorDepthTrace(1, msg.txKey(), err)
	}
}

func (s *Server) response(conn conn, msg *Message, status, phrase string) error {
	// start line
	msg.StartLine[0] = SIPVersion
	msg.StartLine[1] = string(status)
	msg.StartLine[2] = phrase
	if msg.StartLine[2] == "" {
		msg.StartLine[2] = StatusPhrase(status)
	}
	// to tag
	if msg.Header.To.Tag == "" {
		msg.Header.To.Tag = fmt.Sprintf("%d", uid.SnowflakeID())
	}
	// via
	msg.Header.Via[0].RPort = strconv.Itoa(conn.RemotePort())
	msg.Header.Via[0].Received = conn.RemoteIP()
	//
	msg.Header.UserAgent = s.UserAgent
	//
	return s.writeMsg(conn, msg)
}

func (s *Server) writeMsg(conn conn, msg *Message) error {
	// 格式化
	var buf bytes.Buffer
	msg.Enc(&buf)
	// 发送
	return conn.write(buf.Bytes())
}

// checkTxDuration 封装代码
func (s *Server) checkTxDuration() time.Duration {
	dur := s.MsgTimeout / 4
	if dur < time.Second {
		return time.Second
	}
	return dur
}
