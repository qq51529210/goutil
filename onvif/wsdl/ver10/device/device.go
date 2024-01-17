package device

import (
	"encoding/xml"
	"goutil/soap"
)

const (
	// Namespace 命名空间
	Namespace = "http://www.onvif.org/ver10/device/wsdl"
)

var (
	envelopeAttr = []*xml.Attr{
		soap.NewNamespaceAttr(),
		NewNamespaceAttr(),
	}
)

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:tds",
		},
		Value: Namespace,
	}
}
