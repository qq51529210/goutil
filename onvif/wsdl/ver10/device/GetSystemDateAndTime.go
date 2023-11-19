package device

import (
	"context"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

type getSystemDateAndTimeReq struct {
	soap.Envelope[any, struct {
		XMLName string `xml:"tt:GetSystemDateAndTime"`
	}]
}

type getSystemDateAndTimeRes struct {
	soap.Envelope[any, getSystemDateAndTimeResponse]
}

type getSystemDateAndTimeResponse struct {
	schema.SystemDateTime
}

// GetSystemDateAndTime
// This operation gets the device system date and time. The device shall support the return of
// the daylight saving setting and of the manual system date and time (if applicable) or indication
// of NTP time (if applicable) through the GetSystemDateAndTime command.
// A device shall provide the UTCDateTime information.
func GetSystemDateAndTime(ctx context.Context, url string) (*schema.SystemDateTime, error) {
	// 消息
	var _req getSystemDateAndTimeReq
	_req.Envelope.Attr = envelopeAttr
	var _res getSystemDateAndTimeRes
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
	return &_res.Body.Data.SystemDateTime, nil
}
