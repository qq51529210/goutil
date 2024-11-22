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
	lg *log.Logger
}

// NewLog 用于接收  的日志
func NewLog(lg *log.Logger) logger.Interface {
	return &Log{lg: lg}
}

// LogMode 实现接口
func (g *Log) LogMode(logger.LogLevel) logger.Interface {
	return g
}

// Info 实现接口
func (g *Log) Info(ctx context.Context, str string, args ...interface{}) {
	g.lg.Info("", 0, "", str)
}

// Warn 实现接口
func (g *Log) Warn(ctx context.Context, str string, args ...interface{}) {
	g.lg.Warn("", 0, "", str)
}

// Error 实现接口
func (g *Log) Error(ctx context.Context, str string, args ...interface{}) {
	g.lg.Error("", 0, "", str)
}

// Trace 实现接口
func (g *Log) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	g.lg.Debug("", time.Since(begin), "", sql)
	//
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			g.lg.Error("", 0, err)
			return
		}
	}
}
