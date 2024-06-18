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
	// 保存调用链函数
	handleFunc []HandleRequestFunc
	// 当前调用的函数下标
	handleIdx int
}

// Next 执行调用链中剩下的所有函数
func (c *Request) Next() {
	for c.handleIdx < len(c.handleFunc) {
		c.handleFunc[c.handleIdx](c)
		c.handleIdx++
	}
}

// callback 执行调用链中剩下的所有函数
func (c *Request) callback() {
	for c.handleIdx < len(c.handleFunc) {
		c.handleFunc[c.handleIdx](c)
		c.handleIdx++
	}
}

// ResponseMsg 发送响应，中断调用链，msg 为 nil 不会发送数据
func (c *Request) Response(msg *Message) error {
	c.tx.finish(ErrFinish)
	c.handleIdx = len(c.handleFunc)
	if msg == nil {
		return nil
	}
	return c.tx.writeMsg(c.conn, msg)
}

// ResponseNil 什么都不返回
func (c *Request) ResponseNil() {
	c.tx.finish(ErrFinish)
	c.handleIdx = len(c.handleFunc)
}

// NewResponse 根据自身返回
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
	// 保存调用链函数
	handleFunc []HandleResponseFunc
	// 当前调用的函数下标
	handleIdx int
}

// Next 执行调用链中剩下的所有函数
func (c *Response) Next() {
	for c.handleIdx < len(c.handleFunc) {
		c.handleFunc[c.handleIdx](c)
		c.handleIdx++
	}
}

// Finish 结束调用链
func (c *Response) Finish(err error) {
	if err == nil {
		err = ErrFinish
	}
	c.tx.finish(err)
	c.handleIdx = len(c.handleFunc)
}

// callback 执行调用链中剩下的所有函数
func (c *Response) callback() {
	for c.handleIdx < len(c.handleFunc) {
		c.handleFunc[c.handleIdx](c)
		c.handleIdx++
	}
}

// Status 返回 StartLine[1]
func (c *Response) Status() string {
	return c.Message.StartLine[1]
}
