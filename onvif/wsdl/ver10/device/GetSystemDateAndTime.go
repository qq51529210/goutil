package device

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// GetSystemDateAndTime 获取时间日期
func (d *Device) GetSystemDateAndTime(ctx context.Context) (*schema.SystemDateTime, error) {
	// 请求体
	var req soap.ReqEnvelope[any, struct {
		XMLName xml.Name `xml:"tds:GetSystemDateAndTime"`
	}]
	req.Attr = envelopeAttr
	// 响应体
	var res soap.ResEnvelope[any, struct {
		XMLName           xml.Name `xml:"GetSystemDateAndTimeResponse"`
		SystemDateAndTime schema.SystemDateTime
	}]
	// 请求
	err := soap.Do(ctx, d.url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 错误
	if res.Body.Fault != nil {
		return nil, res.Body.Fault
	}
	// 成功
	// return nil, nil
	return &res.Body.Data.SystemDateAndTime, nil
}
