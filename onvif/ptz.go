package onvif

import (
	"context"
	"errors"
	"time"
)

// 错误
var (
	ErrPTZCapabilityUnsupported = errors.New("ptz capability unsupported")
)

// IsPTzServiceOK 服务是否支持
func (d *Device) IsPTzServiceOK() bool {
	return d.Capabilities != nil && d.Capabilities.Media != nil && d.Capabilities.Media.XAddr != ""
}

// ContinuousMove 云台控制，速度在 0-1 之间
func (d *Device) ContinuousMove(ctx context.Context, profileToken string, xSpeed, ySpeed, zSpeed float64, timeout time.Duration) error {
	if !d.IsMediaServiceOK() {
		return ErrPTZCapabilityUnsupported
	}
	//
	return nil
}
