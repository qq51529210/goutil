package gorm

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

var (
	// InitQueryTag 是 InitQuery 解析 tag 的名称
	InitQueryTag = "gq"
	// InitQueryFunc 是 InitQuery 处理函数
	InitQueryFunc = map[string]func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB{
		"in": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			if kind == reflect.Slice || kind == reflect.Array {
				if !value.IsZero() {
					return db.Where(fmt.Sprintf("`%s` IN ?", field), value.Interface())
				}
			}
			return db
		},
		"nin": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			if kind == reflect.Slice || kind == reflect.Array {
				if !value.IsZero() {
					return db.Where(fmt.Sprintf("`%s` NOT IN ?", field), value.Interface())
				}
			}
			return db
		},
		"eq": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`=?", field), value.Interface())
		},
		"neq": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`!=?", field), value.Interface())
		},
		"like": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			if kind == reflect.String {
				s := value.String()
				if s != "" {
					return db.Where(fmt.Sprintf("`%s` LIKE ?", field), fmt.Sprintf("%%%s%%", s))
				}
			}
			return db
		},
		"gt": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`<?", field), value.Interface())
		},
		"gte": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`<=?", field), value.Interface())
		},
		"lt": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`>?", field), value.Interface())
		},
		"lte": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			return db.Where(fmt.Sprintf("`%s`>=?", field), value.Interface())
		},
		"null": func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
			ok := false
			if kind >= reflect.Int && kind <= reflect.Uint64 {
				ok = value.Interface() == 1
			} else if kind == reflect.Bool {
				ok = value.Interface().(bool)
			} else if kind == reflect.String {
				ok = value.String() == "true"
			} else {
				return db
			}
			if ok {
				return db.Where(fmt.Sprintf("`%s` IS NULL", field))
			}
			return db.Where(fmt.Sprintf("`%s` IS NOT NULL", field))
		},
	}
)

// InitQuery 将 v 格式化到 where ，全部是 AND ，略过空值
// 其他条件，自己添加 InitQueryFunc
//
//	type query struct {
//	  F1 *int64 `gq:"eq"` db.Where("`F1`=?", F1)
//	  F2 string `gq:"like"` db.Where("`F2` LiKE %%s%%", F2)
//	  F3 *int64 `gq:"gt=F"` db.Where("`F`<?", F3)
//	  F4 *int64 `gq:"gte=F"` db.Where("`F`<=?", F4)
//	  F5 *int64 `gq:"lt=F"` db.Where("`F`>?", F5)
//	  F6 *int64 `gq:"let=F"` db.Where("`F`>=?", F6)
//	  F7 *int64 `gq:"neq"` db.Where("`F`!=?", F7)
//	  F8 *int8 `gq:"null"` if F8==0/1 true/false db.Where("`F8` IS NULL/IS NOT NULL")
//	  F9 []int64 `gq:"in=F"` db.Where("`F` IN ?", F9)
//	  F10 []int64 `gq:"nin=F"` db.Where("`F` NOT IN ?", F10)
//	}
//
// 先这样，以后遇到再加
func InitQuery(db *gorm.DB, q any) *gorm.DB {
	return InitQueryWithTag(db, q, InitQueryTag)
}

// InitQueryWithTag 自定义 tag
func InitQueryWithTag(db *gorm.DB, q any, tag string) *gorm.DB {
	v := reflect.ValueOf(q)
	vk := v.Kind()
	if vk == reflect.Pointer {
		if v.IsNil() {
			return db
		}
		v = v.Elem()
		vk = v.Kind()
	}
	if vk != reflect.Struct {
		panic("q must be struct")
	}
	return initQuery(db, v, tag)
}

// initQuery 是 InitQuery 的实现
func initQuery(db *gorm.DB, v reflect.Value, tag string) *gorm.DB {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		// 类型
		ft := vt.Field(i)
		// 值
		fv := v.Field(i)
		// 数据类型
		fk := fv.Kind()
		if fk == reflect.Pointer {
			fv = fv.Elem()
			// 无效值
			if !fv.IsValid() {
				continue
			}
			fk = fv.Kind()
		}
		// 嵌入不是结构不处理
		if ft.Anonymous {
			if fk == reflect.Struct {
				db = initQuery(db, fv, tag)
			}
			continue
		}
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// 没有 tag 不处理
		tag := ft.Tag.Get(tag)
		if tag == "" {
			continue
		}
		// eq=F
		var name string
		j := strings.Index(tag, "=")
		if j < 0 {
			name = ft.Name
		} else {
			name = tag[j+1:]
			tag = tag[:j]
		}
		// 处理
		fun := InitQueryFunc[tag]
		if fun == nil {
			continue
		}
		db = fun(db, name, fv, fk)
	}
	//
	return db
}
