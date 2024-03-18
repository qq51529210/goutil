package ptz

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/onvif/wsdl/xs"
	"goutil/soap"
	"time"
)

// ContinuousMove 云台控制继续移动
func ContinuousMove(ctx context.Context, url string, security *soap.Security, profileToken string, xSpeed, ySpeed, zSpeed float64, timeout time.Duration) error {
	// 请求体
	var req soap.Envelope[*soap.Security, struct {
		XMLName      xml.Name        `xml:"tptz:ContinuousMove"`
		ProfileToken string          `xml:"tptz:ReferenceToken"`
		Velocity     schema.PTZSpeed `xml:"tptz:Velocity"`
		Timeout      string          `xml:"tptz:Timeout"`
	}]
	req.SetSoapTag()
	req.Attr = append(envelopeAttr, soap.NewSecurityNamespaceAttr())
	req.Attr = append(req.Attr, xs.NewNamespaceAttr())
	req.Header.Data = security
	req.Body.Data.ProfileToken = profileToken
	req.Body.Data.Velocity = schema.PTZSpeed{
		PanTilt: schema.Vector2D{
			X:     xSpeed,
			Y:     ySpeed,
			Space: PanTiltSpaceVelocityGenericSpace,
		},
		Zoom: schema.Vector1D{
			X:     zSpeed,
			Space: ZoomSpacesVelocityGenericSpace,
		},
	}
	req.Body.Data.Timeout = xs.Duration(timeout).String()
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName      xml.Name `xml:"ContinuousMoveResponse"`
		Capabilities schema.Capabilities
	}]
	// 发送
	return soap.Do(ctx, url, &req, &res)
}
