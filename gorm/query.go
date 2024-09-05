package gorm

import (
	"fmt"
	gr "goutil/reflect"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	// InitQueryTag 是 InitQuery 解析 tag 的名称
	InitQueryTag = "gq"
	// InitQueryFunc 是 InitQuery 处理函数
	InitQueryFunc = map[string]func(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB{
		"in":     QueryIN,
		"nin":    QueryNIN,
		"eq":     QueryEQ,
		"neq":    QueryNEQ,
		"like":   QueryLIKE,
		"llike":  QueryLLIKE,
		"rlike":  QueryRLIKE,
		"gt":     QueryGT,
		"gte":    QueryGTE,
		"lt":     QueryLT,
		"lte":    QueryLTE,
		"null":   QueryNULL,
		"select": QuerySelect,
		"omit":   QueryOmit,
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
//	  F6 *int64 `gq:"lte=F"` db.Where("`F`>=?", F6)
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
		// 是零值不处理
		if fv.IsZero() {
			continue
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

// QueryIN field in ?
func QueryIN(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.Slice || kind == reflect.Array {
		return db.Where(fmt.Sprintf("`%s` IN ?", field), value.Interface())
	}
	return db
}

// QueryNIN field not in ?
func QueryNIN(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.Slice || kind == reflect.Array {
		return db.Where(fmt.Sprintf("`%s` NOT IN ?", field), value.Interface())
	}
	return db
}

// QueryEQ field = ?
func QueryEQ(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` = ?", field), value.Interface())
}

// QueryNEQ field != ?
func QueryNEQ(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` != ?", field), value.Interface())
}

// QueryLIKE field like %?%
func QueryLIKE(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.String {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", field), fmt.Sprintf("%%%s%%", value.String()))
	}
	return db
}

// QueryLLIKE field like %?
func QueryLLIKE(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.String {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", field), fmt.Sprintf("%%%s", value.String()))
	}
	return db
}

// QueryRLIKE field like ?%
func QueryRLIKE(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.String {
		return db.Where(fmt.Sprintf("`%s` LIKE ?", field), fmt.Sprintf("%s%%", value.String()))
	}
	return db
}

// QueryGT field < ?
func QueryGT(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` < ?", field), value.Interface())
}

// QueryGTE field <= ?
func QueryGTE(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` <= ?", field), value.Interface())
}

// QueryLT field > ?
func QueryLT(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` > ?", field), value.Interface())
}

// QueryLTE field >= ?
func QueryLTE(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	return db.Where(fmt.Sprintf("`%s` >= ?", field), value.Interface())
}

// QueryNULL field IS [NOT]NULL
func QueryNULL(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	ok := false
	if kind >= reflect.Int && kind <= reflect.Int64 {
		ok = value.Int() == 1
	} else if kind >= reflect.Uint && kind <= reflect.Uint64 {
		ok = value.Uint() == 1
	} else if kind == reflect.Bool {
		ok = value.Interface().(bool)
	} else if kind == reflect.String {
		ok = value.String() == "1"
	} else {
		return db
	}
	if ok {
		return db.Where(fmt.Sprintf("`%s` IS NOT NULL", field))
	}
	return db.Where(fmt.Sprintf("`%s` IS NULL", field))
}

// QuerySelect select fields
func QuerySelect(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.Array {
		if vs, ok := value.Interface().([]string); ok {
			return db.Select(vs)
		}
	}
	return db
}

// QueryOmit omit fields
func QueryOmit(db *gorm.DB, field string, value reflect.Value, kind reflect.Kind) *gorm.DB {
	if kind == reflect.Array {
		if vs, ok := value.Interface().([]string); ok {
			return db.Omit(vs...)
		}
	}
	return db
}

// InitQuery 将 v 格式化到 where ，全部是 AND ，略过空值
// 其他条件，自己添加 InitQueryFunc
//
//	type query struct {
//	  F1 *int64 `json:"f1,omitempty"`
//	  F2 string `json:"ff"`
//	}
//
//	gorm.Expr("JSON_SET(`f1`, '$.f1', ?, '$.ff', ?)", F1, F2)
//
// 先这样，以后遇到再加
func InitMysqlJSONSet(sv any, field string) clause.Expr {
	return InitMysqlJSONSetWithTag(sv, field, "json")
}

func InitMysqlJSONSetWithTag(sv any, field string, tag string) clause.Expr {
	v := reflect.ValueOf(sv)
	k := v.Kind()
	if k == reflect.Pointer {
		// 空指针
		if v.IsNil() {
			panic("nil pointer")
		}
		v = v.Elem()
		k = v.Kind()
	}
	if k != reflect.Struct {
		panic("input must be struct")
	}
	// 检索
	data := make(map[string]any)
	initMysqlJSON(v, "", tag, data)
	if len(data) < 1 {
		panic("invalid input struct")
	}
	//
	var str strings.Builder
	fmt.Fprintf(&str, "JSON_SET(`%s`", field)
	var vs []any
	for k, v := range data {
		str.WriteString(",'$.")
		str.WriteString(k)
		str.WriteString("',?")
		vs = append(vs, v)
	}
	str.WriteString(")")
	return gorm.Expr(str.String(), vs...)
}

// sv 解析的结构
// pname 所属的结构名称，sv 是字段
// tag 就是 tag
// data 用于装结果
func initMysqlJSON(sv reflect.Value, pname string, tag string, data map[string]any) {
	st := sv.Type()
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// tag
		name, omitempty, ignore := gr.ParseTag(&ft, tag)
		// 忽略字段
		if ignore {
			continue
		}
		fv := sv.Field(i)
		fk := fv.Kind()
		if fk == reflect.Pointer {
			// 空指针就算了
			if fv.IsNil() {
				continue
			}
			// 有指针，就不算零值
			fv = fv.Elem()
			fk = fv.Kind()
		} else {
			// 忽略零值
			if fv.IsZero() && omitempty {
				continue
			}
		}
		// 嵌入
		if ft.Anonymous {
			// 嵌入只处理结构
			if fk == reflect.Struct {
				// 虽然是嵌入，但是 tag 有名称
				if name != "" {
					if pname != "" {
						initMysqlJSON(fv, pname+"."+name, tag, data)
					} else {
						initMysqlJSON(fv, name, tag, data)
					}
				} else {
					initMysqlJSON(fv, pname, tag, data)
				}
			}
			continue
		}
		// 没有 tag ，使用字段名称
		if name == "" {
			name = ft.Name
		}
		// 结构
		if fk == reflect.Struct {
			if pname != "" {
				initMysqlJSON(fv, pname+"."+name, tag, data)
			} else {
				initMysqlJSON(fv, name, tag, data)
			}
			continue
		}
		// 其他
		if pname != "" {
			data[pname+"."+name] = fv.Interface()
		} else {
			data[name] = fv.Interface()
		}
	}
}
