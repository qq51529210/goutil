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
	msg, addr, err := request.New(m.Cascade, "", sip.MethodRegister, "")
	if err != nil {
		return err
	}
	if m.Authorization != "" {
		msg.Header.Set(StrAuthorization, m.Authorization)
	}
	// 这两个是一样的
	msg.Header.To.URI = msg.Header.From.URI
	msg.Header.Expires = m.Expires
	//
	return m.Ser.RequestWithContext(ctx, msg, addr, m)
}
