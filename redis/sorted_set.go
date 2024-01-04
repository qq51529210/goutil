package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ZRangeIncr 获取最小值 member 然后 + 1
func ZRangeIncr(ctx context.Context, db redis.UniversalClient, key string) (string, error) {
	res, err := redis.NewScript(`
-- 获取
local res = redis.call("ZRange", KEYS[1], 0, 0)
if (res ~= nil and #res > 0) then
	-- 预加 1 个，防止并发
	if redis.call("ZINCRBY", KEYS[1], 1, res[1]) then
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

// ZScanAllMember 获取所有 member
func ZScanAllMember(ctx context.Context, db redis.UniversalClient, key string) ([]string, error) {
	var ids []string
	n := 0
	// 扫描
	it := db.ZScan(ctx, key, 0, "*", 0).Iterator()
	for it.Next(ctx) {
		if n%2 == 0 {
			ids = append(ids, it.Val())
		}
		n++
	}
	if err := it.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}
