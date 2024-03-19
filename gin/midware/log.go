package gin

import (
	"fmt"
	"goutil/log"
	"goutil/uid"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Log 日志中间件
type Log struct {
	// Logger 日志
	Logger *log.Logger
	// CtxKeyTraceID 用于追踪 id
	CtxKeyTraceID string
	// CtxKeyRequestData 用于保存请求的的数据
	CtxKeyRequestData string
	// CtxKeyResponseData 用于保存响应的数据
	CtxKeyResponseData string
	// CtxKeyError 用于保存处理中发生的错误
	CtxKeyHandleError string
	// HeaderNameRemoteAddr 代理服务透传的客户端地址头名称
	HeaderNameRemoteAddr string
	// HeaderNameTraceID 如果有则使用它
	HeaderNameTraceID string
}

// 实现接口
func (h *Log) ServeHTTP(ctx *gin.Context) {
	// 清理
	defer func() {
		// 异常
		if h.Logger.Recover(recover()) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	// 追踪 id
	var traceID string
	if h.HeaderNameTraceID != "" {
		traceID = ctx.GetHeader(h.HeaderNameTraceID)
	}
	if traceID == "" {
		traceID = uid.SnowflakeIDString()
	}
	ctx.Set(h.CtxKeyTraceID, traceID)
	old := time.Now()
	// 执行
	ctx.Next()
	// 花费时间
	cost := time.Since(old)
	// 如果有代理，代理必须使用这个字段来透传客户端 ip
	remoteAddr := ctx.Request.RemoteAddr
	if addr := ctx.GetHeader(h.HeaderNameRemoteAddr); addr != "" {
		remoteAddr = addr
	}
	// 日志
	var str strings.Builder
	fmt.Fprintf(&str, "[%v] %s %s %s", cost, remoteAddr, ctx.Request.Method, ctx.Request.URL.Path)
	// 提交的数据
	if data, ok := ctx.Value(h.CtxKeyRequestData).(string); ok && data != "" {
		str.WriteString("\nrequest data: ")
		str.WriteString(data)
	}
	// 返回的数据
	if data, ok := ctx.Value(h.CtxKeyResponseData).(string); ok && data != "" {
		str.WriteString("\nresponse data: ")
		str.WriteString(data)
	}
	// 如果有错误
	if data, ok := ctx.Value(h.CtxKeyHandleError).(string); ok && data != "" {
		str.WriteString("\nhandle error: ")
		str.WriteString(data)
	}
	// 输出
	h.Logger.DebugTrace(traceID, str.String())
}
