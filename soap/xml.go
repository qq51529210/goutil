package soap

import (
	"encoding/xml"
	"fmt"
	"strings"
)

const (
	// Namespace 命名空间
	Namespace = "http://www.w3.org/2003/05/soap-envelope"
)

// ReqHeader 可以用于 Envelope 的 Header 字段
type ReqHeader[Data any] struct {
	XMLName xml.Name    `xml:"soap:Header"`
	Attr    []*xml.Attr `xml:",attr"`
	Data    Data
}

// ReqBody 可以用于 Envelope 的 Body 字段
type ReqBody[Data any] struct {
	XMLName xml.Name    `xml:"soap:Body"`
	Attr    []*xml.Attr `xml:",attr"`
	Data    Data
}

// ReqEnvelope 表示整个 xml 消息
type ReqEnvelope[H, B any] struct {
	XMLName xml.Name     `xml:"soap:Envelope"`
	Attr    []*xml.Attr  `xml:",attr"`
	Header  ReqHeader[H] `xml:",omitempty"`
	Body    ReqBody[B]   `xml:",omitempty"`
}

// ResHeader 可以用于 Envelope 的 Header 字段
type ResHeader[Data any] struct {
	XMLName xml.Name    `xml:"Header"`
	Attr    []*xml.Attr `xml:",attr"`
	Data    Data
}

// ResBody 可以用于 Envelope 的 Body 字段
type ResBody[Data any] struct {
	XMLName xml.Name    `xml:"Body"`
	Attr    []*xml.Attr `xml:",attr"`
	Data    Data
	Fault   *Fault `xml:",omitempty"`
}

// ResEnvelope 表示整个 xml 消息
type ResEnvelope[H, B any] struct {
	XMLName xml.Name     `xml:"Envelope"`
	Attr    []*xml.Attr  `xml:",attr"`
	Header  ResHeader[H] `xml:",omitempty"`
	Body    ResBody[B]   `xml:",omitempty"`
}

// Fault 表示 Envelope.Body 中的错误
type Fault struct {
	XMLName xml.Name     `xml:"Fault,omitempty"`
	Code    *FaultCode   `xml:"Code,omitempty"`
	Reason  *FaultReason `xml:"Reason,omitempty"`
	Detail  *FaultDetail `xml:"Detail,omitempty"`
}

func (c *Fault) Error() string {
	if c.Detail != nil {
		return fmt.Sprintf("code: %s, %s",
			c.Detail.ErrorCode, c.Detail.Description)
	}
	if c.Reason != nil {
		return c.Reason.Text
	}
	if c.Code != nil {
		var str strings.Builder
		fmt.Fprintf(&str, "code: %s", c.Code.Value)
		if c.Code.Subcode != nil {
			fmt.Fprintf(&str, " sub code: %s", c.Code.Subcode.Value)
		}
		return str.String()
	}
	return "unknown fault"
}

// FaultCode 表示 Fault 的 Code 字段
type FaultCode struct {
	Value   string     `xml:"Value"`
	Subcode *FaultCode `xml:"Subcode,omitempty"`
}

// FaultReason 表示 Fault 的 Reason 字段
type FaultReason struct {
	Text string `xml:"Text"`
}

// FaultDetail 表示 Fault 的 Detail 字段
type FaultDetail struct {
	ErrorCode   string `xml:"ErrorCode"`
	Description string `xml:"Description"`
}

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:soap",
		},
		Value: Namespace,
	}
}
