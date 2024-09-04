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
}

// SendRegister 注册
func SendRegister(ctx context.Context, m *Register) error {
	msg := request.NewRegister(m.Cascade, m.Expires)
	// 网络地址
	addr, err := m.Cascade.GetNetAddr()
	if err != nil {
		return err
	}
	if m.Authorization != "" {
		msg.Header.Set(StrAuthorization, m.Authorization)
	}
	//
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
