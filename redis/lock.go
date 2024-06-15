package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Lock 分布式锁
func Lock(ctx context.Context, db redis.UniversalClient, key, value string, expire time.Duration) (bool, error) {
	cmd := redis.NewScript(`
-- 参数
local key, value, expire = KEYS[1], ARGV[1], ARGV[2]
-- 抢锁
local result = redis.call("SET", key, value, "NX", "EX", expire)
-- 抢到返回
if result then
	return 1
end
-- 没抢到
result = redis.call("GET", key)
-- 检查是否自己
if result == value then
	-- 更新过期时间
	redis.call("EXPIRE", key, expire)
	return 1
end
-- 不是自己
return 0
`).Run(ctx, db, []string{key}, value, expire/time.Second)
	err := cmd.Err()
	if err != nil {
		return false, err
	}
	n, _ := cmd.Int()
	return n == 1, nil
}
