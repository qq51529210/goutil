package redis

import "github.com/redis/go-redis/v9"

// IsDataNotFound 是否没有数据
func IsDataNotFound(err error) bool {
	return err == redis.Nil
}
