package gorm

import (
	"gorm.io/gorm"
)

// Time 创建和更新时间
type Time struct {
	// 创建时间戳，单位秒
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 更新时间戳，单位秒
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// PageQuery 分页查询参数
type PageQuery struct {
	// 偏移，小于 0 不匹配
	Offset int `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count int `json:"count,omitempty" form:"count" binding:"omitempty,min=1"`
	// 排序，"column1 [desc], column2..."
	Order string `json:"order,omitempty" form:"order"`
	// 是否需要返回总数
	Total string `json:"total,omitempty" form:"total" binding:"omitempty,oneof=0 1"`
}

// HasTotal 是否有总数
func (m *PageQuery) HasTotal() bool {
	return m.Total == "1"
}

// InitDB 初始化
func (m *PageQuery) InitDB(db *gorm.DB) *gorm.DB {
	// 分页
	if m.Offset > 0 {
		db = db.Offset(m.Offset)
	}
	// 数量
	if m.Count > 0 {
		db = db.Limit(m.Count)
	}
	// 排序
	if m.Order != "" {
		db = db.Order(m.Order)
	}
	return db
}

// PageResult 是 Page 的返回值
type PageResult[M any] struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Data []M `json:"data"`
}

// Page 用于分页查询
func Page[M any](db *gorm.DB, page *PageQuery, res *PageResult[M], fields ...string) error {
	if page != nil {
		// 总数
		if page.HasTotal() {
			if err := db.Count(&res.Total).Error; err != nil {
				return err
			}
		}
		// 分页
		db = page.InitDB(db)
	}
	// 查询
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if err := db.Find(&res.Data).Error; err != nil {
		return err
	}
	//
	return nil
}

// All 用于查询全部
func All[M any](db *gorm.DB, query any) (ms []M, err error) {
	// 查询条件
	if query != nil {
		db = InitQuery(db, query)
	}
	// 查询
	err = db.Find(&ms).Error
	return
}

// AllOrder 用于查询全部
func AllOrder[M any](db *gorm.DB, query any, order string) (ms []M, err error) {
	// 查询条件
	if query != nil {
		db = InitQuery(db, query)
	}
	// 排序
	if order != "" {
		db = db.Order(order)
	}
	// 查询
	err = db.Find(&ms).Error
	return
}
