package gorm

import (
	"context"
	"goutil/log"
	"time"
)

// ITERATOR

// Iterator 封装循环查询代码
type Iterator[M, Q any] struct {
	// 日志
	Trace string
	// 条件
	Query Q
	// 查询函数
	QueryFunc func(context.Context, Q, *PageResult[M]) error
	// 是否继续
	InterruptFunc func() bool
	// 处理函数
	HandleFunc func(Q, []M, time.Duration) (retry bool, goon bool)
}

// Do 查询数据库然后回调，返回 true 表示全部遍历完成
func (it *Iterator[M, Q]) Do(page, retry int, dataTO time.Duration) bool {
	// 失败重试
	errCount := -1
	// 数据
	for !it.InterruptFunc() {
		// 查询
		var res PageResult[M]
		ctx, cancel := context.WithTimeout(context.Background(), dataTO)
		err := it.QueryFunc(ctx, it.Query, &res)
		cancel()
		if err != nil {
			log.ErrorfTrace(it.Trace, "get db list %v", err)
			// 重试
			errCount++
			// 重试达标
			if errCount >= retry {
				return false
			}
			// 等一秒
			time.Sleep(time.Second)
			continue
		}
		// 没有数据
		if len(res.Data) < 1 {
			return true
		}
		// 处理
		for it.InterruptFunc() {
			_retry, goon := it.HandleFunc(it.Query, res.Data, dataTO)
			if !goon {
				// 不继续
				return true
			}
			if !_retry {
				// 不重试
				break
			}
			// 重试
			errCount++
			// 重试达标
			if errCount >= retry {
				return false
			}
		}
		// 不分页/到底了
		if page < 1 || len(res.Data) < page {
			return true
		}
	}
	return false
}
