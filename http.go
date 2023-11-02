package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
)

// HTTPStatusError 表示状态错误
type HTTPStatusError int

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("status code %d", e)
}

var (
	// HTTPQueryTag 是 HTTPQuery 解析 tag 的名称
	HTTPQueryTag = "query"
)

// HTTPQuery 将结构体 v 格式化到 url.Values
// 只扫描一层，并略过空值
func HTTPQuery(v any, q url.Values) url.Values {
	if q == nil {
		q = make(url.Values)
	}
	rv := reflect.ValueOf(v)
	vk := rv.Kind()
	if vk == reflect.Pointer {
		rv = rv.Elem()
		vk = rv.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	return httpQuery(rv, q)
}

func httpQuery(v reflect.Value, q url.Values) url.Values {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构，只一层
		if fvk == reflect.Struct {
			continue
		}
		if fvk == reflect.String {
			// 空字符串
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(HTTPQueryTag)
		if tn == "" || tn == "-" {
			continue
		}
		q.Add(tn, fmt.Sprintf("%v", fv.Interface()))
	}
	return q
}

// HTTPClient 封装 http 请求代码
type HTTPClient[reqData any] struct {
	// 客户端，需要初始化
	C *http.Client
	// 地址，需要初始化
	URL string
	// 结果回调，需要初始化
	OnRes func(*http.Response) error
}

// Do 发送请求
func (c *HTTPClient[reqData]) Do(ctx context.Context, method string, query url.Values, reqBody *reqData) error {
	return HTTPWithContext[reqData](ctx, c.C, method, c.URL, query, reqBody, c.OnRes)
}

// HTTPWithContext 封装 http 操作
// ctx 超时上下文
// client 使用的客户端
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// onResponse 处理结果
func HTTPWithContext[reqData any](ctx context.Context, client *http.Client, method, url string, query url.Values, reqBody *reqData, onResponse func(res *http.Response) error) error {
	// body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		q := req.URL.Query()
		for k, vs := range query {
			for _, v := range vs {
				q.Add(k, v)
			}
		}
		req.URL.RawQuery = q.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 回调
	return onResponse(res)
}

// HTTPServer 封装代码
type HTTPServer struct {
	S http.Server
	// 证书路径
	CertFile string
	KeyFile  string
}

// Serve 如果证书路径不为空，监听 tls
func (s *HTTPServer) Serve() error {
	if s.CertFile != "" && s.KeyFile != "" {
		return s.S.ListenAndServeTLS(s.CertFile, s.KeyFile)
	}
	return s.S.ListenAndServe()
}
