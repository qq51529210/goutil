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
	msg, addr := request.NewRegister(m.Cascade, m.Expires)
	if m.Authorization != "" {
		msg.Header.Set(StrAuthorization, m.Authorization)
	}
	//
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
