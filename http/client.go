package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// JSONClient 封装 http 请求代码
type JSONClient[Data any] struct {
	// 客户端，需要初始化
	C *http.Client
	// 地址，需要初始化
	URL string
	// 结果回调，需要初始化
	OnRes func(*http.Response) error
}

// Do 发送请求
func (c *JSONClient[Data]) Do(ctx context.Context, method string, query url.Values, body *Data) error {
	return Request[Data](ctx, c.C, method, c.URL, query, body, c.OnRes)
}

// Request 封装 http 操作
// ctx 超时上下文
// client 使用的客户端
// method 方法
// url 请求地址
// query 请求参数
// body 格式化 json 后写入 body
// onResponse 处理结果
func Request[Data any](ctx context.Context, client *http.Client, method, url string, query url.Values, data *Data, onResponse func(res *http.Response) error) error {
	// body
	var body io.Reader = nil
	if data != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(data)
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
