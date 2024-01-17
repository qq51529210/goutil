package device

import (
	"context"
	"encoding/xml"
	"goutil/soap"
)

// Information 是 GetDeviceInformation 的返回值
type Information struct {
	Manufacturer    string
	Model           string
	FirmwareVersion string
	SerialNumber    string
	HardwareID      string
}

// GetDeviceInformation 获取设备基本信息
func (d *Device) GetDeviceInformation(ctx context.Context) (*Information, error) {
	// 请求体
	var req soap.ReqEnvelope[*soap.Security, struct {
		XMLName xml.Name `xml:"tds:GetDeviceInformation"`
	}]
	req.Attr = append(envelopeAttr, soap.NewSecurityNamespaceAttr())
	req.Header.Data = d.security
	// 响应体
	var res soap.ResEnvelope[any, struct {
		XMLName xml.Name `xml:"GetDeviceInformationResponse"`
		Information
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
	return &res.Body.Data.Information, nil
}
