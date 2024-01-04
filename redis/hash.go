package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HGet 查询哈希，不指定字段则查询全部
func HGet[M any](ctx context.Context, db redis.UniversalClient, key string, m *M, fields ...string) error {
	if len(fields) < 1 {
		cmd := db.HGetAll(ctx, key)
		data, err := cmd.Result()
		if err != nil {
			return err
		}
		if len(data) > 0 {
			return cmd.Scan(m)
		}
	} else {
		cmd := db.HMGet(ctx, key, fields...)
		data, err := cmd.Result()
		if err != nil {
			return err
		}
		if len(data) > 0 {
			return cmd.Scan(m)
		}
	}
	// 没有数据
	return redis.Nil
}

// PageHGet 分页查询哈希，不指定字段则查询全部，注意 redis 这个包无法 scan 到指针字段
func PageHGet[M any](ctx context.Context, db redis.UniversalClient, key string, cursor, count uint64, fields ...string) ([]*M, uint64, error) {
	// 扫描 key
	cmd := db.Scan(ctx, cursor, key, int64(count))
	keys, _cursor, err := cmd.Result()
	if err != nil {
		return nil, cursor, err
	}
	// 根据 keys 再查询
	var ms []*M
	for i := 0; i < len(keys); i++ {
		m := new(M)
		err := HGet(ctx, db, keys[i], m, fields...)
		if err != nil {
			if err != redis.Nil {
				return nil, cursor, err
			}
			continue
		}
		ms = append(ms, m)
	}
	// 返回
	return ms, _cursor, nil
}
