package onvif

import (
	"context"
	owvm "goutil/onvif/wsdl/ver10/media"
	"goutil/soap"
	"net/url"
)

// IsMediaServiceOK 服务是否支持
func (d *Device) IsMediaServiceOK() bool {
	return d.Capabilities != nil && d.Capabilities.Media != nil && d.Capabilities.Media.XAddr != ""
}

// GetRTSPStreamURL 返回 rtsp 地址
func (d *Device) GetRTSPStreamURL(ctx context.Context, profileToken string) (string, error) {
	if !d.IsMediaServiceOK() {
		return "", ErrMediaCapabilityUnsupported
	}
	m, err := owvm.GetStreamURL(ctx, d.Capabilities.Media.XAddr, soap.NewSecurity(d.username, d.password),
		profileToken, owvm.StreamProtocolRTSP, owvm.StreamTypeRTPUnicast)
	if err != nil {
		return "", err
	}
	u, _ := url.Parse(m.URL)
	u.User = url.UserPassword(d.username, d.password)
	u.Host = d.replaceIP(u.Host)
	return u.String(), nil
}
