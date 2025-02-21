package register

import (
	"context"
	"goutil/gb28181/request"
	"goutil/sip"
)

// Register 是 SendRegister 的参数
type Register struct {
	Ser     *sip.Server
	Cascade request.Request
	// 过期时间
	Expires string
	// Header.Authorization
	Authorization string
	// 追踪标识
	TraceID string
}

// SendRegister 注册
func SendRegister(ctx context.Context, m *Register) error {
	// 网络地址
	addr, err := m.Cascade.GetNetAddr()
	if err != nil {
		return err
	}
	// 消息
	msg := request.NewRegister(m.Cascade, addr.Network(), m.Expires)
	if m.Authorization != "" {
		msg.Header.Set(StrAuthorization, m.Authorization)
	}
	//
	return m.Ser.RequestWithContext(ctx, m.TraceID, msg, addr, m)
}
