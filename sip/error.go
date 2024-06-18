package sip

import (
	"errors"
	"fmt"
)

// 一些错误
var (
	// 这个不是错误，
	// 因为 context.Context 在 Done() 返回后，
	// 要求 Err() 必须是错误
	ErrFinish = errors.New("finish")
	//
	ErrLargeMessage     = errors.New("large message")
	ErrUnknownAddress   = errors.New("only tcp/udp supported")
	ErrServerShutdown   = errors.New("server shutdown")
	ErrTransactionExist = errors.New("transaction exists")
	//
	errMissHeaderVia           = errors.New("miss header via")
	errMissHeaderFrom          = errors.New("miss header from")
	errMissHeaderTo            = errors.New("miss header to")
	errMissHeaderCSeq          = errors.New("miss header cseq")
	errMissHeaderCallID        = errors.New("miss header call-id")
	errMissHeaderContentLength = errors.New("miss header content-length")
	//
)

// ResponseError 表示 sip 响应消息的错误
type ResponseError struct {
	Status string
	Phrase string
	Err    string
}

func (e *ResponseError) Error() string {
	if e.Err != "" {
		return e.Err
	}
	e.Err = fmt.Sprintf("stauts %s phrase %s", e.Status, e.Phrase)
	return e.Err
}
