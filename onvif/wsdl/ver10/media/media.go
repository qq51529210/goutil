package media

import (
	"encoding/xml"
	"goutil/soap"
)

const (
	Namespace = "http://www.onvif.org/ver10/media/wsdl"
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
			Local: "xmlns:trt",
		},
		Value: Namespace,
	}
}
