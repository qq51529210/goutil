package util

import (
	"context"
	"sync"

	"gorm.io/gorm"
)

// GORMCache 用于缓存数据
type GORMCache[K comparable, M any] struct {
	sync.Mutex
	// 数据库
	db *gorm.DB
	// 是否开启缓存
	cache bool
	// 数据
	D map[K]M
	// gorm.Model 要的模型
	M M
	// 标记数据 d 是否需要重新加载
	ok bool
	// 创建函数
	new func() M
	// 返回 M 的主键
	key func(M) K
	// 返回 M 的主键，用于单个，查询，修改，删除
	whereKey func(*gorm.DB, K) *gorm.DB
	// 返回 M 的主键列表，用于批量删除
	whereKeys func(*gorm.DB, []K) *gorm.DB
}

// NewGORMCache 返回新的缓存，enable 为 false 则不开启缓存
func NewGORMCache[K comparable, M any](
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) *GORMCache[K, M] {
	c := new(GORMCache[K, M])
	c.Init(db, cache, newFunc, keyFunc, whereKeyFunc, whereKeysFunc)
	return c
}

// Init 初始化字段
func (c *GORMCache[K, M]) Init(
	db *gorm.DB,
	cache bool,
	newFunc func() M,
	keyFunc func(M) K,
	whereKeyFunc func(*gorm.DB, K) *gorm.DB,
	whereKeysFunc func(*gorm.DB, []K) *gorm.DB,
) {
	c.db = db
	c.cache = cache
	c.D = make(map[K]M)
	c.new = newFunc
	c.key = keyFunc
	c.whereKey = whereKeyFunc
	c.whereKeys = whereKeysFunc
	c.M = newFunc()
}

// IsCache 返回是否启用
func (c *GORMCache[K, M]) IsCache() bool {
	return c.cache
}

// DB 返回 db
func (c *GORMCache[K, M]) DB() *gorm.DB {
	return c.db
}

// Model 返回加载模型的 db
func (c *GORMCache[K, M]) Model(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx).Model(c.M)
}

// LoadWhere 加载并替换指定条件的数据，同步
func (c *GORMCache[K, M]) LoadWhere(ctx context.Context, whereFunc func(db *gorm.DB) *gorm.DB) error {
	// 不启用
	if !c.cache {
		return nil
	}
	//
	db := c.Model(ctx)
	c.Lock()
	err := c.loadWhere(db)
	c.Unlock()
	//
	return err
}

// loadWhere 加载并替换指定条件的数据
func (c *GORMCache[K, M]) loadWhere(db *gorm.DB) error {
	// 查询
	var models []M
	err := db.Find(&models).Error
	if err != nil {
		c.ok = false
		return err
	}
	// 加载或替换
	for _, model := range models {
		c.D[c.key(model)] = model
	}
	//
	return nil
}

// Check 检查内存数据是否需要重新加载，同步
func (c *GORMCache[K, M]) LoadAll(ctx context.Context) error {
	// 不启用
	if !c.cache {
		return nil
	}
	// 加载
	c.Lock()
	err := c.loadAll(ctx)
	c.Unlock()
	//
	return err
}

// loadAll 检查内存数据是否需要重新加载
func (c *GORMCache[K, M]) loadAll(ctx context.Context) error {
	// 数据是否有效
	if !c.ok {
		var models []M
		// 数据库
		err := c.Model(ctx).Find(&models).Error
		if err != nil {
			c.ok = false
			return err
		}
		// 重置
		c.D = make(map[K]M)
		for _, model := range models {
			c.D[c.key(model)] = model
		}
		c.ok = true
	}
	return nil
}

// Load 尝试加载，添加和修改时候调用，同步
func (c *GORMCache[K, M]) Load(ctx context.Context, k K) {
	// 不启用
	if !c.cache {
		return
	}
	// 加载
	c.Lock()
	c.load(ctx, k)
	c.Unlock()
}

// load 尝试加载，添加和修改时候调用
func (c *GORMCache[K, M]) load(ctx context.Context, k K) error {
	// 读取
	m := c.new()
	err := c.whereKey(c.Model(ctx), k).First(m).Error
	// 失败
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.ok = true
			return nil
		}
		c.ok = false
		return err
	}
	// 成功
	c.D[k] = m
	c.ok = true
	//
	return nil
}

// All 返回所有内存，不要修改返回的指针，同步
func (c *GORMCache[K, M]) All(ctx context.Context) ([]M, error) {
	// 不启用
	if !c.cache {
		var models []M
		err := c.Model(ctx).Find(&models).Error
		if err != nil {
			return nil, err
		}
		return models, nil
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return nil, err
	}
	// 列表
	var models []M
	for _, v := range c.D {
		models = append(models, v)
	}
	//
	return models, nil
}

// Get 返回指定，不要修改返回的指针，同步
func (c *GORMCache[K, M]) Get(ctx context.Context, k K) (m M, err error) {
	// 不启用
	if !c.cache {
		mm := c.new()
		err = c.whereKey(c.Model(ctx), k).First(mm).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return
		}
		m = mm
		//
		return
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保加载
	err = c.loadAll(ctx)
	if err != nil {
		return
	}
	m = c.D[k]
	//
	return
}

// Add 添加，同步
func (c *GORMCache[K, M]) Add(ctx context.Context, m M) (int64, error) {
	// 数据库
	db := c.Model(ctx).Create(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(ctx, c.key(m))
		c.Unlock()
	}
	//
	return db.RowsAffected, nil
}

// Update 更新，同步
func (c *GORMCache[K, M]) Update(ctx context.Context, m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.Model(ctx), k).Updates(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(ctx, c.key(m))
		c.Unlock()
	}
	//
	return db.RowsAffected, nil
}

// UpdateCache 更新内存，同步。回调有可能为 nil
func (c *GORMCache[K, M]) UpdateCache(k K, fn func(M)) {
	// 不启用
	if !c.cache {
		return
	}
	//
	c.Lock()
	fn(c.D[k])
	c.Unlock()
}

// Save 保存，同步
func (c *GORMCache[K, M]) Save(ctx context.Context, m M) (int64, error) {
	// 数据库
	k := c.key(m)
	db := c.whereKey(c.Model(ctx), k).Save(m)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		c.load(ctx, c.key(m))
		c.Unlock()
	}
	//
	return db.RowsAffected, nil
}

// BatchSave 事务保存，同步
func (c *GORMCache[K, M]) BatchSave(ctx context.Context, ms []M) (int64, error) {
	// 数据库
	err := c.Model(ctx).Transaction(func(tx *gorm.DB) error {
		for _, m := range ms {
			db := tx.Save(m)
			if db.Error != nil {
				return db.Error
			}
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	// 内存
	c.LoadAll(ctx)
	//
	return int64(len(ms)), nil
}

// Delete 删除，同步
func (c *GORMCache[K, M]) Delete(ctx context.Context, k K) (int64, error) {
	// 数据库
	db := c.whereKey(c.Model(ctx), k).Delete(c.M)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		delete(c.D, k)
		c.Unlock()
	}
	//
	return db.RowsAffected, nil
}

// DeleteCache 删除内存
func (c *GORMCache[K, M]) DeleteCache(k K) {
	// 不启用
	if !c.cache {
		return
	}
	c.Lock()
	delete(c.D, k)
	c.Unlock()
}

// BatchDelete 批量删除，同步
func (c *GORMCache[K, M]) BatchDelete(ctx context.Context, ks []K) (int64, error) {
	// 数据库
	db := c.whereKeys(c.Model(ctx), ks).Delete(c.M)
	if db.Error != nil {
		return db.RowsAffected, db.Error
	}
	// 内存
	if c.cache && db.RowsAffected > 0 {
		c.Lock()
		for _, k := range ks {
			delete(c.D, k)
		}
		c.Unlock()
	}
	//
	return db.RowsAffected, nil
}

// BatchDeleteCache 删除内存，同步
func (c *GORMCache[K, M]) BatchDeleteCache(ks []K) {
	// 不启用
	if !c.cache {
		return
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 删除
	for i := 0; i < len(ks); i++ {
		delete(c.D, ks[i])
	}
}

// ForeachCache 遍历缓存吗，同步
func (c *GORMCache[K, M]) ForeachCache(ctx context.Context, fc func(M)) error {
	// 不启用
	if !c.cache {
		return nil
	}
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return err
	}
	// 循环
	for _, m := range c.D {
		fc(m)
	}
	//
	return nil
}

// Search 在内存中查找，同步
func (c *GORMCache[K, M]) SearchCache(ctx context.Context, match func(M) bool) ([]M, error) {
	var mm []M
	// 启用
	if c.cache {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 确保数据
		err := c.loadAll(ctx)
		if err != nil {
			return nil, err
		}
		// 查找
		for _, m := range c.D {
			if match(m) {
				mm = append(mm, m)
			}
		}
	}
	//
	return mm, nil
}

// SearchCacheIn 在内存中查找，同步
func (c *GORMCache[K, M]) SearchCacheIn(ctx context.Context, ks []K) ([]M, error) {
	var mm []M
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return nil, err
	}
	// 查找
	for i := 0; i < len(ks); i++ {
		m, ok := c.D[ks[i]]
		if ok {
			mm = append(mm, m)
		}
	}
	//
	return mm, nil
}

// Search 在内存中查找，同步
func (c *GORMCache[K, M]) SearchCacheOne(ctx context.Context, match func(M) bool) (m M, err error) {
	// 启用
	if c.cache {
		// 上锁
		c.Lock()
		defer c.Unlock()
		// 确保数据
		err = c.loadAll(ctx)
		if err != nil {
			return
		}
		// 查找
		for _, v := range c.D {
			if match(v) {
				m = v
				break
			}
		}
	}
	//
	return
}

// Count 返回内存匹配数量，同步
func (c *GORMCache[K, M]) CacheCount(ctx context.Context, match func(M) bool) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return 0, err
	}
	// 查找
	var n int64
	for _, v := range c.D {
		if match(v) {
			n++
		}
	}
	//
	return n, nil
}

// Total 返回内存总量，同步
func (c *GORMCache[K, M]) CacheTotal(ctx context.Context) (int64, error) {
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return 0, err
	}
	//
	return int64(len(c.D)), nil
}

// List 查询数据库
func (c *GORMCache[K, M]) List(ctx context.Context, page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	return gormList(c.Model(ctx), page, query, res)
}

// GORMSearchCache 模板化的 Search
func GORMSearchCache[T any, K comparable, M any](ctx context.Context, c *GORMCache[K, M], match func(M) (bool, T)) ([]T, error) {
	var vv []T
	// 上锁
	c.Lock()
	defer c.Unlock()
	// 确保数据
	err := c.loadAll(ctx)
	if err != nil {
		return nil, err
	}
	// 查找
	for _, m := range c.D {
		o, v := match(m)
		if o {
			vv = append(vv, v)
		}
	}
	//
	return vv, nil
}
