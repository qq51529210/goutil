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
func (m *Model[V]) Get(ctx context.Context, v V) error {
	return m.D.WithContext(ctx).Take(v).Error
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
func (m *Model[V]) Update(ctx context.Context, v V) (int64, error) {
	db := m.D.WithContext(ctx).Updates(v)
	return db.RowsAffected, db.Error
}

// UpdateFields 更新指定字段
func (m *Model[V]) UpdateFields(ctx context.Context, v V, fs ...string) (int64, error) {
	db := m.D.WithContext(ctx).Select(fs).Updates(v)
	return db.RowsAffected, db.Error
}

// BatchUpdate 批量更新
func (m *Model[V]) BatchUpdate(ctx context.Context, vs []V) (int64, error) {
	var row int64
	return row, m.D.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range vs {
			db := tx.Updates(v)
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
