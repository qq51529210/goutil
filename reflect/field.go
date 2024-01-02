package reflect

import "reflect"

// var (
// 	// StructFieldNameTagName 是 StructFieldName 结构解析的 tag 名称
// 	StructFieldNameTagName = "field"
// )

// StructFieldName 按顺序提取所有字段名称
func StructFieldName(v any, tag string) []string {
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
	return structFieldName(vv.Type(), tag, make([]string, 0))
}

// structFieldName  是 StructFieldName 的实现
func structFieldName(st reflect.Type, tag string, a []string) []string {
	// 所有字段
	for i := 0; i < st.NumField(); i++ {
		// 类型
		ft := st.Field(i)
		// tag
		name, _, ignore := parseTag(&ft, tag)
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
				a = structFieldName(t, tag, a)
			}
			continue
		}
		// 其他
		a = append(a, name)
	}
	return a
}
