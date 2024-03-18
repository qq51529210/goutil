package ptz

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// StopPanTilt 停止水平
func StopPanTilt(ctx context.Context, url string, security *soap.Security, profileToken string) error {
	return stop(ctx, url, security, profileToken, true, false)
}

// StopZoom 停止纵深
func StopZoom(ctx context.Context, url string, security *soap.Security, profileToken string) error {
	return stop(ctx, url, security, profileToken, false, true)
}

// StopPanTiltZoom 停止所有
func StopPanTiltZoom(ctx context.Context, url string, security *soap.Security, profileToken string) error {
	return stop(ctx, url, security, profileToken, true, true)
}

// stop 停止
func stop(ctx context.Context, url string, security *soap.Security, profileToken string, panTilt, zoom bool) error {
	// 请求体
	var req soap.Envelope[*soap.Security, struct {
		XMLName      xml.Name `xml:"tptz:Stop"`
		ProfileToken string   `xml:"tptz:ReferenceToken"`
		PanTilt      bool     `xml:"tptz:PanTilt"`
		Zoom         bool     `xml:"tptz:Zoom"`
	}]
	req.SetSoapTag()
	req.Attr = append(envelopeAttr, soap.NewSecurityNamespaceAttr())
	req.Header.Data = security
	req.Body.Data.ProfileToken = profileToken
	req.Body.Data.PanTilt = panTilt
	req.Body.Data.Zoom = zoom
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName      xml.Name `xml:"StopResponse"`
		Capabilities schema.Capabilities
	}]
	// 发送
	return soap.Do(ctx, url, &req, &res)
}
