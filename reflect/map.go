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
	sv := reflect.ValueOf(v)
	if sv.Kind() == reflect.Pointer {
		sv = sv.Elem()
	}
	if sv.Kind() != reflect.Struct {
		panic("v must be struct")
	}
	m := make(map[string]any)
	return structToMap(sv, tag, m)
}

// structToMap 封装 StructToMap 的代码
func structToMap(sv reflect.Value, tag string, m map[string]any) map[string]any {
	// 类型
	st := sv.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// tag
		name, omitempty, ignore := ParseTag(&ft, tag)
		// 忽略
		if ignore {
			continue
		}
		// 值类型
		fv := sv.Field(i)
		fk := fv.Kind()
		if fk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				// 忽略零值
				if omitempty {
					continue
				}
				// 嵌入的空指针，不处理
				if ft.Anonymous {
					continue
				}
				// 设置
				if name != "" {
					m[name] = nil
				} else {
					m[ft.Name] = nil
				}
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
		// 嵌入字段
		if ft.Anonymous {
			// 结构
			if fk == reflect.Struct {
				structToMap(fv, tag, m)
			}
			// 嵌入的不是结构不处理
			continue
		}
		// 结构
		if fk == reflect.Struct {
			_data := structToMap(fv, tag, make(map[string]any))
			if len(_data) > 0 {
				if name != "" {
					m[name] = _data
				} else {
					m[ft.Name] = _data
				}
			}
			continue
		}
		// 其他
		if name != "" {
			m[name] = fv.Interface()
		} else {
			m[ft.Name] = fv.Interface()
		}
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
//
// 使用 slice，map 做值等等可能会有问题
func StructFromMap(v any, m map[string]any) {
	StructFromMapWithTag(v, m, StructToMapTagName)
}

// StructFromMapWithTag 使用自定义 tag
func StructFromMapWithTag(v any, m map[string]any, tag string) {
	sv := reflect.ValueOf(v)
	if sv.Kind() == reflect.Pointer {
		sv = sv.Elem()
		if sv.Kind() == reflect.Struct {
			structFromMap(sv, tag, m)
			return
		}
	}
	panic("v must be struct pointer")
}

// structFromMap 封装 StructFromMap 的代码
func structFromMap(sv reflect.Value, tag string, data map[string]any) bool {
	hasValue := false
	// 类型
	st := sv.Type()
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// tag
		name, _, ignore := ParseTag(&ft, tag)
		// 忽略
		if ignore {
			continue
		}
		// 字段值类型
		fv := sv.Field(i)
		fk := fv.Kind()
		isNil := false
		// 指针类型，可能需要 new
		if fk == reflect.Pointer {
			if fv.IsNil() {
				isNil = true
				fk = ft.Type.Elem().Kind()
			} else {
				fv = fv.Elem()
				fk = fv.Kind()
			}
		}
		// 嵌入
		if ft.Anonymous {
			if isNil {
				// 空指针，new
				nfv := reflect.New(ft.Type.Elem())
				if structFromMap(nfv.Elem(), tag, data) {
					// 有值才设置
					fv.Set(nfv)
					hasValue = true
				}
			} else {
				// 只处理结构
				if fk == reflect.Struct {
					structFromMap(fv, tag, data)
				}
			}
			continue
		}
		// map 的值
		if name == "" {
			name = ft.Name
		}
		mv, ok := data[name]
		if !ok || mv == nil {
			// 不存在/或者是
			continue
		}
		// 结构
		if fk == reflect.Struct {
			// 那么 mv 是 map[string]any 才能继续
			if _v, ok := mv.(map[string]any); ok {
				if isNil {
					// 空指针
					nfv := reflect.New(ft.Type.Elem())
					if structFromMap(nfv.Elem(), tag, _v) {
						// 有值才设置
						fv.Set(nfv)
						hasValue = true
					}
				} else {
					structFromMap(fv, tag, _v)
				}
			}
			continue
		}
		// mvv 是指针
		mvv := reflect.ValueOf(mv)
		mvk := mvv.Kind()
		if mvk == reflect.Pointer {
			mvv = mvv.Elem()
			mvk = mvv.Kind()
		}
		// 整形
		if (fk == reflect.Int || fk == reflect.Int8 || fk == reflect.Int16 || fk == reflect.Int32 || fk == reflect.Int64) &&
			(mvk == reflect.Int || mvk == reflect.Int8 || mvk == reflect.Int16 || mvk == reflect.Int32 || mvk == reflect.Int64) {
			if isNil {
				fv.Set(reflect.New(ft.Type.Elem()))
			}
			fv.SetInt(mvv.Int())
			continue
		}
		if (fk == reflect.Uint || fk == reflect.Uint8 || fk == reflect.Uint16 || fk == reflect.Uint32 || fk == reflect.Uint64) &&
			(mvk == reflect.Uint || mvk == reflect.Uint8 || mvk == reflect.Uint16 || mvk == reflect.Uint32 || mvk == reflect.Uint64) {
			if isNil {
				fv.Set(reflect.New(ft.Type.Elem()))
			}
			fv.SetUint(mvv.Uint())
			continue
		}
		// 浮点
		if (fk == reflect.Float32 || fk == reflect.Float64) && (mvk == reflect.Float32 || mvk == reflect.Float64) {
			fv.SetFloat(mvv.Float())
			continue
		}
		// 其他的类型相同才赋值
		if isNil {
			ftt := ft.Type.Elem()
			if ftt == mvv.Type() {
				nfv := reflect.New(ftt).Elem()
				nfv.Set(mvv)
				sv.Field(i).Set(nfv.Addr())
			}
		} else {
			if ft.Type == mvv.Type() {
				fv.Set(mvv)
			}
		}
		//
		hasValue = true
	}
	return hasValue
}
