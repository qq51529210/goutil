package gin

import (
	"net/http"

	gh "goutil/http"

	"github.com/gin-gonic/gin"
)

// Forward 转发请求和响应，baseURL 是 http(s)://host:port
func Forward(ctx *gin.Context, client *http.Client, baseURL string) error {
	return gh.Forward(ctx, client, baseURL, ctx.Request, ctx.Writer)
}
