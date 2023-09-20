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
	"time"
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

// DoWithContext 发送请求，query 使用 HTTPQuery 来组成参数
func (c *HTTPClient[reqData]) Do(method string, query any, reqBody *reqData, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return c.DoWithContext(ctx, method, query, reqBody)
}

// DoWithContext 发送请求，query 使用 HTTPQuery 来组成参数
func (c *HTTPClient[reqData]) DoWithContext(ctx context.Context, method string, query any, reqBody *reqData) error {
	// 请求 body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequestWithContext(ctx, method, c.URL, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = HTTPQuery(query, req.URL.Query()).Encode()
	}
	// 发送
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 回调
	return c.OnRes(res)
}
