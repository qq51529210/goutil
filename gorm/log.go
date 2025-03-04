package gorm

import (
	"context"
	"goutil/log"
	"time"

	"gorm.io/gorm/logger"
)

// Log 用于接收 gorm 日志
type Log struct {
	lg *log.Logger
	tk string
}

// NewLog 用于接收  的日志
func NewLog(lg *log.Logger, traceKey string) logger.Interface {
	return &Log{lg: lg, tk: traceKey}
}

// LogMode 实现接口
func (g *Log) LogMode(logger.LogLevel) logger.Interface {
	return g
}

// Info 实现接口
func (g *Log) Info(ctx context.Context, str string, args ...interface{}) {
	trace, _ := ctx.Value(g.tk).(string)
	g.lg.Info(-1, trace, 0, str)
}

// Warn 实现接口
func (g *Log) Warn(ctx context.Context, str string, args ...interface{}) {
	trace, _ := ctx.Value(g.tk).(string)
	g.lg.Warn(-1, trace, 0, str)
}

// Error 实现接口
func (g *Log) Error(ctx context.Context, str string, args ...interface{}) {
	trace, _ := ctx.Value(g.tk).(string)
	g.lg.Error(-1, trace, 0, str)
}

// Trace 实现接口
func (g *Log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	//
	trace, _ := ctx.Value(g.tk).(string)
	if err != nil {
		g.lg.Errorf(-1, trace, time.Since(begin), "%s error: %s", sql, err.Error())
	} else {
		g.lg.Debug(-1, trace, time.Since(begin), sql)
	}
}
