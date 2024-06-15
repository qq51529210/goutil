package onvif

import (
	"context"
	"errors"
	ptz "goutil/onvif/wsdl/ver20/ptz/wsdl"
	"goutil/soap"
	"time"
)

// 错误
var (
	ErrPTZCapabilityUnsupported = errors.New("ptz capability unsupported")
)

// IsPTzServiceOK 服务是否支持
func (d *Device) IsPTzServiceOK() bool {
	return d.Capabilities != nil && d.Capabilities.PTZ != nil && d.Capabilities.PTZ.XAddr != ""
}

// ContinuousMove 云台控制，速度在 0-1 之间
func (d *Device) ContinuousMove(ctx context.Context, profileToken string, xSpeed, ySpeed, zSpeed float64, timeout time.Duration) error {
	if !d.IsPTzServiceOK() {
		return ErrPTZCapabilityUnsupported
	}
	return ptz.ContinuousMove(ctx, d.Capabilities.PTZ.XAddr, soap.NewSecurity(d.username, d.password), profileToken, xSpeed, ySpeed, zSpeed, timeout)
}

// StopPTZ 停止云台控制
func (d *Device) StopPTZ(ctx context.Context, profileToken string) error {
	if !d.IsPTzServiceOK() {
		return ErrPTZCapabilityUnsupported
	}
	return ptz.Stop(ctx, d.Capabilities.PTZ.XAddr, soap.NewSecurity(d.username, d.password), profileToken, true, true)
}
