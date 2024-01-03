package http

import "fmt"

// Result 表示返回的结果
type Result[T any] struct {
	// 状态码
	Code int `json:"code,omitempty"`
	// 错误短语
	Msg string `json:"msg,omitempty"`
	// 没有错误时候的数据
	Data T `json:"data,omitempty"`
}

// Error 实现 error 接口
func (c *Result[T]) Error() string {
	return fmt.Sprintf("code: %d, msg: %s", c.Code, c.Msg)
}

// StatusError 表示状态错误
type StatusError int

func (e StatusError) Error() string {
	return fmt.Sprintf("status code %d", e)
}
