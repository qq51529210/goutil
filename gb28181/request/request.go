package request

import (
	"context"
	"fmt"
	"goutil/gb28181/xml"
	"goutil/sip"
	"goutil/uid"
	"io"
	"net"
)

// 常量
const (
	ContentTypeSDP      = "Application/SDP"
	ContentTypeXML      = "Application/MANSCDP+xml"
	ContentTypeMANSRTSP = "Application/MANSRTSP"
	//
	MaxForwards = "70"
)

// Invite 用于 bye/info 消息
type Invite interface {
	GetFromTag() string
	GetToTag() string
	GetCallID() string
}

// Request 用于组装 sip.Message
type Request interface {
	// 网络类型
	GetNetwork() string
	// 网络地址
	GetNetAddr() (net.Addr, error)
	// 本地国标编号
	GetFromID() string
	// 本地国标域名
	GetFromDomain() string
	// 对方国标编号
	GetToID() string
	// 对方国标域名
	GetToDomain() string
	// Contact 的地址，一般是 ip:port
	GetContactAddress() string
	// xml 编码
	GetXMLEncoding() string
}

// New 创建新的请求消息，toID 用于自定义，传空字符串使用 req.GetToID()
// toID 覆盖 req.GetToID() 的值
// tcpProto 是 tcp/udp
func New(req Request, toID, method, contentType string) *sip.Message {
	// 消息
	m := new(sip.Message)
	// from
	fromID := req.GetFromID()
	fromDomain := req.GetFromDomain()
	// to
	_toID := req.GetToID()
	// 如果指定，使用指定
	if toID != "" {
		_toID = toID
	}
	toDomain := req.GetToDomain()
	contact := req.GetContactAddress()
	// start line
	m.StartLine[0] = method
	m.StartLine[1] = fmt.Sprintf("sip:%s@%s", _toID, toDomain)
	m.StartLine[2] = sip.SIPVersion
	// via
	var proto string
	if req.GetNetwork() == "tcp" {
		proto = sip.TCP
	} else {
		proto = sip.UDP
	}
	m.Header.Via = append(m.Header.Via, &sip.Via{
		Proto:   proto,
		Address: contact,
		Branch:  fmt.Sprintf("%s%d", sip.BranchPrefix, uid.SnowflakeID()),
	})
	// From
	m.Header.From.URI.Scheme = sip.SIP
	m.Header.From.URI.Name = fromID
	m.Header.From.URI.Domain = fromDomain
	m.Header.From.Tag = uid.SnowflakeIDString()
	// To
	m.Header.To.URI.Scheme = sip.SIP
	m.Header.To.URI.Name = _toID
	m.Header.To.URI.Domain = toDomain
	// Call-ID
	m.Header.CallID = uid.SnowflakeIDString()
	// CSeq
	m.Header.CSeq.SN = sip.GetSNString()
	m.Header.CSeq.Method = method
	// Max-Forwards
	m.Header.MaxForwards = MaxForwards
	// Content-Type
	m.Header.ContentType = contentType
	// Contact
	m.Header.Contact.Scheme = sip.SIP
	m.Header.Contact.Name = fromID
	m.Header.Contact.Domain = contact
	//
	return m
}

// NewBye 返回新的 bye 方法的请求
func NewBye(req Request, channelID string, invite Invite) *sip.Message {
	// 消息
	msg := New(req, channelID, sip.MethodBye, "")
	// 重新设置头
	msg.Header.From.Tag = invite.GetFromTag()
	msg.Header.To.Tag = invite.GetToTag()
	msg.Header.CallID = invite.GetCallID()
	//
	return msg
}

// SendBye 封装请求
func SendBye(ctx context.Context, ser *sip.Server, req Request, channelID string, invite Invite, data any) error {
	// 地址
	addr, err := req.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := NewBye(req, channelID, invite)
	// 发送
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// NewInfo 返回新的 info 方法的请求
func NewInfo(req Request, channelID string, invite Invite) *sip.Message {
	// 消息
	msg := New(req, channelID, sip.MethodInfo, ContentTypeMANSRTSP)
	// 重新设置头
	msg.Header.From.Tag = invite.GetFromTag()
	msg.Header.To.Tag = invite.GetToTag()
	msg.Header.CallID = invite.GetCallID()
	//
	return msg
}

// SendInfo 封装请求
func SendInfo(ctx context.Context, ser *sip.Server, req Request, channelID string, invite Invite, body io.Reader, data any) error {
	// 地址
	addr, err := req.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := NewInfo(req, channelID, invite)
	// body
	if _, err := io.Copy(&msg.Body, body); err != nil {
		return err
	}
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// NewInvite 返回新的 invite 方法的请求
func NewInvite(req Request, channelID string) *sip.Message {
	return New(req, channelID, sip.MethodInvite, ContentTypeSDP)
}

// NewAck 返回新的 ack 方法的请求
func NewAck(req Request, channelID string, invite Invite) *sip.Message {
	// 消息
	msg := New(req, channelID, sip.MethodACK, "")
	// 重新设置头
	msg.Header.From.Tag = invite.GetFromTag()
	msg.Header.To.Tag = invite.GetToTag()
	msg.Header.CallID = invite.GetCallID()
	//
	return msg
}

// NewRegister 返回新的 register 方法的请求
func NewRegister(req Request, expires string) *sip.Message {
	// 消息
	msg := New(req, "", sip.MethodRegister, "")
	// 重新设置头 from 和 to 一样
	msg.Header.To.URI = msg.Header.From.URI
	msg.Header.Expires = expires
	//
	return msg
}

// NewMessage 返回新的 message 方法的请求
func NewMessage(req Request, body *xml.Message) *sip.Message {
	// 消息
	msg := New(req, "", sip.MethodMessage, ContentTypeXML)
	// body
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	return msg
}

// SendMessage 发送 message 请求并等待结果
func SendMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any) error {
	// 地址
	addr, err := req.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := NewMessage(req, body)
	// 发送
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// SendReplyMessage 发送有应答的 message 请求并等待结果
func SendReplyMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any) error {
	// 应答
	rep := AddReply(body.DeviceID, body.SN, data, ser.MsgTimeout())
	// 请求
	if err := SendMessage(ctx, ser, req, body, rep); err != nil {
		rep.Finish(err, nil)
		return err
	}
	// 等待响应请求结果
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-rep.Done():
		return rep.Err()
	}
}

// NewSubscribe 返回新的 subscribe 方法的请求
func NewSubscribe(req Request, body *xml.Subscribe) *sip.Message {
	// 消息
	msg := New(req, "", sip.MethodSubscribe, ContentTypeXML)
	// body
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	//
	return msg
}

// SendSubscribe 封装请求
func SendSubscribe(ctx context.Context, ser *sip.Server, req Request, body *xml.Subscribe, expire int64, data any) error {
	// 地址
	addr, err := req.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := NewSubscribe(req, body)
	// 特别的
	msg.Header.Expires = fmt.Sprintf("%d", expire)
	msg.Header.Set("Event", "presence")
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}
