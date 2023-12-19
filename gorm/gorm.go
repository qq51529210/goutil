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
	Offset *int `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `json:"count,omitempty" form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `json:"order,omitempty" form:"order"`
	// 是否需要返回总数
	Total *int8 `json:"total,omitempty" form:"total" binding:"omitempty,oneof=0 1"`
}

// NextPage 下一页，n 是当前页的数据量
func (m *PageQuery) NextPage(n int) bool {
	// 分页
	if m.HasCount() && m.Offset != nil {
		// 最后的一页，不继续了
		if *m.Count > n {
			return false
		}
		*m.Offset += *m.Count
		return true
	}
	// 不分页
	return false
}

// HasCount 是否有分页
func (m *PageQuery) HasCount() bool {
	return m.Count != nil && *m.Count > 0
}

// HasTotal 是否有总数
func (m *PageQuery) HasTotal() bool {
	return m.Total != nil && *m.Total == 1
}

// PageResult 是 Page 的返回值
type PageResult[M any] struct {
	// 总数
	Total int64 `json:"total,omitempty"`
	// 列表
	Data []M `json:"data"`
}

// Page 用于分页查询
func Page[M any](db *gorm.DB, page *PageQuery, res *PageResult[M]) (err error) {
	if page != nil {
		// 总数
		if page.HasTotal() {
			err = db.Count(&res.Total).Error
			if err != nil {
				return err
			}
		}
		// 分页
		if page.Offset != nil {
			db = db.Offset(*page.Offset)
		}
		if page.HasCount() {
			db = db.Limit(*page.Count)
		}
		// 排序
		if page.Order != "" {
			db = db.Order(page.Order)
		}
	}
	// 查询
	err = db.Find(&res.Data).Error
	if err != nil {
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
	err = db.Scan(&ms).Error
	return
}
