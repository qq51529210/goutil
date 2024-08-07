package subscribe

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// Catalog 是 SendCatalog 的参数
type Catalog struct {
	Ser    *sip.Server
	Device request.Request
	// 过期时间
	Expire int64
	// 起始时间
	StartTime string `json:"startTime" binding:"omitempty,gb_time"`
	// 终止时间
	EndTime string `json:"endTime" binding:"omitempty,gb_time"`
	// 结果
	result string
}

func (m *Catalog) SetResult(s string) {
	m.result = s
}

// SendCatalog 发送 Request-Subscribe-Catalog 请求
func SendCatalog(ctx context.Context, m *Catalog) (string, error) {
	// 消息
	var body xml.Subscribe
	body.XMLName.Local = xml.TypeQuery
	body.CmdType = xml.CmdCatalog
	body.DeviceID = m.Device.GetToID()
	body.SN = sip.GetSNString()
	body.StartTime = m.StartTime
	body.EndTime = m.EndTime
	//
	var result request.XMLResult
	return result.Result, request.SendSubscribe(ctx, m.Ser, m.Device, &body, m.Expire, &result)
}
