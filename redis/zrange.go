package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

func ZRangeIncr(ctx context.Context, db redis.UniversalClient, key string) (string, error) {
	res, err := redis.NewScript(`
-- 参数
local key = KEYS[1]
-- 获取
local res = redis.call("ZRange", key, 0, 0)
if (res ~= nil and #res > 0) then
	-- 预加 1 个，防止并发
	if redis.call("ZINCRBY", key, 1, res[1]) then
		return res[1]
	end
end
-- 返回
return ""
`).Run(context.Background(), db, []string{key}).Result()
	if err != nil {
		return "", err
	}
	return res.(string), nil
}
