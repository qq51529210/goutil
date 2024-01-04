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
-- 返回值
local result = 0
-- 设置锁
if not redis.call("SET", key, value, "NX", "EX", expire) then
	-- 失败，检查是否自己
	if value == redis.call("GET", key) then
		result = 1
		-- 更新过期时间
		redis.call("EXPIRE", key, expire)
	end
end
-- 返回
return result
`).Run(ctx, db, []string{key}, value, expire/time.Second)
	err := cmd.Err()
	if err != nil {
		return false, err
	}
	n, _ := cmd.Int()
	return n == 1, nil
}
