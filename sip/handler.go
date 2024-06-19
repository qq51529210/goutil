package sip

// HandleRequestFunc 请求回调函数
type HandleRequestFunc func(*Request)

// HandleResponseFunc 响应回调函数
type HandleResponseFunc func(*Response)

type handleFunc struct {
	// 请求消息回调
	reqFunc map[string][]HandleRequestFunc
	// 响应消息回调
	resFunc map[string][]HandleResponseFunc
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
