package zlm

import (
	"errors"
	"fmt"
	gh "goutil/http"
)

const (
	// 正确码
	codeTrue = 0
)

// 定义一些错误以便全局使用，看名称猜意思
var (
	ErrServerNotAvailable = errors.New("server not available")
	ErrMediaNotFound      = errors.New("media not found")
	ErrToken              = errors.New("error token")
)

// apiError 用于接收接口错误
type apiError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Path string `json:"-"`
}

func (e *apiError) Error() string {
	return fmt.Sprintf("call %s code %d msg %s", e.Path, e.Code, e.Msg)
}

// IsZLMError 检查是否 zlm 定义的错误
func IsZLMError(err error) bool {
	if err == ErrServerNotAvailable {
		return true
	}
	if _, ok := err.(*apiError); ok {
		return true
	}
	if _, ok := err.(gh.StatusError); ok {
		return true
	}
	return false
}
