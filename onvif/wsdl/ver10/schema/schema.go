package schema

import "encoding/xml"

const (
	// Namespace 命名空间
	Namespace = "http://www.onvif.org/ver10/schema"
	// NamespacePrefix 命名空间前缀
	NamespacePrefix = "tt"
)

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:tt",
		},
		Value: Namespace,
	}
}
