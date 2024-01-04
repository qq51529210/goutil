package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// ScanFirst 扫描并返回第一个
func ScanFirst(ctx context.Context, db redis.UniversalClient, scanKey string) (string, error) {
	// 用迭代器
	it := db.Scan(ctx, 0, scanKey, 0).Iterator()
	for it.Next(ctx) {
		return it.Val(), nil
	}
	return "", it.Err()
}
