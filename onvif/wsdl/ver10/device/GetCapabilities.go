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

// GetCapabilities 获取设备能力
func GetCapabilities(ctx context.Context, url string, security *soap.Security, categories ...CapabilityCategory) (*schema.Capabilities, error) {
	// 请求体
	var req soap.Envelope[*soap.Security, struct {
		XMLName  xml.Name             `xml:"tds:GetCapabilities"`
		Category []CapabilityCategory `xml:"tds:Category"`
	}]
	req.SetSoapTag()
	req.Attr = append(envelopeAttr, soap.NewSecurityNamespaceAttr())
	req.Header.Data = security
	// 不传就获取所有
	if len(categories) < 1 {
		req.Body.Data.Category = append(req.Body.Data.Category, CapabilityCategoryAll)
	}
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName      xml.Name `xml:"GetCapabilitiesResponse"`
		Capabilities schema.Capabilities
	}]
	// 发送
	err := soap.Do(ctx, url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 成功
	return &res.Body.Data.Capabilities, nil
}
