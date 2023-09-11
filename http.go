package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// HTTPStatusError 表示状态错误
type HTTPStatusError int

func (e HTTPStatusError) Error() string {
	return fmt.Sprintf("status code %d", e)
}

// HTTP 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
// timeout 超时
func HTTP[reqData, resData any](method, url string, query url.Values, reqBody *reqData, resBody *resData, statusCode int, timeout time.Duration) error {
	// 请求
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return nil
}

// HTTPTo 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 写入响应的 body 数据
// statusCode 用于判断状态码
// timeout 超时
func HTTPTo[reqData any](method, url string, query url.Values, reqBody *reqData, resBody io.Writer, statusCode int, timeout time.Duration) error {
	// 请求
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		_, err = io.Copy(resBody, res.Body)
	}
	return err
}

// HTTPFrom 封装 http 操作
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
// timeout 超时
func HTTPFrom[resData any](method, url string, query url.Values, reqBody io.Reader, resBody *resData, statusCode int, timeout time.Duration) error {
	// 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	// 参数
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return err
}

// HTTPWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
func HTTPWithContext[reqData, resData any](ctx context.Context, method, url string, query url.Values, reqBody *reqData, resBody *resData, statusCode int) error {
	// body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return nil
}

// HTTPToWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 格式化 json 后写入 body
// resBody 写入响应的 body 数据
// statusCode 用于判断状态码
func HTTPToWithContext[reqData any](ctx context.Context, method, url string, query url.Values, reqBody *reqData, resBody io.Writer, statusCode int) error {
	// body
	var body io.Reader = nil
	if reqBody != nil {
		buf := bytes.NewBuffer(nil)
		json.NewEncoder(buf).Encode(reqBody)
		body = buf
	}
	// 请求
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		_, err = io.Copy(resBody, res.Body)
	}
	return err
}

// HTTPFromWithContext 封装 http 操作
// ctx 超时上下文
// method 方法
// url 请求地址
// query 请求参数
// reqBody 用于读取发送 body
// resBody 用于解析响应 body 中的 json
// statusCode 用于判断状态码
func HTTPFromWithContext[resData any](ctx context.Context, method, url string, query url.Values, reqBody io.Reader, resBody *resData, statusCode int) error {
	// 请求
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return err
	}
	// query
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}
	// 发送
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 状态码
	if res.StatusCode != statusCode {
		return HTTPStatusError(res.StatusCode)
	}
	// 解析
	if resBody != nil {
		return json.NewDecoder(res.Body).Decode(resBody)
	}
	return err
}
