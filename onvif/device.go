package onvif

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	owvd "goutil/onvif/wsdl/ver10/device"
	owvm "goutil/onvif/wsdl/ver10/media"
	owvs "goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// 错误
var (
	ErrMediaCapabilityUnsupported = errors.New("media capability unsupported")
)

// Device 表示设备
type Device struct {
	URL string
	// 需要主动调用赋值
	*owvs.Capabilities
	// 设备地址
	host string
	// 用户名
	username string
	// 密码
	password string
}

// NewDevice 初始化字段
func NewDevice(host, username, password string) *Device {
	// 初始化
	d := new(Device)
	d.URL = fmt.Sprintf("http://%s/onvif/device_service", host)
	d.username = username
	d.password = password
	d.host = host
	// 返回
	return d
}

// NewDeviceWithCapabilities 初始化字段，同时获取能力
func NewDeviceWithCapabilities(ctx context.Context, host, username, password string) (*Device, error) {
	d := NewDevice(host, username, password)
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
func (d *Device) ReplaceXAddrHost(xaddr string) string {
	u, _ := url.Parse(xaddr)
	u.Host = d.host
	return u.String()
}

func (d *Device) replaceIP(host string) string {
	n := strings.LastIndex(host, ":")
	// host 是一个 ip
	if n < 1 {
		return host
	}
	// 提取 ip
	ip := d.host
	if i := strings.LastIndex(ip, ":"); i > 0 {
		ip = ip[:i]
	}
	return ip + host[n:]
}

// GetSystemDateAndTime 查询日期时间
func (d *Device) GetSystemDateAndTime(ctx context.Context) (*owvs.SystemDateTime, error) {
	return owvd.GetSystemDateAndTime(ctx, d.URL)
}

// GetCapabilities 查询能力
func (d *Device) GetCapabilities(ctx context.Context, categories ...owvd.CapabilityCategory) (*owvs.Capabilities, error) {
	return owvd.GetCapabilities(ctx, d.URL, soap.NewSecurity(d.username, d.password))
}

// GetDeviceInformation 查询信息
func (d *Device) GetDeviceInformation(ctx context.Context) (*owvd.Information, error) {
	return owvd.GetDeviceInformation(ctx, d.URL, soap.NewSecurity(d.username, d.password))
}

// GetProfiles 查询媒体属性
func (d *Device) GetProfiles(ctx context.Context) ([]*owvs.Profile, error) {
	if !d.IsMediaServiceOK() {
		return nil, ErrMediaCapabilityUnsupported
	}
	return owvm.GetProfiles(ctx, d.Capabilities.Media.XAddr, soap.NewSecurity(d.username, d.password))
}
