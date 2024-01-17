package media

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// GetProfiles 查询媒体属性
func GetProfiles(ctx context.Context, url string, security *soap.Security) ([]*schema.Profile, error) {
	// 请求体
	var req soap.Envelope[*soap.Security, struct {
		XMLName xml.Name `xml:"trt:GetProfiles"`
	}]
	req.SetSoapTag()
	req.Attr = append(envelopeAttr, soap.NamespaceAttr)
	req.Header.Data = security
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName  xml.Name `xml:"GetProfilesResponse"`
		Profiles []*schema.Profile
	}]
	// 发送
	err := soap.Do(ctx, url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 成功
	return res.Body.Data.Profiles, nil
}
