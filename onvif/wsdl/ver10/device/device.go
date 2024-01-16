package device

import (
	"encoding/xml"
	"fmt"
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

// Device 设备
type Device struct {
	soap.Security
	url string
}

// NewDevice 返回新的设备
func NewDevice(host, username, password string) *Device {
	d := new(Device)
	d.Security.Init(username, password)
	d.url = fmt.Sprintf("http://%s/onvif/device_service", host)
	return d
}
