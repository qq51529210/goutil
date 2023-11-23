package gin

import (
	"encoding/json"
	"fmt"
	"goutil/log"
	"goutil/uid"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	// CtxKeySubmitData 用于 log 保存提交的 body 的数据
	CtxKeySubmitData = "SumbitData"
	// CtxKeyError 用于 log 保存处理中发生的错误
	CtxKeyError = "Error"
	// CtxKeyTraceID 用于 log 追踪 id
	CtxKeyTraceID = "TraceID"
)

var (
	// ProxyRemoteAddrHeader 代理服务透传的客户端地址头名称
	ProxyRemoteAddrHeader = "X-Remote-Addr"
	// Logger 日志
	Logger *log.Logger
)

// Log 日志中间件
func Log(ctx *gin.Context) {
	old := time.Now()
	// 追踪 id
	traceID := uid.SnowflakeIDString()
	ctx.Set(CtxKeyTraceID, traceID)
	// 清理
	defer func() {
		// 异常
		r := recover()
		if r != nil {
			if Logger != nil {
				Logger.Recover(r)
			}
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
	}()
	// 执行
	ctx.Next()
	// 日志
	if Logger == nil {
		return
	}
	// 花费时间
	cost := time.Since(old)
	// 如果有代理，代理必须使用这个字段来透传客户端 ip
	remoteAddr := ctx.Request.RemoteAddr
	if addr := ctx.GetHeader(ProxyRemoteAddrHeader); addr != "" {
		remoteAddr = addr
	}
	// 日志
	var str strings.Builder
	fmt.Fprintf(&str, "%s %s %s cost %v", remoteAddr, ctx.Request.Method, ctx.Request.URL.Path, cost)
	// 提交的数据
	submitData := ctx.Value(CtxKeySubmitData)
	if submitData != nil {
		d, err := json.Marshal(submitData)
		if err == nil {
			str.WriteString("\nsubmit data: ")
			str.Write(d)
		}
	}
	// 如果有错误
	errData := ctx.Value(CtxKeyError)
	if errData != nil {
		fmt.Fprintf(&str, "\nhandle error: %v", errData)
	}
	Logger.DebugTrace(traceID, str.String())
}
