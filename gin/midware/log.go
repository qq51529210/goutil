package gin

import (
	"goutil/log"
	"goutil/uid"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// LogLogger 日志对象，必须设置
	LogLogger *log.Logger
	// 上下文拿到追踪标识的 key
	LogCtxTraceID = "TraceID"
	// 拿到上层服务透传的追踪标识的 header name
	LogHeaderTrace = "X-Trace"
	// 拿到上层服务透传的客户端地址的 header name
	LogHeaderRemoteAddr = "X-Remote-Addr"
)

// Logger 全局中间件，一般放在第一
// 拿到/生成追踪标识，放到 ctx.Value(LogCtxTraceID)
// 如果有 LogHeaderRemoteAddr 的值，修改 ctx.Request.RemoteAddr
// LogLogger 必须设置，不然空指针
func Logger(ctx *gin.Context) {
	defer func() {
		// 异常
		if LogLogger.Recover(recover()) {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	// 追踪 id
	traceID := ctx.GetHeader(LogHeaderTrace)
	if traceID == "" {
		traceID = uid.SnowflakeIDString()
	}
	ctx.Set(LogCtxTraceID, traceID)
	// 如果有代理，代理必须使用这个字段来透传客户端 ip
	if addr := ctx.GetHeader(LogHeaderRemoteAddr); addr != "" {
		ctx.Request.RemoteAddr = addr
	}
	// 日志
	LogLogger.DebugfTrace(traceID, "%s %s %s", ctx.Request.RemoteAddr, ctx.Request.Method, ctx.Request.URL.Path)
	// 执行
	old := time.Now()
	ctx.Next()
	// 日志
	LogLogger.DebugfTrace(traceID, "cost %v", time.Since(old))
}
