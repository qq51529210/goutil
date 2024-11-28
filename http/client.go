package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

const (
	// ContentTypeJSON header Content-Type
	ContentTypeJSON = "application/json; charset=utf-8"
)

// JSONClient 封装 http 请求代码
type JSONClient struct {
	// 客户端，需要初始化
	C *http.Client
	// 地址，需要初始化
	URL string
	// 结果回调，需要初始化
	OnRes func(*http.Response) error
}

// Do 发送请求
func (c *JSONClient) Do(ctx context.Context, method string, query url.Values, body any) error {
	return JSONRequest(ctx, c.C, method, c.URL, query, body, c.OnRes)
}

// JSONRequest 封装 http json 请求
func JSONRequest(ctx context.Context, client *http.Client, method, url string, query url.Values, data any, onResponse func(res *http.Response) error) error {
	// body
	var body io.Reader = nil
	if data != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(data)
		body = buf
	}
	return Request(ctx, client, method, url, query, map[string]string{"Content-Type": ContentTypeJSON}, body, onResponse)
}

// Request 封装 http 请求
func Request(ctx context.Context, client *http.Client, method, url string, query url.Values, header map[string]string, body io.Reader, onResponse func(res *http.Response) error) error {
	// 请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	for k, v := range header {
		req.Header.Set(k, v)
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
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 回调
	return onResponse(res)
}

// HandleResponse 封装处理代码，判断 status 是否 200 ，然后解析 body json 到 data
func HandleResponse[M any](res *http.Response, data *Result[M]) error {
	// 状态码
	if res.StatusCode != http.StatusOK {
		return StatusError(res.StatusCode)
	}
	// 解析数据
	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		return err
	}
	if data.Code != 0 {
		return data
	}
	return nil
}

// HandleResponseStatus 封装处理代码，只判断 status 是否 200
func HandleResponseStatus(res *http.Response) error {
	// 状态码
	if res.StatusCode != http.StatusOK {
		return StatusError(res.StatusCode)
	}
	return nil
}
