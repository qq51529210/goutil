package util

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// GORMTime 创建和更新时间
type GORMTime struct {
	// 数据库的创建时间，时间戳
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 数据库的更新时间，时间戳
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// GORMPage 分页查询参数
type GORMPage struct {
	// 偏移，小于 0 不匹配
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `form:"order"`
}

// GORMQuery 是 All 函数格式化查询参数的接口
type GORMQuery interface {
	Init(*gorm.DB) *gorm.DB
}

// GORMList 是 GORMList 的返回值
type GORMList[M any] struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Data []M `json:"data"`
}

// GORMDB 模板 api
type GORMDB[K, M any] struct {
	db *gorm.DB
	m  M
}

// NewGORMDB 返回新的 GORMDB
func NewGORMDB[K, M any](db *gorm.DB, m M) *GORMDB[K, M] {
	return &GORMDB[K, M]{
		db: db,
		m:  m,
	}
}

// Init 初始化
func (g *GORMDB[K, M]) Init(db *gorm.DB, m M) {
	g.db = db
	g.m = m
}

// Model 返回
func (g *GORMDB[K, M]) Model() *gorm.DB {
	return g.db.Model(g.m)
}

// All 返回列表查询结果
func (g *GORMDB[K, M]) All(query GORMQuery) ([]M, error) {
	db := g.db.Model(g.m)
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 列表
	var models []M
	err := db.Find(&models).Error
	if err != nil {
		return nil, err
	}
	// 返回
	return models, nil
}

// List 返回列表查询结果
func (g *GORMDB[K, M]) List(page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	return gormList(g.db, page, query, res)
}

// Save 添加
func (g *GORMDB[K, M]) Save(m M) (int64, error) {
	db := g.db.Save(m)
	return db.RowsAffected, db.Error
}

// Add 添加
func (g *GORMDB[K, M]) Add(m M) (int64, error) {
	db := g.db.Create(m)
	return db.RowsAffected, db.Error
}

// Update 根据主键更新
func (g *GORMDB[K, M]) Update(m M) (int64, error) {
	db := g.db.Updates(m)
	return db.RowsAffected, db.Error
}

// Delete 根据主键删除
func (g *GORMDB[K, M]) Delete(k K) (int64, error) {
	db := g.db.Delete(g.m, k)
	return db.RowsAffected, db.Error
}

// BatchDelete 根据主键批量删除
func (g *GORMDB[K, M]) BatchDelete(ks []K) (int64, error) {
	db := g.db.Delete(g.m, ks)
	return db.RowsAffected, db.Error
}

// Get 根据主键查询
func (g *GORMDB[K, M]) Get(m M) (bool, error) {
	err := g.db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Select 根据主键查询，可选择列
func (g *GORMDB[K, M]) Select(m M, c ...string) (bool, error) {
	db := g.db
	if len(c) > 0 {
		db = db.Select(c)
	}
	err := db.First(m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// In 根据主键查询，where in ks
func (g *GORMDB[K, M]) In(ks []K) ([]M, error) {
	var ms []M
	err := g.db.Find(&ms, ks).Error
	if err != nil {
		return nil, err
	}
	return ms, nil
}

var (
	// GORMInitQueryTag 是 GORMInitQuery 解析 tag 的名称
	GORMInitQueryTag = "gq"
)

// Query 将 v 格式化到 where ，全部是 AND ，略过空值
//
//	type query struct {
//	  A *int64 `gq:"eq"` db.Where("`A` = ?", A)
//	  B string `gq:"like"` db.Where("`B` LiKE %%s%%", B)
//	  C *int64 `gq:"gt=A"` db.Where("`A` < ?", C)
//	  D *int64 `gq:"gte=A"` db.Where("`A` <= ?", D)
//	  E *int64 `gq:"lt=A"` db.Where("`A` > ?", E)
//	  F *int64 `gq:"let=A"` db.Where("`A` >= ?", F)
//	  G *int64 `gq:"neq"` db.Where("`G` != ?", G)
//	}
//
// 先这样，以后遇到再加
func GORMInitQuery(db *gorm.DB, q any) *gorm.DB {
	v := reflect.ValueOf(q)
	vk := v.Kind()
	if vk == reflect.Pointer {
		v = v.Elem()
		vk = v.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	return gormInitQuery(db, v)
}

func gormInitQuery(db *gorm.DB, v reflect.Value) *gorm.DB {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构
		if fvk == reflect.Struct {
			gormInitQuery(db, fv)
			continue
		}
		if fvk == reflect.String {
			// 空值
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(GORMInitQueryTag)
		p := strings.TrimPrefix(tn, "eq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` = ?", p), fv.Interface())
			continue
		}
		if tn == "eq" {
			db = db.Where(fmt.Sprintf("`%s` = ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "neq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` != ?", p), fv.Interface())
			continue
		}
		if tn == "neq" {
			db = db.Where(fmt.Sprintf("`%s` != ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "like=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", p), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		if tn == "like" {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", ft.Name), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		p = strings.TrimPrefix(tn, "gt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` < ?", p), fv.Interface())
			continue
		}
		if tn == "gt" {
			db = db.Where(fmt.Sprintf("`%s` < ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "gte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` <= ?", p), fv.Interface())
			continue
		}
		if tn == "gte" {
			db = db.Where(fmt.Sprintf("`%s` <= ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` > ?", p), fv.Interface())
			continue
		}
		if tn == "lt" {
			db = db.Where(fmt.Sprintf("`%s` > ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` >= ?", p), fv.Interface())
			continue
		}
		if tn == "lte" {
			db = db.Where(fmt.Sprintf("`%s` >= ?", ft.Name), fv.Interface())
			continue
		}
	}
	//
	return db
}

func gormList[M any](db *gorm.DB, page *GORMPage, query GORMQuery, res *GORMList[M]) error {
	// 条件
	if query != nil {
		db = query.Init(db)
	}
	// 总数
	err := db.Count(&res.Total).Error
	if err != nil {
		return err
	}
	if page != nil {
		// 分页
		if page.Offset != nil {
			db = db.Offset(*page.Offset)
		}
		if page.Count != nil {
			db = db.Limit(*page.Count)
		}
		// 排序
		if page.Order != "" {
			db = db.Order(page.Order)
		}
	}
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}
