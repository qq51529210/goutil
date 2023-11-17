package reflect

import "reflect"

var (
	// StructToSliceTagName 是 StructToSlice 结构解析的 tag 名称
	StructToSliceTagName = "slice"
)

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
func StructToSlice(v any) []any {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	m := make([]any, 0)
	return structToSlice(vv, m)
}

// structToSlice 封装 StructToSlice 的代码
func structToSlice(v reflect.Value, m []any) []any {
	// 类型
	st := v.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, omitempty, ignore := parseTag(&ft, StructToSliceTagName)
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
				m = structToSlice(fv, m)
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

var (
	// StructFieldNameSliceTagName 是 StructFieldNameSlice 结构解析的 tag 名称
	StructFieldNameSliceTagName = "field"
)

// StructFieldNameSlice 按顺序提取所有字段名称
func StructFieldNameSlice(v any) []string {
	vv := reflect.ValueOf(v)
	if !vv.IsValid() {
		panic("vv is invalid")
	}
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	return structFieldNameSlice(vv.Type(), make([]string, 0))
}

// structFieldNameSlice  是 StructFieldNameSlice 的实现
func structFieldNameSlice(st reflect.Type, a []string) []string {
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, _, ignore := parseTag(&ft, StructFieldNameSliceTagName)
		// 忽略 / 不可导出
		if ignore || !ft.IsExported() {
			continue
		}
		// 嵌入
		if ft.Anonymous {
			t := ft.Type
			if ft.Type.Kind() == reflect.Pointer {
				t = t.Elem()
			}
			// 结构
			if t.Kind() == reflect.Struct {
				a = structFieldNameSlice(t, a)
			}
			continue
		}
		// 其他
		a = append(a, name)
	}
	return a
}
