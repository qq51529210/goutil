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
	return StructToMapWithTag(v, StructToMapTagName)
}

// StructToMapWithTag 使用自定义 tag
func StructToMapWithTag(v any, tag string) map[string]any {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
	}
	if vv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	m := make(map[string]any)
	return structToMap(vv, tag, m)
}

// structToMap 封装 StructToMap 的代码
func structToMap(v reflect.Value, tag string, m map[string]any) map[string]any {
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
				structToMap(fv, tag, m)
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
			m[name] = structToMap(fv, tag, make(map[string]any))
			continue
		}
		// 其他
		m[name] = fv.Interface()
	}
	return m
}

// StructFromMap 使用 m 填充 v ，v 必须是结构体指针
//
//	type S1 struct {
//	   A string -> A=map["A"] 默认使用字段名称
//	   B string `map:"b"` -> B=map["b"] 有 tag 则使用 tag
//	   D *string -> D=map["D"]  指针 new
//	   E string `map:"-"` -> 忽略
//	}
func StructFromMap(v any, m map[string]any) {
	StructFromMapWithTag(v, m, StructToMapTagName)
}

// StructFromMapWithTag 使用自定义 tag
func StructFromMapWithTag(v any, m map[string]any, tag string) {
	vv := reflect.ValueOf(v)
	if vv.Kind() == reflect.Pointer {
		vv = vv.Elem()
		if vv.Kind() == reflect.Struct {
			structFromMap(vv, tag, m)
			return
		}
	}
	panic("v must be struct pointer")
}

var (
	structFromMapType = reflect.TypeOf(make(map[string]any))
)

// structFromMap 封装 StructFromMap 的代码
func structFromMap(v reflect.Value, tag string, m map[string]any) {
	// 类型
	st := v.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, _, ignore := parseTag(&ft, tag)
		// 忽略
		if ignore {
			continue
		}
		// 值
		fv := v.Field(i)
		// 数据类型是否一致
		fk := ft.Type.Kind()
		if fk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				fv.Set(reflect.New(ft.Type.Elem()))
			}
			fv = fv.Elem()
			fk = fv.Kind()
		}
		// 嵌入
		if ft.Anonymous {
			// 必须是结构
			if fk == reflect.Struct {
				structFromMap(fv, tag, m)
			}
			continue
		}
		// map 值
		mv, ok := m[name]
		if !ok {
			continue
		}
		mvv := reflect.ValueOf(mv)
		mvk := mvv.Kind()
		if mvk == reflect.Pointer {
			mvv = mvv.Elem()
			mvk = mvv.Kind()
		}
		// 结构对应 map
		if fk == reflect.Struct {
			if mvv.Type() == structFromMapType {
				structFromMap(fv, tag, mv.(map[string]any))
			}
			continue
		}
		// 数据类型是否一致
		if fk != mvk {
			continue
		}
		fv.Set(mvv)
	}
}
