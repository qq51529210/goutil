package zlm

import (
	"fmt"
)

// var (
// 	// ErrServerNotFound 服务不存在
// 	ErrServerNotFound = &Error{Err: errors.New("server not found"), Msg: "流媒体服务不存在"}
// 	// ErrServerNotAvailable 服务不可用
// 	ErrServerNotAvailable = &Error{Err: errors.New("server not available"), Msg: "流媒体服务暂时不可用"}
// 	// ErrServerNotAssign 服务未分配
// 	ErrServerNotAssign = &Error{Err: errors.New("server not assign"), Msg: "流媒体服务未分配"}
// 	// ErrUnknownStream 未知的流
// 	ErrUnknownStream = &Error{Err: errors.New("unknown stream"), Msg: "未知的媒体流"}
// )

// // 中文
// const (
// 	MediaInfoNotFound = "媒体流不存在"
// 	SanpshotError     = "读取快照异常"
// )

// // Error 用于友好提示前端
// type Error struct {
// 	Err error
// 	// 提示
// 	Msg string
// }

// func (e *Error) Error() string {
// 	return e.Err.Error()
// }

// apiError 用于接收接口错误
type apiError struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	SerID string `json:"-"`
	Path  string `json:"-"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("zlm %s api call %s code %d msg %s", e.SerID, e.Path, e.Code, e.Msg)
}
