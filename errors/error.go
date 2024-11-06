package errors

import (
	"runtime"
)

type Error struct {
	// 发生错误的文件路径
	Path string
	// 发生错误的行
	Line int
	// 错误代码
	Code int
	// 错误简要
	Msg string
	// 错误详细
	Err string
}

// Error 实现 errors 接口
func (e *Error) Error() string {
	return e.Err
}

// New 返回具有堆栈信息的 Error
func New(code int, msg string, err error) *Error {
	_, path, line, ok := runtime.Caller(1)
	if !ok {
		path = "???"
		line = -1
	}
	return &Error{
		Path: path,
		Line: line,
		Code: code,
		Msg:  msg,
		Err:  err.Error(),
	}
}
