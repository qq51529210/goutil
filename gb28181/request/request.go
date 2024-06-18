package request

import (
	"context"
	"errors"
	"fmt"
	"goutil/gb28181/xml"
	"goutil/sip"
	"goutil/uid"
	"net"
	"time"
)

// 常量
const (
	ContentTypeSDP      = "Application/SDP"
	ContentTypeXML      = "Application/MANSCDP+xml"
	ContentTypeMANSRTSP = "Application/MANSRTSP"
	//
	MaxForwards = "70"
)

var (
	errorIPAddress = errors.New("error ip address")
)

// Request 用于组装
type Request interface {
	// 网络类型
	GetNetwork() string
	// ip
	GetIP() string
	// port
	GetPort() int
	// 本地编号
	GetLocalID() string
	// 本地域
	GetLocalDomain() string
	// 本地的可访问到的监听地址
	GetLocalAddress() string
	// 对方编号
	GetRemoteID() string
	// 对方域
	GetRemoteDomain() string
	// xml 编号
	GetXMLEncoding() string
}

// New 创建新的请求消息，channelID 用于 invite ，其他传空字符串即可
func New(ser *sip.Server, req Request, channelID, method, contentType string) (*sip.Message, net.Addr, error) {
	// 地址
	ip := net.ParseIP(req.GetIP())
	if ip == nil {
		return nil, nil, errorIPAddress
	}
	var addr net.Addr
	var proto string
	if req.GetNetwork() == "tcp" {
		addr = &net.TCPAddr{IP: ip, Port: req.GetPort()}
		proto = sip.TCP
	} else {
		addr = &net.UDPAddr{IP: ip, Port: req.GetPort()}
		proto = sip.UDP
	}
	// 消息
	m := new(sip.Message)
	fromID := req.GetLocalID()
	fromDomain := req.GetLocalDomain()
	toID := req.GetRemoteID()
	if channelID != "" {
		toID = channelID
	}
	toDomain := req.GetRemoteDomain()
	contact := req.GetLocalAddress()
	// start line
	m.StartLine[0] = method
	m.StartLine[1] = fmt.Sprintf("sip:%s@%s", toID, toDomain)
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
	m.Header.To.URI.Name = toID
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
	return m, addr, nil
}

// SendMessage 发送 message 请求并等待结果
func SendMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any) error {
	msg, addr, err := New(ser, req, body.DeviceID, sip.MethodMessage, ContentTypeXML)
	if err != nil {
		return err
	}
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}

// SendReplyMessage 发送有应答的 message 请求并等待结果
func SendReplyMessage(ctx context.Context, ser *sip.Server, req Request, body *xml.Message, data any, timeout time.Duration) error {
	msg, addr, err := New(ser, req, body.DeviceID, sip.MethodMessage, ContentTypeXML)
	if err != nil {
		return err
	}
	xml.Encode(&msg.Body, req.GetXMLEncoding(), body)
	// 应答
	rep := AddReply(body.DeviceID, body.SN, data, timeout)
	defer rep.Finish(nil)
	// 请求
	if err := ser.RequestWithContext(ctx, msg, addr, rep); err != nil {
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

// SendRegister 发送 register 请求并等待结果
func SendRegister(ctx context.Context, ser *sip.Server, req Request, expires string, data any) error {
	msg, addr, err := New(ser, req, "", sip.MethodRegister, "")
	if err != nil {
		return err
	}
	// 这两个是一样的
	msg.Header.To.URI = msg.Header.From.URI
	msg.Header.Expires = expires
	//
	return ser.RequestWithContext(ctx, msg, addr, data)
}
