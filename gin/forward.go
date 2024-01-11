package gin

import (
	"io"
	"net/http"

	gh "goutil/http"

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
	return gh.ForwardResponse(ctx, client, method, url, header, body, ctx.Writer)
}

// ForwardResponseWithRequest 使用新的请求转发响应
func ForwardResponseWithRequest(ctx *gin.Context, client *http.Client, req *http.Request) error {
	return gh.ForwardResponseWithRequest(ctx, client, req, ctx.Writer)
}

// ForwardResponseWithReader 使用 reader 作为 body 构造新的请求转发响应
func ForwardResponseWithReader(ctx *gin.Context, client *http.Client, method, url string, header http.Header, body io.Reader) error {
	return gh.ForwardResponseWithReader(ctx, client, method, url, header, body, ctx.Writer)
}

// ForwardResponseWithJSONBody 使用参数构造新的 json body 请求后转发响应
func ForwardResponseWithJSONBody(ctx *gin.Context, client *http.Client, method, url string, header http.Header, body any) error {
	return gh.ForwardResponseWithJSONBody(ctx, client, method, url, header, body, ctx.Writer)
}
