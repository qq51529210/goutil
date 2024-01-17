package onvif

import (
	"context"
	"fmt"
	"goutil/soap"
)

// Device 表示设备
type Device struct {
	URL string
	*soap.Security
}

// NewDevice 初始化字段
func NewDevice(ctx context.Context, host, username, password string) *Device {
	// 初始化
	d := new(Device)
	d.URL = fmt.Sprintf("http://%s/onvif/device_service", host)
	d.Security = new(soap.Security)
	d.Security.Init(username, password)
	// 返回
	return d
}
