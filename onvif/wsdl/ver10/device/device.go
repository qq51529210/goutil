package device

import (
	"context"
	"encoding/xml"
	"fmt"
	"goutil/onvif/wsdl/ver10/schema"
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
	url          string
	security     *soap.Security
	capabilities *schema.Capabilities
}

// NewDevice 返回新的设备，同时获取 all 设备能力
func NewDevice(ctx context.Context, host, username, password string) (*Device, error) {
	// 初始化
	d := new(Device)
	d.url = fmt.Sprintf("http://%s/onvif/device_service", host)
	d.security = new(soap.Security)
	d.security.Init(username, password)
	// 获取设备能力
	capabilities, err := d.GetCapabilities(ctx, CapabilityCategoryAll)
	if err != nil {
		return nil, err
	}
	d.capabilities = capabilities
	// 返回
	return d, nil
}
