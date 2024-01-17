package onvif

import (
	"context"
	"errors"
	"fmt"
	"goutil/soap"
	"net/url"

	owvd "goutil/onvif/wsdl/ver10/device"
	owvm "goutil/onvif/wsdl/ver10/media"
	owvs "goutil/onvif/wsdl/ver10/schema"
)

// 错误
var (
	ErrMediaCapabilityUnsupported = errors.New("media capability unsupported")
)

// Device 表示设备
type Device[Data any] struct {
	URL string
	*soap.Security
	// 需要主动调用赋值
	*owvs.Capabilities
	// 自定义数据
	Data Data
	// 设备地址
	host string
}

// NewDevice 初始化字段
func NewDevice[Data any](host, username, password string, data Data) *Device[Data] {
	// 初始化
	d := new(Device[Data])
	d.URL = fmt.Sprintf("http://%s/onvif/device_service", host)
	d.Security = new(soap.Security)
	d.Security.Init(username, password)
	d.Data = data
	d.host = host
	// 返回
	return d
}

// NewDeviceWithCapabilities 初始化字段，同时获取能力
func NewDeviceWithCapabilities[Data any](ctx context.Context, host, username, password string, data Data) (*Device[Data], error) {
	// 初始化
	d := new(Device[Data])
	d.URL = fmt.Sprintf("http://%s/onvif/device_service", host)
	d.Security = new(soap.Security)
	d.Security.Init(username, password)
	d.Data = data
	d.host = host
	// 获取能力
	c, err := d.GetCapabilities(ctx)
	if err != nil {
		return nil, err
	}
	// 转换一下 host
	if c.Analytics != nil {
		c.Analytics.XAddr = d.ReplaceXAddrHost(c.Analytics.XAddr)
	}
	if c.Device != nil {
		c.Device.XAddr = d.ReplaceXAddrHost(c.Device.XAddr)
	}
	if c.Events != nil {
		c.Events.XAddr = d.ReplaceXAddrHost(c.Events.XAddr)
	}
	if c.Imaging != nil {
		c.Imaging.XAddr = d.ReplaceXAddrHost(c.Imaging.XAddr)
	}
	if c.Media != nil {
		c.Media.XAddr = d.ReplaceXAddrHost(c.Media.XAddr)
	}
	if c.PTZ != nil {
		c.PTZ.XAddr = d.ReplaceXAddrHost(c.PTZ.XAddr)
	}
	if c.Extension != nil {
		if c.Extension.DeviceIO != nil {
			c.Extension.DeviceIO.XAddr = d.ReplaceXAddrHost(c.Extension.DeviceIO.XAddr)
		}
		if c.Extension.Display != nil {
			c.Extension.Display.XAddr = d.ReplaceXAddrHost(c.Extension.Display.XAddr)
		}
		if c.Extension.Recording != nil {
			c.Extension.Recording.XAddr = d.ReplaceXAddrHost(c.Extension.Recording.XAddr)
		}
		if c.Extension.Search != nil {
			c.Extension.Search.XAddr = d.ReplaceXAddrHost(c.Extension.Search.XAddr)
		}
		if c.Extension.Replay != nil {
			c.Extension.Replay.XAddr = d.ReplaceXAddrHost(c.Extension.Replay.XAddr)
		}
		if c.Extension.Receiver != nil {
			c.Extension.Receiver.XAddr = d.ReplaceXAddrHost(c.Extension.Receiver.XAddr)
		}
		if c.Extension.AnalyticsDevice != nil {
			c.Extension.AnalyticsDevice.XAddr = d.ReplaceXAddrHost(c.Extension.AnalyticsDevice.XAddr)
		}
	}
	d.Capabilities = c
	// 返回
	return d, nil
}

// ReplaceXAddrHost 替换掉 xaddr url 中的 host
func (d *Device[Data]) ReplaceXAddrHost(xaddr string) string {
	u, _ := url.Parse(xaddr)
	u.Host = d.host
	return u.String()
}

// GetSystemDateAndTime 查询日期时间
func (d *Device[Data]) GetSystemDateAndTime(ctx context.Context) (*owvs.SystemDateTime, error) {
	return owvd.GetSystemDateAndTime(ctx, d.URL)
}

// GetCapabilities 查询能力
func (d *Device[Data]) GetCapabilities(ctx context.Context, categories ...owvd.CapabilityCategory) (*owvs.Capabilities, error) {
	return owvd.GetCapabilities(ctx, d.URL, d.Security)
}

// GetDeviceInformation 查询信息
func (d *Device[Data]) GetDeviceInformation(ctx context.Context) (*owvd.Information, error) {
	return owvd.GetDeviceInformation(ctx, d.URL, d.Security)
}

// GetProfiles 查询媒体属性
func (d *Device[Data]) GetProfiles(ctx context.Context) ([]*owvs.Profile, error) {
	if d.Capabilities == nil || d.Capabilities.Media == nil || d.Capabilities.Media.XAddr == "" {
		return nil, ErrMediaCapabilityUnsupported
	}
	return owvm.GetProfiles(ctx, d.Capabilities.Media.XAddr, d.Security)
}
