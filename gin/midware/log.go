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

// Log 上下文 key
const (
	LogCtxKeySubmitData     = "SubmitData"
	LogCtxKeyResponseData   = "ResponseData"
	LogCtxKeyError          = "Error"
	LogCtxKeyTraceID        = "TraceID"
	LogHeaderNameRemoteAddr = "X-Remote-Addr"
)

// Log 日志中间件
type Log struct {
	// CtxKeySubmitData 用于保存提交的 body 的数据
	CtxKeySubmitData string
	// CtxKeyResponseData 用于保存返回的 body 的数据
	CtxKeyResponseData string
	// CtxKeyError 用于保存处理中发生的错误
	CtxKeyError string
	// CtxKeyTraceID 用于追踪 id
	CtxKeyTraceID string
	// HeaderNameRemoteAddr 代理服务透传的客户端地址头名称
	HeaderNameRemoteAddr string
	// Logger 日志
	Logger *log.Logger
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
	traceID := uid.SnowflakeIDString()
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
	submitData := ctx.Value(h.CtxKeySubmitData)
	if submitData != nil {
		d, err := json.Marshal(submitData)
		if err == nil {
			str.WriteString("\nsubmit data: ")
			str.Write(d)
		}
	}
	// 返回的数据
	responseData := ctx.Value(h.CtxKeyResponseData)
	if responseData != nil {
		d, err := json.Marshal(responseData)
		if err == nil {
			str.WriteString("\nresponse data: ")
			str.Write(d)
		}
	}
	// 如果有错误
	errData := ctx.Value(h.CtxKeyError)
	if errData != nil {
		fmt.Fprintf(&str, "\nhandle error: %v", errData)
	}
	h.Logger.DebugTrace(traceID, str.String())
}
