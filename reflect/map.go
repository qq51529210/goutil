package reflect

import "reflect"

var (
	// StructToMapTagName 是 StructToMap 结构解析的 tag 名称
	StructToMapTagName = "map"
)

// StructToMap 将 v 转换为 map，v 必须是结构体
// 嵌入的字段不是结构，或者是 nil 的结构指针，不处理
// 不可导出的字段，不处理
//
//	type S1 struct {
//	   A string -> map["A"]=A 默认使用字段名称
//	   B string `map:"b"` -> map["b"]=B 有 tag 则使用 tag
//	   C string `map:"omitempty"` -> 忽略零值
//	   D *string `map:"d,omitempty"` -> 指针不为 nil 不算零值
//	   E string `map:"-"` -> 忽略
//	}
//
//	type S2 struct {
//	   S1 -> 直接嵌入 map["A"]=A , map["b"]=B ...
//	   F S1 -> map["F"]=F
//	}
func StructToMap(v any) map[string]any {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	m := make(map[string]any)
	return structToMap(vv, m)
}

// structToMap 封装 StructToMap 的代码
func structToMap(v reflect.Value, m map[string]any) map[string]any {
	// 类型
	st := v.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, omitempty, ignore := parseTag(&ft, StructToMapTagName)
		// 忽略
		if ignore {
			continue
		}
		// 值
		fv := v.Field(i)
		// 数据类型
		fk := ft.Type.Kind()
		if fk == reflect.Pointer {
			fk = ft.Type.Elem().Kind()
			// nil 指针
			if fv.IsNil() {
				// 忽略
				if omitempty {
					continue
				}
				// 嵌入不处理 / 不可导出
				if ft.Anonymous || !ft.IsExported() {
					continue
				}
				// 不是嵌入
				m[ft.Name] = nil
				continue
			}
			fv = fv.Elem()
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
				structToMap(fv, m)
			}
			// 嵌入的不是结构不处理
			continue
		}
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// 结构
		if fk == reflect.Struct {
			m[name] = structToMap(fv, make(map[string]any))
			continue
		}
		// 其他
		m[name] = fv.Interface()
	}
	return m
}
