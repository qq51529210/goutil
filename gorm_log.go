package util

import (
	"context"
	"time"

	"util/log"

	"gorm.io/gorm/logger"
)

// GORMLog 用于接收 gorm 的日志
type GORMLog struct {
}

func (lg *GORMLog) LogMode(logger.LogLevel) logger.Interface {
	return lg
}

func (lg *GORMLog) Info(ctx context.Context, str string, args ...interface{}) {
	log.Info(str)
}

func (lg *GORMLog) Warn(ctx context.Context, str string, args ...interface{}) {
	log.Warn(str)
}

func (lg *GORMLog) Error(ctx context.Context, str string, args ...interface{}) {
	log.Error(str)
}

func (lg *GORMLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	log.Debugf("%s cost %v", sql, time.Since(begin))
	//
	if err != nil {
		log.Error(err)
		return
	}
}
