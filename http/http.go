package http

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Result 表示返回的结果
type Result[T any] struct {
	// 追踪标识，用于日志快速定位
	Trace string `json:"trace,omitempty"`
	// 服务类型，用于日志快速定位发生错误的服务
	Service string `json:"service,omitempty"`
	// 错误码
	Code int `json:"code"`
	// 正确时返回的数据
	Data T `json:"data,omitempty"`
	// 错误短语
	Msg string `json:"msg,omitempty"`
	// 错误详细
	Err string `json:"err,omitempty"`
	// 额外的数据，用户自定义
	Ext string `json:"ext,omitempty"`
}

// Error 实现 error 接口
func (c *Result[T]) Error() string {
	if c.Err != "" {
		return c.Err
	}
	if c.Msg != "" {
		return c.Msg
	}
	return fmt.Sprintf("code %d", c.Code)
}

// StatusError 表示状态错误
type StatusError int

func (e StatusError) Error() string {
	return fmt.Sprintf("status code %d", e)
}

func Json(data any) string {
	var str strings.Builder
	enc := json.NewEncoder(&str)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(data)
	s := str.String()
	if s != "" {
		n := len(s) - 1
		if s[n] == '\n' {
			s = s[:n]
		}
	}
	return s
}
