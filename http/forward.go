package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

func forwardResponse(client *http.Client, req *http.Request, res http.ResponseWriter) error {
	// 发送
	_res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer _res.Body.Close()
	// status code
	res.WriteHeader(_res.StatusCode)
	// header
	header := res.Header()
	for k, v := range _res.Header {
		for _, vv := range v {
			header.Add(k, vv)
		}
	}
	// body
	_, err = io.Copy(res, _res.Body)
	return err
}

// Forward 转发请求和响应
func Forward(ctx context.Context, client *http.Client, baseURL string, req *http.Request, res http.ResponseWriter) error {
	// 请求
	_req, err := http.NewRequestWithContext(ctx, req.Method, baseURL, req.Body)
	if err != nil {
		return err
	}
	// path + query
	_req.URL.Path = req.URL.Path
	_req.URL.RawQuery = req.URL.RawQuery
	// header
	for k, v := range req.Header {
		for _, vv := range v {
			_req.Header.Add(k, vv)
		}
	}
	// 转发
	return forwardResponse(client, _req, res)
}

// ForwardResponse 使用参数构造新的请求后转发响应
func ForwardResponse(ctx context.Context, client *http.Client, method, url string, header http.Header, body io.Reader, res http.ResponseWriter) error {
	// 请求
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return err
	}
	// header
	for k, v := range header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	// 转发
	return forwardResponse(client, req, res)
}

// ForwardResponseWithJSONBody 使用参数构造新的 json body 请求后转发响应
func ForwardResponseWithJSONBody(ctx context.Context, client *http.Client, method, url string, header http.Header, body any, res http.ResponseWriter) error {
	// body
	var data io.ReadWriter
	if body != nil {
		data = bytes.NewBuffer(nil)
		json.NewEncoder(data).Encode(body)
	}
	if header == nil {
		header = make(http.Header)
	}
	if header.Get("Content-Type") == "" {
		header.Add("Content-Type", "application/json; charset=utf-8")
	}
	return ForwardResponse(ctx, client, method, url, header, data, res)
}
