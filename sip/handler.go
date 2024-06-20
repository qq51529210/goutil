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

// reqFuncChain 请求调用链
type reqFuncChain struct {
	// 保存调用链函数
	f []HandleRequestFunc
	// 当前调用的函数下标
	i int
}

// Next 执行调用链中剩下的所有函数
func (c *reqFuncChain) Next(r *Request) {
	for c.i < len(c.f) {
		c.f[c.i](r)
		c.i++
	}
}

// Abort 结束调用链
func (c *reqFuncChain) Abort() {
	c.i = len(c.f)
}

// resFuncChain 响应调用链
type resFuncChain struct {
	// 保存调用链函数
	f []HandleResponseFunc
	// 当前调用的函数下标
	i int
}

// Next 执行调用链中剩下的所有函数
func (c *resFuncChain) Next(r *Response) {
	for c.i < len(c.f) {
		c.f[c.i](r)
		c.i++
	}
}

// Abort 结束调用链
func (c *resFuncChain) Abort() {
	c.i = len(c.f)
}
