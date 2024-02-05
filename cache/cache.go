package goutil

import (
	"context"
	"goutil/sync"
)

// CacheLoadFunc 加载数据操作
type CacheLoadFunc[K comparable, M any] func(context.Context, *Cache[K, M]) error

// Cache 封装数据库缓存
type Cache[K comparable, M any] struct {
	D sync.Map[K, M]
	// 加载，被调用之前，已经加锁了
	Load CacheLoadFunc[K, M]
	// 缓存是否有效
	ok bool
}

// Init 初始化
func (c *Cache[K, M]) Init(loadFunc CacheLoadFunc[K, M]) {
	c.D.Init()
	c.Load = loadFunc
}

// Invalid 设置数据无效
func (c *Cache[K, M]) Invalid() {
	c.D.Lock()
	c.ok = false
	c.D.Unlock()
}

// BatchDel 批量删除
func (c *Cache[K, M]) BatchDel(k []K) {
	c.D.Lock()
	for i := 0; i < len(k); i++ {
		delete(c.D.D, k[i])
	}
	c.D.Unlock()
}

// DelBy 删除符合 f 的第一个返回
func (c *Cache[K, M]) DelBy(f func(m M) bool) {
	c.D.Lock()
	defer c.D.Unlock()
	for k, v := range c.D.D {
		if f(v) {
			delete(c.D.D, k)
			break
		}
	}
}

// BatchDeleteBy 删除所有符合 f 的
func (c *Cache[K, M]) BatchDelBy(f func(m M) bool) {
	c.D.Lock()
	defer c.D.Unlock()
	for k, v := range c.D.D {
		if f(v) {
			delete(c.D.D, k)
		}
	}
}

// Get 查询
func (c *Cache[K, M]) Get(ctx context.Context, k K) (m M, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	return c.D.D[k], nil
}

// GetIn 查询
func (c *Cache[K, M]) GetIn(ctx context.Context, k []K) (ms []M, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for i := 0; i < len(k); i++ {
		m, o := c.D.D[k[i]]
		if o {
			ms = append(ms, m)
		}
	}
	return
}

// Count 符合 f 的数量
func (c *Cache[K, M]) Count(ctx context.Context, f func(M) bool) (n int64, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		if f(v) {
			n++
		}
	}
	return
}

// Total 总数
func (c *Cache[K, M]) Total(ctx context.Context, k K) (n int64, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	n = int64(len(c.D.D))
	return
}

// Foreach 遍历
func (c *Cache[K, M]) Foreach(ctx context.Context, f func(M)) (err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		f(v)
	}
	return nil
}

// Search 遍历查询符合 f 的第一个返回
func (c *Cache[K, M]) Search(ctx context.Context, f func(M) bool) (m M, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		if f(v) {
			m = v
			break
		}
	}
	return
}

// BatchSearch 遍历查询所有符合 f 的
func (c *Cache[K, M]) BatchSearch(ctx context.Context, f func(M) bool) (ms []M, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		if f(v) {
			ms = append(ms, v)
		}
	}
	return
}

// SearchCache 遍历查询符合 f 的第一个返回
func SearchCache[T any, K comparable, M any](ctx context.Context, c *Cache[K, M], f func(M) (bool, T)) (m T, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		o, t := f(v)
		if o {
			m = t
			break
		}
	}
	return
}

// BatchSearchCache 遍历查询所有符合 f 的
func BatchSearchCache[T any, K comparable, M any](ctx context.Context, c *Cache[K, M], f func(M) (bool, T)) (ms []T, err error) {
	c.D.Lock()
	defer c.D.Unlock()
	if !c.ok {
		if err = c.Load(ctx, c); err != nil {
			return
		}
		c.ok = true
	}
	for _, v := range c.D.D {
		o, t := f(v)
		if o {
			ms = append(ms, t)
		}
	}
	return
}
