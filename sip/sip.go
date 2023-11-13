package sip

import (
	"errors"
	"fmt"
	"io"
)

// 方法
const (
	MethodRegister  string = "REGISTER"
	MethodInvite    string = "INVITE"
	MethodACK       string = "ACK"
	MethodBye       string = "BYE"
	MethodMessage   string = "MESSAGE"
	MethodNotify    string = "NOTIFY"
	MethodSubscribe string = "SUBSCRIBE"
	MethodInfo      string = "INFO"
)

// 一些常量
const (
	BranchPrefix  = "z9hG4bK"
	SIPVersion    = "SIP/2.0"
	MaxMessageLen = 1024 * 10
	TCP           = "SIP/2.0/TCP"
	TCPS          = "SIPS/2.0/TCP"
	UDP           = "SIP/2.0/UDP"
	UDPS          = "SIPS/2.0/UDP"
	CRLF          = "\r\n"
	SIP           = "sip"
	SIPS          = "sips"
)

// 一些错误
var (
	errLargeMessage            = errors.New("large message")
	errMissHeaderVia           = errors.New("miss header via")
	errMissHeaderFrom          = errors.New("miss header from")
	errMissHeaderTo            = errors.New("miss header to")
	errMissHeaderCSeq          = errors.New("miss header cseq")
	errMissHeaderCallID        = errors.New("miss header call-id")
	errMissHeaderContentLength = errors.New("miss header content-length")
	errTransactionExists       = errors.New("transaction exists")
	errUnknownAddress          = errors.New("unknown address")
	errServerClosed            = errors.New("server closed")
)

// 状态码列表，不全，需要的时候再添加
const (
	// 1xx
	StatusTrying               = "100"
	StatusDialogEstablishment  = "101"
	StatusRinging              = "180"
	StatusCallIsBeingForwarded = "181"
	StatusQueued               = "182"
	StatusSessionProgress      = "183"
	// 2xx
	StatusOK = "200"
	// 3xx
	StatusMultipleChoices    = "300"
	StatusMovedPermently     = "301"
	StatusMovedTemporarily   = "302"
	StatusUseProxy           = "305"
	StatusAlternativeService = "380"
	// 4xx
	StatusBadRequest                    = "400"
	StatusUnauthorized                  = "401"
	StatusPaymentRequired               = "402"
	StatusForbidden                     = "403"
	StatusNotFound                      = "404"
	StatusMethodNotAllowed              = "405"
	StatusNotAcceptable                 = "406"
	StatusProxyAuthenticationRequired   = "407"
	StatusRequestTimeout                = "408"
	StatusConflict                      = "409"
	StatusGone                          = "410"
	StatusLengthRequired                = "411"
	StatusRequestEntityTooLarge         = "413"
	StatusRequestURITooLong             = "414"
	StatusUnsupportedMediaType          = "415"
	StatusUnsupportedURIScheme          = "416"
	StatusBadExtension                  = "420"
	StatusExtensionRequired             = "421"
	StatusIntervalTooBrief              = "423"
	StatusTemporarilyUnavalilable       = "480"
	StatusCallOrTransactionDoesNotExist = "481"
	StatusLoopDetected                  = "482"
	StatusTooManyHops                   = "483"
	StatusAddressIncomplete             = "484"
	StatusAmbiguous                     = "485"
	StatusBusyHere                      = "486"
	StatusRequetTerminated              = "487"
	StatusNotAcceptableHere             = "488"
	StatusRequestPending                = "489"
	StatusUndecipherable                = "490"
	// 5xx
	StatusServerInternalError = "500"
	StatusNotImplemented      = "501"
	StatusBadGateway          = "502"
	StatusServiceUnavailable  = "503"
	StatusServerTimeout       = "504"
	StatusVersionNotSupported = "505"
	StatusMessageTooLarge     = "513"
	// 6xx
	StatusBusyEverywhere       = "600"
	StatusDecline              = "603"
	StatusDoesNotExistAnywhere = "604"
)

// StatusPhrase 返回 s 的短语
func StatusPhrase(s string) string {
	switch s {
	case StatusTrying:
		return "Trying"
	case StatusDialogEstablishment:
		return "Dialog Establishment"
	case StatusRinging:
		return "Ringing"
	case StatusCallIsBeingForwarded:
		return "Call Is Being Forwarded"
	case StatusQueued:
		return "Queued"
	case StatusSessionProgress:
		return "Session Progress"
	case StatusOK:
		return "OK"
	case StatusMultipleChoices:
		return "Multiple Choices"
	case StatusMovedPermently:
		return "Moved Permently"
	case StatusMovedTemporarily:
		return "Moved Temporarily"
	case StatusUseProxy:
		return "Use Proxy"
	case StatusAlternativeService:
		return "Alternative Service"
	case StatusBadRequest:
		return "Bad Request"
	case StatusUnauthorized:
		return "Unauthorized"
	case StatusPaymentRequired:
		return "Payment Required"
	case StatusForbidden:
		return "Forbidden"
	case StatusNotFound:
		return "Not Found"
	case StatusMethodNotAllowed:
		return "Method Not Allowed"
	case StatusNotAcceptable:
		return "Not Acceptable"
	case StatusProxyAuthenticationRequired:
		return "Proxy Authentication Required"
	case StatusRequestTimeout:
		return "Request Timeout"
	case StatusConflict:
		return "Conflict"
	case StatusGone:
		return "Gone"
	case StatusLengthRequired:
		return "Length Required"
	case StatusRequestEntityTooLarge:
		return "Request Entity Too Large"
	case StatusRequestURITooLong:
		return "Request-URI Too Long"
	case StatusUnsupportedMediaType:
		return "Unsupported Media Type"
	case StatusUnsupportedURIScheme:
		return "Unsupported URI Scheme"
	case StatusBadExtension:
		return "Bad Extension"
	case StatusExtensionRequired:
		return "Extension Required"
	case StatusIntervalTooBrief:
		return "Interval Too Brief"
	case StatusTemporarilyUnavalilable:
		return "Temporarily Unavalilable"
	case StatusCallOrTransactionDoesNotExist:
		return "Call/Transaction Does Not Exist"
	case StatusLoopDetected:
		return "Loop Detected"
	case StatusTooManyHops:
		return "Too Many Hops"
	case StatusAddressIncomplete:
		return "Address Incomplete"
	case StatusAmbiguous:
		return "Ambiguous"
	case StatusBusyHere:
		return "Busy Here"
	case StatusRequetTerminated:
		return "Requet Terminated"
	case StatusNotAcceptableHere:
		return "Not Acceptable Here"
	case StatusRequestPending:
		return "Request Pending"
	case StatusUndecipherable:
		return "Undecipherable"
	case StatusServerInternalError:
		return "Server Internal Error"
	case StatusNotImplemented:
		return "Not Implemented"
	case StatusBadGateway:
		return "Bad Gateway"
	case StatusServiceUnavailable:
		return "Service Unavailable"
	case StatusServerTimeout:
		return "Server Timeout"
	case StatusVersionNotSupported:
		return "Version Not Supported"
	case StatusMessageTooLarge:
		return "Message Too Large"
	case StatusBusyEverywhere:
		return "Busy Everywhere"
	case StatusDecline:
		return "Decline"
	case StatusDoesNotExistAnywhere:
		return "Does Not Exist Anywhere"
	}
	return "Unknown Status Code"
}

// StatusCode 返回 s 的整数值
func StatusCode(s string) int {
	switch s {
	case StatusTrying:
		return 100
	case StatusDialogEstablishment:
		return 101
	case StatusRinging:
		return 180
	case StatusCallIsBeingForwarded:
		return 181
	case StatusQueued:
		return 182
	case StatusSessionProgress:
		return 183
	case StatusOK:
		return 200
	case StatusMultipleChoices:
		return 300
	case StatusMovedPermently:
		return 301
	case StatusMovedTemporarily:
		return 302
	case StatusUseProxy:
		return 305
	case StatusAlternativeService:
		return 380
	case StatusBadRequest:
		return 400
	case StatusUnauthorized:
		return 401
	case StatusPaymentRequired:
		return 402
	case StatusForbidden:
		return 403
	case StatusNotFound:
		return 404
	case StatusMethodNotAllowed:
		return 405
	case StatusNotAcceptable:
		return 406
	case StatusProxyAuthenticationRequired:
		return 407
	case StatusRequestTimeout:
		return 408
	case StatusConflict:
		return 409
	case StatusGone:
		return 410
	case StatusLengthRequired:
		return 411
	case StatusRequestEntityTooLarge:
		return 413
	case StatusRequestURITooLong:
		return 414
	case StatusUnsupportedMediaType:
		return 415
	case StatusUnsupportedURIScheme:
		return 416
	case StatusBadExtension:
		return 420
	case StatusExtensionRequired:
		return 421
	case StatusIntervalTooBrief:
		return 423
	case StatusTemporarilyUnavalilable:
		return 480
	case StatusCallOrTransactionDoesNotExist:
		return 481
	case StatusLoopDetected:
		return 482
	case StatusTooManyHops:
		return 483
	case StatusAddressIncomplete:
		return 484
	case StatusAmbiguous:
		return 485
	case StatusBusyHere:
		return 486
	case StatusRequetTerminated:
		return 487
	case StatusNotAcceptableHere:
		return 488
	case StatusRequestPending:
		return 489
	case StatusUndecipherable:
		return 490
	case StatusServerInternalError:
		return 500
	case StatusNotImplemented:
		return 501
	case StatusBadGateway:
		return 502
	case StatusServiceUnavailable:
		return 503
	case StatusServerTimeout:
		return 504
	case StatusVersionNotSupported:
		return 505
	case StatusMessageTooLarge:
		return 513
	case StatusBusyEverywhere:
		return 600
	case StatusDecline:
		return 603
	case StatusDoesNotExistAnywhere:
		return 604
	}
	return 0
}

// Writer 主要用来格式化 message
type Writer interface {
	WriteByte(byte) error
	WriteString(string) (int, error)
	io.Writer
}

// Reader 主要用于解析 message
type Reader interface {
	ReadLine() (string, error)
	io.Reader
}

// ResponseError 表示 sip 响应消息的错误
type ResponseError struct {
	Status string
	Phrase string
	Err    string
}

func (e ResponseError) Error() string {
	if e.Err != "" {
		return e.Err
	}
	e.Err = fmt.Sprintf("stauts %s phrase %s", e.Status, e.Phrase)
	return e.Err
}

// NewResponseError 是 ResponseError 构造函数
func NewResponseError(status, phrase, err string) *ResponseError {
	return &ResponseError{
		Status: status,
		Phrase: phrase,
		Err:    err,
	}
}
