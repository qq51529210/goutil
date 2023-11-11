package goutil

import (
	"context"
	"sync"
)

type loader struct {
	c   chan struct{}
	err error
}

// Cache 封装数据库缓存
type Cache[K comparable, M any] struct {
	// 同步锁
	sync.Mutex
	// 数据
	D map[K]M
	// 缓存是否有效
	ok bool
	// 加载所有数据
	loadFunc func() (map[K]M, error)
	// 是否正在加载缓存
	loading bool
	// 协程加载
	loader *loader
}

// Init 初始化
func (c *Cache[K, M]) Init(loadFunc func() (map[K]M, error)) {
	c.D = make(map[K]M)
	c.loadFunc = loadFunc
}

// Query 首先检查数据是否有效，无效则重新加载
func (c *Cache[K, M]) Query(ctx context.Context, query func(c *Cache[K, M])) error {
	c.Lock()
	// 无效，重新加载所有
	if !c.ok {
		var wait *loader
		// 确保并发下只有一个加载操作
		if !c.loading {
			c.loading = true
			c.loader = new(loader)
			c.loader.c = make(chan struct{})
			go c.loadAllRoutine()
		}
		wait = c.loader
		c.Unlock()
		// 等待结果
		select {
		case <-ctx.Done():
			// 上下文结果
			return ctx.Err()
		case <-wait.c:
			// 加载失败，返回
			if wait.err != nil {
				return wait.err
			}
			// 加载成功，往下
		}
	}
	// 有效，查询
	query(c)
	c.Unlock()
	return nil
}

// loadAllRoutine 在协程中加载数据
func (c *Cache[K, M]) loadAllRoutine() {
	data, err := c.loadFunc()
	c.Lock()
	c.loader.err = err
	if err == nil {
		c.D = data
		c.ok = true
	}
	c.loading = false
	// 通知
	close(c.loader.c)
	c.Unlock()
}

// Invalidate 设置数据无效，这会导致下一次操作重新加载数据
func (c *Cache[K, M]) Invalidate() {
	c.Lock()
	c.ok = false
	c.D = make(map[K]M)
	c.Unlock()
}

// Add 添加
func (c *Cache[K, M]) Add(k K, m M) {
	c.Lock()
	c.D[k] = m
	c.Unlock()
}

// Update 更新，f 回调的 m 可能为零值
func (c *Cache[K, M]) Update(k K, f func(m M)) {
	c.Lock()
	f(c.D[k])
	c.Unlock()
}

// Delete 删除
func (c *Cache[K, M]) Delete(k K) {
	c.Lock()
	delete(c.D, k)
	c.Unlock()
}

// BatchDelete 批量删除
func (c *Cache[K, M]) BatchDelete(k []K) {
	c.Lock()
	for i := 0; i < len(k); i++ {
		delete(c.D, k[i])
	}
	c.Unlock()
}

// DeleteBy 删除 f 返回 true 的数据，然后返回
func (c *Cache[K, M]) DeleteBy(f func(m M) bool) {
	c.Lock()
	for k, v := range c.D {
		if f(v) {
			delete(c.D, k)
			break
		}
	}
	c.Unlock()
}

// BatchDeleteBy 循环删除 f 返回 true 的数据
func (c *Cache[K, M]) BatchDeleteBy(f func(m M) bool) {
	c.Lock()
	for k, v := range c.D {
		if f(v) {
			delete(c.D, k)
		}
	}
	c.Unlock()
}

// GetKeys 返回所有 key
func (c *Cache[K, M]) GetKeys(ctx context.Context) (ks []K, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for k := range c.D {
			ks = append(ks, k)
		}
	})
	return
}

// Get 查询
func (c *Cache[K, M]) Get(ctx context.Context, k K) (m M, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		m = c.D[k]
	})
	return
}

// GetIn 查询
func (c *Cache[K, M]) GetIn(ctx context.Context, k []K) (ms []M, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for i := 0; i < len(k); i++ {
			m, o := c.D[k[i]]
			if o {
				ms = append(ms, m)
			}
		}
	})
	return
}

// Count 符合 f 的数量
func (c *Cache[K, M]) Count(ctx context.Context, f func(M) bool) (m int64, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for _, v := range c.D {
			if f(v) {
				m++
			}
		}
	})
	return
}

// Total 总数
func (c *Cache[K, M]) Total(ctx context.Context, k K) (m int64, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		m = int64(len(c.D))
	})
	return
}

// Foreach 遍历
func (c *Cache[K, M]) Foreach(ctx context.Context, f func(M)) error {
	return c.Query(ctx, func(c *Cache[K, M]) {
		for _, v := range c.D {
			f(v)
		}
	})
}

// Search 遍历查询
func (c *Cache[K, M]) Search(ctx context.Context, f func(M) bool) (ms []M, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for _, v := range c.D {
			if f(v) {
				ms = append(ms, v)
			}
		}
	})
	return
}

// SearchOne 遍历查询，f 返回 true 就返回
func (c *Cache[K, M]) SearchOne(ctx context.Context, f func(M) bool) (m M, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for _, v := range c.D {
			if f(v) {
				m = v
				return
			}
		}
	})
	return
}

// SearchCache 遍历查询，返回其他的类型
func SearchCache[T any, K comparable, M any](ctx context.Context, c *Cache[K, M], f func(M) (bool, T)) (ms []T, err error) {
	err = c.Query(ctx, func(c *Cache[K, M]) {
		for _, v := range c.D {
			o, m := f(v)
			if o {
				ms = append(ms, m)
			}
		}
	})
	return
}
