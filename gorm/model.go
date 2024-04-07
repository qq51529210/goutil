package gorm

import (
	"context"

	"gorm.io/gorm"
)

// Model 模型，抽象
type Model[V any] struct {
	M V
	D *gorm.DB
}

// NewModel 构造
func NewModel[V any](db *gorm.DB, m V) *Model[V] {
	return &Model[V]{
		D: db,
		M: m,
	}
}

// Get 单个
func (m *Model[V]) Get(ctx context.Context, v V, fields ...string) error {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	return db.Take(v).Error
}

// First 第一个
func (m *Model[V]) First(ctx context.Context, v V, q any, fields ...string) error {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if q != nil {
		db = InitQuery(db, q)
	}
	return db.First(v).Error
}

// Add 添加
func (m *Model[V]) Add(ctx context.Context, v V) error {
	return m.D.WithContext(ctx).Create(v).Error
}

// BatchAdd 批量添加
func (m *Model[V]) BatchAdd(ctx context.Context, vs []V) error {
	return m.D.WithContext(ctx).Create(vs).Error
}

// Update 更新
func (m *Model[V]) Update(ctx context.Context, v V, fields ...string) (int64, error) {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Updates(v)
	return db.RowsAffected, db.Error
}

// BatchUpdate 批量更新
func (m *Model[V]) BatchUpdate(ctx context.Context, vs []V, fields ...string) (int64, error) {
	var row int64
	return row, m.D.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range vs {
			db := tx
			if len(fields) > 0 {
				db = db.Select(fields)
			}
			db = db.Updates(v)
			if db.Error != nil {
				return db.Error
			}
			row += db.RowsAffected
		}
		return nil
	})
}

// Delete 删除
func (m *Model[V]) Delete(ctx context.Context, v V) (int64, error) {
	db := m.D.WithContext(ctx).Delete(v)
	return db.RowsAffected, db.Error
}

// BatchDelete 批量删除
func (m *Model[V]) BatchDelete(ctx context.Context, query any) (int64, error) {
	db := InitQuery(m.D.WithContext(ctx), query).Delete(m.M)
	return db.RowsAffected, db.Error
}

// Page 分页
func (m *Model[V]) Page(ctx context.Context, page *PageQuery, query any, res *PageResult[V]) error {
	db := m.D.WithContext(ctx).Model(m.M)
	if query != nil {
		db = InitQuery(db, query)
	}
	return Page(db, page, res)
}

// All 所有
func (m *Model[V]) All(ctx context.Context, query any) ([]V, error) {
	return All[V](m.D.WithContext(ctx).Model(m.M), query)
}
