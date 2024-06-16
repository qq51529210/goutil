package sip

import (
	"strings"
)

// HandleRequestFunc 请求回调函数
type HandleRequestFunc func(*Request)

// HandleResponseFunc 响应回调函数
type HandleResponseFunc func(*Response)

// HandleRequestFunc 请求方法未注册时回调函数，res 是回复的消息
type HandleRequestNotFoundFunc func(req *Message) (res *Message)

type handleFunc struct {
	// 请求消息回调
	reqFunc map[string][]HandleRequestFunc
	// 响应消息回调
	resFunc map[string][]HandleResponseFunc
	// 请求消息方法未注册时回调
	reqNotFoundFunc HandleRequestNotFoundFunc
	// 已注册的请求消息方法，用于 len(reqNotFoundFunc)==0 时自动回复
	reqMethods string
}

// RequestFunc 注册请求消息回调函数链，并发不安全，要提前设置好
func (h *handleFunc) RequestFunc(method string, funcs ...HandleRequestFunc) {
	if len(funcs) < 1 {
		panic("invalid request callback func")
	}
	if h.reqFunc == nil {
		h.reqFunc = make(map[string][]HandleRequestFunc)
	}
	f, ok := h.reqFunc[method]
	if !ok {
		f = make([]HandleRequestFunc, 0)
	}
	h.reqFunc[method] = append(f, funcs...)
	//
	methods := make([]string, 0)
	for k := range h.reqFunc {
		methods = append(methods, k)
	}
	h.reqMethods = strings.Join(methods, ",")
}

// ResponseFunc 注册响应消息回调函数链，并发不安全，要提前设置好
func (h *handleFunc) ResponseFunc(method string, funcs ...HandleResponseFunc) {
	if len(funcs) < 1 {
		panic("invalid response callback func")
	}
	if h.resFunc == nil {
		h.resFunc = make(map[string][]HandleResponseFunc)
	}
	f, ok := h.resFunc[method]
	if !ok {
		f = make([]HandleResponseFunc, 0)
	}
	h.resFunc[method] = append(f, funcs...)
}

// RequestNotFoundFunc 注册请求消息回调函数链，检测到没有注册的方法时调用，并发不安全，要提前设置好
func (h *handleFunc) RequestNotFoundFunc(fun HandleRequestNotFoundFunc) {
	h.reqNotFoundFunc = fun
}
