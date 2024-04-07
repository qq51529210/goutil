package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client[M any] struct {
	db redis.UniversalClient
}

func NewClient[M any](db redis.UniversalClient) *Client[M] {
	c := new(Client[M])
	c.db = db
	return c
}

func (c *Client[M]) HGet(ctx context.Context, key string, m M, fields ...string) error {
	if len(fields) < 1 {
		cmd := c.db.HGetAll(ctx, key)
		data, err := cmd.Result()
		if err != nil {
			return err
		}
		if len(data) > 0 {
			return cmd.Scan(m)
		}
	} else {
		cmd := c.db.HMGet(ctx, key, fields...)
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
