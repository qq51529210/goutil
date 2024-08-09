package sip

import (
	"fmt"
	"goutil/uid"
	"strconv"
)

// _Context 回调函数传递的上下文
type _Context struct {
	// 事务
	tx
	// 消息
	*Message
	// 服务
	Ser *Server
	// 对端网络，tcp/udp
	RemoteNetwork string
	// 对端 IP:Port
	RemoteAddr string
	// 对端 IP
	RemoteIP string
	// 对端 Port
	RemotePort int
	// 用于保存上下文数据
	data map[any]any
}

// 实现 context.Context 接口
// 响应回调时， key=nil 返回发起请求时传入的数据
func (c *_Context) Value(key any) any {
	if c.data == nil {
		return nil
	}
	return c.data[key]
}

// SetValue 用于在调用链中传递数据
func (c *_Context) SetValue(key, value any) {
	if c.data == nil {
		c.data = make(map[any]any)
	}
	c.data[key] = value
}

// Request 请求回调上下文
type Request struct {
	_Context
	conn
	// 当前调用链
	f *reqFuncChain
}

// Next 执行调用链中剩下的所有函数
func (c *Request) Next() {
	c.f.Next(c)
}

// IsResponsed 返回是否已经响应
func (c *Request) IsResponsed() bool {
	return c.f.IsAbort()
}

// Change 将当前的中断并改为执行新的回调
func (c *Request) Change(funcs ...HandleRequestFunc) {
	c.f.Abort()
	f := new(reqFuncChain)
	f.f = append(f.f, funcs...)
	c.f = f
	f.handle(c)
}

// ResponseMsg 发送响应，中断调用链，msg 为 nil 不会发送数据
func (c *Request) Response(msg *Message) error {
	// 通知
	c.tx.finish(ErrFinish)
	// 中断
	c.f.Abort()
	//
	if msg == nil {
		return nil
	}
	// 日志
	c.Ser.logger.DebugfTrace(c.ID(), "response to %s %s\n%v", c.RemoteNetwork, c.RemoteAddr, msg)
	// 发送
	return c.tx.writeMsg(c.conn, msg)
}

// NewResponse 根据自身的字段构造响应消息
func (c *Request) NewResponse(status, phrase string) *Message {
	// 消息
	msg := &Message{}
	msg.SetResponseStartLine(status, phrase)
	msg.Header = c.Message.Header
	msg.Header.To.Tag = fmt.Sprintf("%d", uid.SnowflakeID())
	// via
	msg.Header.Via[0].RPort = strconv.Itoa(c.RemotePort)
	msg.Header.Via[0].Received = c.RemoteIP
	//
	msg.Header.UserAgent = c.Ser.userAgent
	//
	return msg
}

// Response 响应回调上下文
type Response struct {
	_Context
	conn
	// 发起请求传入的数据
	ReqData any
	// 当前调用链
	f *resFuncChain
}

// Next 执行调用链中剩下的所有函数
func (c *Response) Next() {
	c.f.Next(c)
}

// Change 将当前的中断并改为执行新的回调
func (c *Response) Change(funcs ...HandleResponseFunc) {
	c.f.Abort()
	f := new(resFuncChain)
	f.f = append(f.f, funcs...)
	c.f = f
}

// Finish 结束调用链
func (c *Response) Finish(err error) {
	if err == nil {
		err = ErrFinish
	}
	c.tx.finish(err)
	c.f.Abort()
}

// Status 返回 StartLine[1]
func (c *Response) Status() string {
	return c.Message.StartLine[1]
}

// Phrase 返回 StartLine[2]
func (c *Response) Phrase() string {
	return c.Message.StartLine[2]
}
