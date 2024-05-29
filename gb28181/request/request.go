package request

import (
	"context"
	"errors"
	"fmt"
	"goutil/gb28181/xml"
	"goutil/sip"
	"goutil/uid"
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

var (
	errorIPAddress = errors.New("error ip address")
)

type request struct {
	network     string
	ip          string
	port        int
	method      string
	contentType string
	fromID      string
	fromDomain  string
	fromTag     string
	toID        string
	toDomain    string
	toTag       string
	callID      string
	contact     string
}

func (r *request) New() (*sip.Request, net.Addr, error) {
	ip := net.ParseIP(r.ip)
	if ip == nil {
		return nil, nil, errorIPAddress
	}
	// 地址
	var addr net.Addr
	var proto string
	if r.network == "tcp" {
		addr = &net.TCPAddr{
			IP:   net.ParseIP(r.ip),
			Port: r.port,
		}
		proto = sip.TCP
	} else {
		addr = &net.UDPAddr{
			IP:   net.ParseIP(r.ip),
			Port: r.port,
		}
		proto = sip.UDP
	}
	// 请求
	m := sip.NewRequest()
	// start line
	m.StartLine[0] = r.method
	m.StartLine[1] = fmt.Sprintf("sip:%s@%s", r.toID, r.toDomain)
	m.StartLine[2] = sip.SIPVersion
	// via
	m.Header.Via = append(m.Header.Via, &sip.Via{
		Proto:   proto,
		Address: r.contact,
		Branch:  fmt.Sprintf("%s%d", sip.BranchPrefix, uid.SnowflakeID()),
	})
	// From
	m.Header.From.URI.Scheme = sip.SIP
	m.Header.From.URI.Name = r.fromID
	m.Header.From.URI.Domain = r.fromDomain
	m.Header.From.Tag = r.fromTag
	if m.Header.From.Tag == "" {
		m.Header.From.Tag = uid.SnowflakeIDString()
	}
	// To
	m.Header.To.URI.Scheme = sip.SIP
	m.Header.To.URI.Name = r.toID
	m.Header.To.URI.Domain = r.toDomain
	m.Header.To.Tag = r.toTag
	// Call-ID
	m.Header.CallID = r.callID
	if m.Header.CallID == "" {
		m.Header.CallID = uid.SnowflakeIDString()
	}
	// CSeq
	m.Header.CSeq.SN = sip.GetSNString()
	m.Header.CSeq.Method = r.method
	// Max-Forwards
	m.Header.MaxForwards = MaxForwards
	// Content-Type
	m.Header.ContentType = r.contentType
	// Contact
	m.Header.Contact.Scheme = sip.SIP
	m.Header.Contact.Name = r.fromID
	m.Header.Contact.Domain = r.contact
	//
	return m, addr, nil
}

// Device 用于 NewDeviceRequest
type Device interface {
	GetNetwork() string
	GetIP() string
	GetPort() int
	GetDeviceID() string
	GetDeviceDomain() string
	GetServerID() string
	GetServerDomain() string
	GetContact() string
	GetXMLEncoding() string
}

// NewDeviceRequest 创建新的设备请求
func NewDeviceRequest(ser *sip.Server, device Device, deviceOrChannelID, method, contentType, fromTag, toTag, callID string) (*sip.Request, net.Addr, error) {
	var req request
	req.network = device.GetNetwork()
	req.ip = device.GetIP()
	req.port = device.GetPort()
	req.callID = callID
	req.contentType = contentType
	req.method = method
	req.toID = deviceOrChannelID
	req.toDomain = device.GetDeviceDomain()
	req.toTag = toTag
	req.fromID = device.GetServerID()
	req.fromDomain = device.GetServerDomain()
	req.fromTag = fromTag
	req.contact = device.GetContact()
	return req.New()
}

// SendDeviceMessageRequest 向设备发送 message 类型的请求并等待结果
func SendDeviceMessageRequest(ctx context.Context, ser *sip.Server, device Device, body *xml.Message, ctxData any) error {
	m, addr, err := NewDeviceRequest(ser, device, body.DeviceID, sip.MethodMessage, ContentTypeXML, "", "", "")
	if err != nil {
		return err
	}
	xml.Encode(&m.Body, device.GetXMLEncoding(), body)
	return ser.Request(ctx, m, addr, ctxData)
}

// SendDeviceMessageReplyRequest 向设备发送 message 类型的应答式请求并等待结果
func SendDeviceMessageReplyRequest(ctx context.Context, ser *sip.Server, device Device, body *xml.Message, ctxData any) error {
	// 消息
	m, addr, err := NewDeviceRequest(ser, device, body.DeviceID, sip.MethodMessage, ContentTypeXML, "", "", "")
	if err != nil {
		return err
	}
	xml.Encode(&m.Body, device.GetXMLEncoding(), body)
	// 有响应的请求
	rep := AddReply(device.GetDeviceID(), body.SN, ctxData, ser.GetTxTimeout())
	// 请求
	if err := ser.Request(ctx, m, addr, nil); err != nil {
		rep.Finish(err)
		return err
	}
	// 等待结果
	select {
	case <-ctx.Done():
		err := ctx.Err()
		rep.Finish(err)
		return err
	case <-rep.Done():
		return rep.Err()
	}
}

// Device 用于 NewDeviceRequest
type Cascade interface {
	GetNetwork() string
	GetIP() string
	GetPort() int
	GetLocalID() string
	GetLocalDomain() string
	GetCascadeID() string
	GetCascadeDomain() string
	GetContact() string
	GetXMLEncoding() string
}

// NewCascadeRequest 创建新的级联请求
func NewCascadeRequest(ser *sip.Server, cascade Cascade, cascadeOrChannelID, method, contentType, fromTag, toTag, callID string) (*sip.Request, net.Addr, error) {
	var req request
	req.network = cascade.GetNetwork()
	req.ip = cascade.GetIP()
	req.port = cascade.GetPort()
	req.callID = callID
	req.contentType = contentType
	req.method = method
	req.toID = cascadeOrChannelID
	req.toDomain = cascade.GetCascadeDomain()
	req.toTag = toTag
	req.fromID = cascade.GetLocalID()
	req.fromDomain = cascade.GetLocalDomain()
	req.fromTag = fromTag
	req.contact = cascade.GetContact()
	return req.New()
}

// SendCascadeMessageRequest 向级联发送 message 类型的请求并等待结果
func SendCascadeMessageRequest(ctx context.Context, ser *sip.Server, cascade Cascade, body *xml.Message, ctxData any) error {
	m, addr, err := NewCascadeRequest(ser, cascade, cascade.GetCascadeID(), sip.MethodMessage, ContentTypeXML, "", "", "")
	if err != nil {
		return err
	}
	xml.Encode(&m.Body, cascade.GetXMLEncoding(), body)
	return ser.Request(ctx, m, addr, ctxData)
}

// SendCascadeRegisterRequest 向级联发送 register 类型的请求并等待结果
func SendCascadeRegisterRequest(ctx context.Context, ser *sip.Server, cascade Cascade, body *xml.Message, ctxData any) error {
	m, addr, err := NewCascadeRequest(ser, cascade, body.DeviceID, sip.MethodMessage, ContentTypeXML, "", "", "")
	if err != nil {
		return err
	}
	xml.Encode(&m.Body, cascade.GetXMLEncoding(), body)
	// 这两个是一样的
	m.Header.To.URI = m.Header.From.URI
	return ser.Request(ctx, m, addr, ctxData)
}
