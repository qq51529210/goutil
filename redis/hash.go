package redis

import (
	"context"
	"time"

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
			ok := false
			for i := 0; i < len(data); i++ {
				if data[i] != nil {
					ok = true
					break
				}
			}
			if ok {
				return cmd.Scan(m)
			}
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
		return nil, _cursor, err
	}
	// 根据 keys 再查询
	var ms []*M
	for i := 0; i < len(keys); i++ {
		m := new(M)
		err := HGet(ctx, db, keys[i], m, fields...)
		if err != nil {
			if err != redis.Nil {
				return nil, _cursor, err
			}
			continue
		}
		ms = append(ms, m)
	}
	// 返回
	return ms, _cursor, nil
}

// HSetNxEx 存在才设置
func HSetNxEx(ctx context.Context, db redis.UniversalClient, key string, expire time.Duration, field string, value any) (bool, error) {
	res, err := redis.NewScript(`
local result = redis.call("EXISTS", KEYS[1])
-- 存在
if result == 1 then
	-- 更新
	redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
	redis.call("EXPIRE", KEYS[1], ARGV[3])
end
-- 返回
return result
`).Run(ctx, db, []string{key}, field, value, expire/time.Second).Result()
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

// HSetNx 存在才设置
func HSetNx(ctx context.Context, db redis.UniversalClient, key, field string, value any) (bool, error) {
	res, err := redis.NewScript(`
local result = redis.call("EXISTS", KEYS[1])
-- 存在
if result == 1 then
	-- 更新
	redis.call("HSET", KEYS[1], ARGV[1], ARGV[2])
end
-- 返回
return result
`).Run(ctx, db, []string{key}, field, value).Result()
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}

// HMSetNxEx 存在才设置
func HMSetNxEx(ctx context.Context, db redis.UniversalClient, key string, expire time.Duration, args ...any) (bool, error) {
	var values []any
	values = append(values, expire/time.Second)
	values = append(values, args...)
	//
	res, err := redis.NewScript(`
local result = redis.call("EXISTS", KEYS[1])
-- 存在
if result == 1 then
	-- 更新
	redis.call("HSET", KEYS[1], ARGV[2], ARGV[3])
	redis.call("EXPIRE", KEYS[1], ARGV[1])
end
-- 返回
return result
`).Run(ctx, db, []string{key}, values...).Result()
	if err != nil {
		return false, err
	}
	return res.(int64) == 1, nil
}
