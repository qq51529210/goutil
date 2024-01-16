package device

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// CapabilityCategory 能力
type CapabilityCategory string

// 能力列表
const (
	CapabilityCategoryAll       CapabilityCategory = "All"
	CapabilityCategoryAnalytics CapabilityCategory = "Analytics"
	CapabilityCategoryDevice    CapabilityCategory = "Device"
	CapabilityCategoryEvents    CapabilityCategory = "Events"
	CapabilityCategoryImaging   CapabilityCategory = "Imaging"
	CapabilityCategoryMedia     CapabilityCategory = "Media"
	CapabilityCategoryPTZ       CapabilityCategory = "PTZ"
)

// GetCapabilities 获取设备能力，NewDevice 自动调用过一次了，all
func (d *Device) GetCapabilities(ctx context.Context, categories ...CapabilityCategory) (*schema.Capabilities, error) {
	// 请求体
	var req soap.ReqEnvelope[*soap.Security, struct {
		XMLName  xml.Name             `xml:"tds:GetCapabilities"`
		Category []CapabilityCategory `xml:"tds:Category"`
	}]
	req.Attr = append(envelopeAttr, soap.NewSecurityNamespaceAttr())
	req.Header.Data = d.security
	// 不传就获取所有
	if len(categories) < 1 {
		req.Body.Data.Category = append(req.Body.Data.Category, CapabilityCategoryAll)
	}
	// 响应体
	var res soap.ResEnvelope[any, struct {
		XMLName      xml.Name `xml:"GetCapabilitiesResponse"`
		Capabilities schema.Capabilities
	}]
	// 发送
	err := soap.Do(ctx, d.url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 错误
	if res.Body.Fault != nil {
		return nil, res.Body.Fault
	}
	// 成功
	d.capabilities = &res.Body.Data.Capabilities
	return &res.Body.Data.Capabilities, nil
}
