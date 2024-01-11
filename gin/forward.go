package gin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Forward 转发请求和响应，baseURL 是 http(s)://host:port
func Forward(ctx *gin.Context, client *http.Client, baseURL string) error {
	// 请求
	req, err := http.NewRequestWithContext(ctx, ctx.Request.Method, baseURL, ctx.Request.Body)
	if err != nil {
		return err
	}
	// path + query
	req.URL.Path = ctx.Request.URL.Path
	req.URL.RawQuery = ctx.Request.URL.RawQuery
	// header
	for k, v := range ctx.Request.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	//
	return ForwardResponseWithRequest(ctx, client, req)
}

// ForwardWithURL 转发请求和响应，
func ForwardWithURL(ctx *gin.Context, client *http.Client, url string) error {
	return ForwardResponse(ctx, client, ctx.Request.Method, url, ctx.Request.Header, ctx.Request.Body)
}

// ForwardResponse 使用参数构造新的请求后转发响应
func ForwardResponse(ctx *gin.Context, client *http.Client, method, url string, header http.Header, body io.Reader) error {
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
	return ForwardResponseWithRequest(ctx, client, req)
}

// ForwardResponseWithRequest 使用新的请求转发响应
func ForwardResponseWithRequest(ctx *gin.Context, client *http.Client, req *http.Request) error {
	// 发送
	_res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer _res.Body.Close()
	// status code
	ctx.Writer.WriteHeader(_res.StatusCode)
	// header
	header := ctx.Writer.Header()
	for k, v := range _res.Header {
		for _, vv := range v {
			header.Add(k, vv)
		}
	}
	// body
	_, err = io.Copy(ctx.Writer, _res.Body)
	return err
}

// ForwardResponseWithJSONBody 使用参数构造新的 json body 请求后转发响应
func ForwardResponseWithJSONBody(ctx *gin.Context, client *http.Client, method, url string, header http.Header, body any) error {
	// body
	var data *bytes.Buffer = nil
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
	return ForwardResponse(ctx, client, method, url, header, data)
}
