package sip

import (
	"context"
	"goutil/log"
	"net"
	"time"
)

type ServerOption struct {
	// 用户代理
	UserAgent string
	// 日志
	Logger *log.Logger
	// 消息最大的字节数，接收到的消息如果大于这个数会被丢弃
	MaxMessageLen int
	// 发起请求的超时时间，或者响应消息缓存的超时时间
	MsgTimeout time.Duration
	// tcp conn 最大空闲时间，就是 read timeout
	TCPMaxIdleTime time.Duration
}

type Server struct {
	// 用户代理
	userAgent string
	// 日志
	logger *log.Logger
	// 消息最大的字节数，接收到的消息如果大于这个数会被丢弃
	maxMessageLen int
	// 发起请求的超时时间，或者响应消息缓存的超时时间
	msgTimeout time.Duration
	// 回调函数
	handleFunc
	// udp 服务
	udp udpServer
	// tcp 服务
	tcp tcpServer
}

// NewServer 必须用这个创建
func NewServer(opt *ServerOption) *Server {
	return &Server{
		userAgent:     opt.UserAgent,
		logger:        opt.Logger,
		maxMessageLen: opt.MaxMessageLen,
		msgTimeout:    opt.MsgTimeout,
	}
}

// MsgTimeout 返回底层判断消息超时的时间
func (s *Server) MsgTimeout() time.Duration {
	return s.msgTimeout
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
func (s *Server) ServeTCP(address string, maxIdleTime time.Duration) error {
	s.tcp.s = s
	s.tcp.maxIdleTime = maxIdleTime
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

// RequestAbort 主动中断请求
func (s *Server) RequestAbort(network, msgTxKey string, err error) {
	if err == nil {
		err = ErrFinish
	}
	if network == "udp" {
		t := s.udp.activeTx.Get(msgTxKey)
		if t != nil {
			t.finish(err)
		}
	} else {
		t := s.tcp.activeTx.Get(msgTxKey)
		if t != nil {
			t.finish(err)
		}
	}
}

// checkTxDuration 封装代码
func (s *Server) checkTxDuration() time.Duration {
	dur := s.msgTimeout / 4
	if dur < time.Second {
		return time.Second
	}
	return dur
}
