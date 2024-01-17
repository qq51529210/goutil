package device

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// GetSystemDateAndTime 获取时间日期
func GetSystemDateAndTime(ctx context.Context, url string) (*schema.SystemDateTime, error) {
	// 请求体
	var req soap.Envelope[any, struct {
		XMLName xml.Name `xml:"tds:GetSystemDateAndTime"`
	}]
	req.SetSoapTag()
	req.Attr = envelopeAttr
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName           xml.Name `xml:"GetSystemDateAndTimeResponse"`
		SystemDateAndTime schema.SystemDateTime
	}]
	// 发送
	err := soap.Do(ctx, url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 错误
	if res.Body.Fault != nil {
		return nil, res.Body.Fault
	}
	// 成功
	return &res.Body.Data.SystemDateAndTime, nil
}
