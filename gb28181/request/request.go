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

// Request 用于组装
type Request interface {
	// tcp/udp
	GetNetwork() string
	GetIP() net.IP
	GetPort() int
	GetFromID() string
	GetFromDomain() string
	GetContactAddress() string
	GetToID() string
	GetToDomain() string
	GetXMLEncoding() string
}

// New 创建新的请求消息，toID 用于自定义，传空字符串使用 req.GetToID()
func New(req Request, toID, method, contentType string) (*sip.Message, net.Addr) {
	// 地址
	var addr net.Addr
	var proto string
	if req.GetNetwork() == "tcp" {
		addr = &net.TCPAddr{IP: req.GetIP(), Port: req.GetPort()}
		proto = sip.TCP
	} else {
		addr = &net.UDPAddr{IP: req.GetIP(), Port: req.GetPort()}
		proto = sip.UDP
	}
	// 消息
	m := new(sip.Message)
	fromID := req.GetFromID()
	fromDomain := req.GetFromDomain()
	_toID := req.GetToID()
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
	return m, addr
}

// NewBye 返回新的 bye 方法的请求
func NewBye(req Request, channelID string, invite Invite) (*sip.Message, net.Addr) {
	msg, addr := New(req, channelID, sip.MethodBye, "")
	//
	msg.Header.From.Tag = invite.GetFromTag()
	msg.Header.To.Tag = invite.GetToTag()
	msg.Header.CallID = invite.GetCallID()
	//
	return msg, addr
}

// SendBye 封装请求
func SendBye(ctx context.Context, ser *sip.Server, req Request, channelID string, invite Invite, data any) error {
	// 消息
	msg, addr := NewBye(req, channelID, invite)
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// NewInfo 返回新的 info 方法的请求
func NewInfo(req Request, channelID string, invite Invite) (*sip.Message, net.Addr) {
	msg, addr := New(req, channelID, sip.MethodInfo, ContentTypeMANSRTSP)
	//
	msg.Header.From.Tag = invite.GetFromTag()
	msg.Header.To.Tag = invite.GetToTag()
	msg.Header.CallID = invite.GetCallID()
	//
	return msg, addr
}

// SendInfo 封装请求
func SendInfo(ctx context.Context, ser *sip.Server, req Request, channelID string, invite Invite, body io.Reader, data any) error {
	// 消息
	msg, addr := NewInfo(req, channelID, invite)
	// body
	if _, err := io.Copy(&msg.Body, body); err != nil {
		return err
	}
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// NewInvite 返回新的 invite 方法的请求
func NewInvite(req Request, channelID string) (*sip.Message, net.Addr) {
	return New(req, channelID, sip.MethodInvite, ContentTypeSDP)
}

// NewRegister 返回新的 register 方法的请求
func NewRegister(req Request, expires string) (*sip.Message, net.Addr) {
	msg, addr := New(req, "", sip.MethodRegister, "")
	// from 和 to 一样
	msg.Header.To.URI = msg.Header.From.URI
	msg.Header.Expires = expires
	//
	return msg, addr
}

// NewMessage 返回新的 message 方法的请求
func NewMessage(req Request, body *xml.Message) (*sip.Message, net.Addr) {
	msg, addr := New(req, "", sip.MethodMessage, ContentTypeXML)
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	return msg, addr
}

// SendMessage 发送 message 请求并等待结果
func SendMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any) error {
	msg, addr := NewMessage(req, body)
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// SendReplyMessage 发送有应答的 message 请求并等待结果
func SendReplyMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any) error {
	// 应答
	rep := AddReply(body.DeviceID, body.SN, data, ser.MsgTimeout())
	defer rep.Finish(nil)
	// 请求
	if err := SendMessage(ctx, ser, req, body, rep); err != nil {
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
func NewSubscribe(req Request, body *xml.Subscribe) (*sip.Message, net.Addr) {
	msg, addr := New(req, "", sip.MethodSubscribe, ContentTypeXML)
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	return msg, addr
}

// SendSubscribe 封装请求
func SendSubscribe(ctx context.Context, ser *sip.Server, req Request, body *xml.Subscribe, expire int64, data any) error {
	// 消息
	msg, addr := NewSubscribe(req, body)
	//
	msg.Header.Expires = fmt.Sprintf("%d", expire)
	msg.Header.Set("Event", "presence")
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}
