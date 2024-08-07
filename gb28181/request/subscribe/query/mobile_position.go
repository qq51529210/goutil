package subscribe

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// MobilePosition 是 SendMobilePosition 的参数
type MobilePosition struct {
	Ser    *sip.Server
	Device request.Request
	// 过期时间
	Expire int64
	// 上报的时间间隔，单位秒
	Interval int64
	// 结果
	result string
}

func (m *MobilePosition) SetResult(s string) {
	m.result = s
}

// SendMobilePosition 发送 Request-Subscribe-MobilePosition 请求
func SendMobilePosition(ctx context.Context, m *MobilePosition) (string, error) {
	// 消息
	var body xml.Subscribe
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdMobilePosition
	body.DeviceID = m.Device.GetToID()
	body.SN = sip.GetSNString()
	body.Interval = m.Interval
	//
	var result request.XMLResult
	return result.Result, request.SendSubscribe(ctx, m.Ser, m.Device, &body, m.Expire, &result)
}
