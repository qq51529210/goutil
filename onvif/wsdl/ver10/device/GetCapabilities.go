package device

import (
	"context"
	"goutil/soap"
	"goutil/wsdl/ver10/schema"
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

type getCapabilitiesReq struct {
	soap.Envelope[soap.Security, struct {
		XMLName  string               `xml:"tds:GetCapabilities"`
		Category []CapabilityCategory `xml:"tt:Category"`
	}]
}

type getCapabilitiesRes struct {
	soap.Envelope[any, getCapabilitiesResponse]
}

type getCapabilitiesResponse struct {
	schema.Capabilities
}

// GetCapabilities
// This method has been replaced by the more generic GetServices method.
// For capabilities of individual services refer to the GetServiceCapabilities methods.
func GetCapabilities(ctx context.Context, url, username, password string, categories ...CapabilityCategory) (*schema.Capabilities, error) {
	// 消息
	var _req getCapabilitiesReq
	_req.Envelope.Attr = envelopeAttr
	_req.Header.Data.Init(username, password)
	_req.Body.Data.Category = append(_req.Body.Data.Category, categories...)
	var _res getCapabilitiesRes
	// 请求
	err := soap.Do(ctx, url, &_req, &_res)
	if err != nil {
		return nil, err
	}
	// 错误
	if _res.Body.Fault != nil {
		return nil, _res.Body.Fault
	}
	// 成功
	return &_res.Body.Data.Capabilities, nil
}
