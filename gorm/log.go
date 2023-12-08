package gorm

import (
	"context"
	"goutil/log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Log 用于接收 gorm 日志
type Log struct {
	traceID string
	lg      *log.Logger
}

// NewLog 用于接收  的日志
func NewLog(traceID string, lg *log.Logger) logger.Interface {
	return &Log{traceID: traceID, lg: lg}
}

// LogMode 实现接口
func (g *Log) LogMode(logger.LogLevel) logger.Interface {
	return g
}

// Info 实现接口
func (g *Log) Info(ctx context.Context, str string, args ...interface{}) {
	g.lg.InfoTrace(g.traceID, str)
}

// Warn 实现接口
func (g *Log) Warn(ctx context.Context, str string, args ...interface{}) {
	g.lg.WarnTrace(g.traceID, str)
}

// Error 实现接口
func (g *Log) Error(ctx context.Context, str string, args ...interface{}) {
	g.lg.ErrorTrace(g.traceID, str)
}

// Trace 实现接口
func (g *Log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	if err != nil {
		g.lg.DebugfTrace(g.traceID, "[%v] %s", time.Since(begin), sql)
		if err != gorm.ErrRecordNotFound {
			g.lg.ErrorTrace(g.traceID, err)
			return
		}
	}
}
