package schema

import "encoding/xml"

const (
	// Namespace 命名空间
	Namespace = "http://www.onvif.org/ver10/schema"
	// NamespacePrefix 命名空间前缀
	NamespacePrefix = "tt"
)

var (
	// NamespaceAttr 命名空间属性
	NamespaceAttr = NewSecurityNamespaceAttr()
)

// NewSecurityNamespaceAttr 返回命名空间属性
func NewSecurityNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:tt",
		},
		Value: Namespace,
	}
}
