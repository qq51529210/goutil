package reflect

import "reflect"

// var (
// 	// StructToSliceTagName 是 StructToSlice 结构解析的 tag 名称
// 	StructToSliceTagName = "slice"
// )

// StructToSlice 将 v 转换为 slice 必须是结构体
// 主要是用于 redis 的 HSET 命令
// 不可导出的字段，不处理
// 不是基本的类型，不处理，嵌入结构除外
// 嵌入的字段不是结构，或者是 nil 的结构指针，不处理
//
//	type S1 struct {
//	   A string -> ["A", A] 默认使用字段名称
//	   B string `slice:"b"` -> ["b", B] 有 tag 则使用 tag
//	   C string `slice:"omitempty"` -> 忽略零值
//	   D *string `slice:"d,omitempty"` -> 指针不为 nil 不算零值
//	   E string `slice:"-"` -> 忽略
//	}
//
//	type S2 struct {
//	   S1 -> 直接嵌入 ["A", A, "b", B ...]
//	   F S1 -> 不处理
//	}
func StructToSlice(v any, tag string) []any {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	m := make([]any, 0)
	return structToSlice(vv, tag, m)
}

// structToSlice 封装 StructToSlice 的代码
func structToSlice(v reflect.Value, tag string, m []any) []any {
	// 类型
	st := v.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, omitempty, ignore := parseTag(&ft, tag)
		// 忽略
		if ignore {
			continue
		}
		// 值
		fv := v.Field(i)
		// 数据类型
		fk := ft.Type.Kind()
		if fk == reflect.Pointer {
			fv = fv.Elem()
			if fv.IsValid() {
				continue
			}
			fk = fv.Kind()
		} else {
			// 忽略零值
			if fv.IsZero() && omitempty {
				continue
			}
		}
		// 嵌入字段
		if ft.Anonymous {
			// 结构
			if fk == reflect.Struct {
				m = structToSlice(fv, tag, m)
			}
			// 嵌入的不是结构不处理
			continue
		}
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		if (fk >= reflect.Bool && fk <= reflect.Uint64) ||
			(fk >= reflect.Float32 && fk <= reflect.Float64) ||
			fk == reflect.String {
			m = append(m, name)
			m = append(m, fv.Interface())
		}
	}
	return m
}
