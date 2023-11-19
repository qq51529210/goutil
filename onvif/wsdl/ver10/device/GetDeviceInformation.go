package device

import (
	"context"
	"goutil/soap"
)

type getDeviceInformationReq struct {
	soap.Envelope[any, struct {
		XMLName string `xml:"tds:GetDeviceInformation"`
	}]
}

type getDeviceInformationRes struct {
	soap.Envelope[any, DeviceInformation]
}

// DeviceInformation -> GetDeviceInformationResponse
type DeviceInformation struct {
	Manufacturer    string
	Model           string
	FirmwareVersion string
	SerialNumber    string
	HardwareID      string
}

// GetDeviceInformation
// This operation gets basic device information from the device.
func GetDeviceInformation(ctx context.Context, url string) (*DeviceInformation, error) {
	// 消息
	var _req getDeviceInformationReq
	_req.Envelope.Attr = envelopeAttr
	var _res getDeviceInformationRes
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
	return &_res.Body.Data, nil
}
