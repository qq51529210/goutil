package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// IsNilOrEmpty 如果 v 是空指针或者零值，返回 true
// 指针有零值不算零值
func IsNilOrEmpty(v any) bool {
	return isNilOrEmpty(reflect.ValueOf(v))
}

// isNilOrEmpty 是 IsNilOrEmpty 的实现
func isNilOrEmpty(v reflect.Value) bool {
	// 无效
	if !v.IsValid() || v.IsZero() {
		return true
	}
	k := v.Kind()
	// 指针
	for k == reflect.Pointer {
		v = v.Elem()
		k = v.Kind()
	}
	// 类型
	switch k {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		// 没有数据算空
		return v.Len() == 0
	case reflect.Struct:
		// 所有的字段为空才算空
		for i := 0; i < v.NumField(); i++ {
			if !isNilOrEmpty(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		// 其他
		return !v.IsValid() || v.IsZero()
	}
}

var (
	// StructToMapTagName 是 StructToMap 结构解析的 tag 名称
	StructToMapTagName = "map"
)

// StructToMap 将 v 转换为 map，v 必须是结构体
// 嵌入的字段必须是结构，否则不处理
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
	m := make(map[string]any)
	return structToMap(reflect.ValueOf(v), m)
}

// structToMap 封装 StructToMap 的代码
func structToMap(v reflect.Value, m map[string]any) map[string]any {
	// 无效
	if !v.IsValid() {
		return m
	}
	// 指针
	vk := v.Kind()
	if vk == reflect.Pointer {
		v = v.Elem()
		vk = v.Kind()
	}
	// 必须是结构
	if vk != reflect.Struct {
		panic("v type must be struct")
	}
	// 类型
	t := v.Type()
	// 所有字段
	for i := 0; i < t.NumField(); i++ {
		// 类型
		ft := t.Field(i)
		// 不可导出，但是嵌入的可能可以导出
		if !ft.IsExported() && !ft.Anonymous {
			continue
		}
		// 值
		fv := v.Field(i)
		// 无效
		if !fv.IsValid() {
			continue
		}
		// tag
		name, omitempty, ignore := structToMapTag(&ft)
		// 忽略字段 / 零值
		if ignore || (omitempty && fv.IsZero()) {
			continue
		}
		// 指针
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构
		if fvk == reflect.Struct {
			// 嵌入的
			if ft.Anonymous {
				m = structToMap(fv, m)
			} else {
				m[name] = structToMap(fv, make(map[string]any))
			}
			continue
		}
		// 其他值
		if fv.IsValid() {
			m[name] = fv.Interface()
		} else {
			m[name] = nil
		}
	}
	return m
}

// structToMapTag 是 structToMap 解析 tag 的封装
func structToMapTag(f *reflect.StructField) (name string, omitempty, ignore bool) {
	name = f.Name
	tag := f.Tag.Get(StructToMapTagName)
	for tag != "" {
		var s string
		i := strings.IndexByte(tag, ',')
		if i < 0 {
			s = tag
			tag = ""
		} else {
			s = tag[:i]
			tag = tag[i+1:]
		}
		switch s {
		case "omitempty":
			omitempty = true
		case "-":
			ignore = true
		default:
			if s != "" {
				name = s
			}
		}
	}
	return
}

// StructFieldValue 给字段赋值，v 必须是结构指针
func StructFieldValue(v any) {
	vv := reflect.ValueOf(v)
	if vv.Kind() != reflect.Pointer {
		panic("v must be pointer")
	}
	vv = vv.Elem()
	if vv.Kind() != reflect.Struct {
		panic("v must be struct pointer")
	}
	//
	st := vv.Type()
	for i := 0; i < st.NumField(); i++ {
		ft := st.Field(i)
		tag := ft.Tag.Get("value")
		if tag != "" {
			structFieldtDefaultValue(vv.Field(i), tag)
			continue
		}
		tag = ft.Tag.Get("field")
		if tag != "" {
			structFieldDefaultField(vv, vv.Field(i), tag)
			continue
		}
		tag = ft.Tag.Get("func")
		if tag != "" {
			structFieldDefaultFunc(vv, vv.Field(i), tag)
			continue
		}
	}
}

// structFieldtDefaultValue 设置字段的为指定的值
func structFieldtDefaultValue(fieldValue reflect.Value, defaultValue string) {
	kind := fieldValue.Kind()
	// 指针
	if kind == reflect.Pointer {
		// nil new
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		fieldValue = fieldValue.Elem()
		kind = fieldValue.Kind()
	}
	// 字符串->数字
	if kind >= reflect.Uint && kind <= reflect.Uint64 {
		n, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			panic(err)
		}
		fieldValue.SetUint(uint64(n))
		return
	}
	if kind >= reflect.Int && kind <= reflect.Int64 {
		n, err := strconv.ParseInt(defaultValue, 10, 64)
		if err != nil {
			panic(err)
		}
		fieldValue.SetInt(n)
		return
	}
	if kind >= reflect.Float32 && kind <= reflect.Float64 {
		n, err := strconv.ParseFloat(defaultValue, 64)
		if err != nil {
			panic(err)
		}
		fieldValue.SetFloat(n)
		return
	}
	// 字符串
	if kind == reflect.String {
		fieldValue.SetString(defaultValue)
		return
	}
	//
	if kind == reflect.Bool {
		fieldValue.SetBool(defaultValue == "true")
		return
	}
	panic(fmt.Sprintf("unsupported field type %v", kind))
}

// structFieldDefaultField 设置字段的为指定的字段的值
func structFieldDefaultField(structValue, fieldValue reflect.Value, fieldName string) {
	// 找到标记的字段
	srcFieldValue := structValue.FieldByName(fieldName)
	// 无效就算了
	if !srcFieldValue.IsValid() {
		return
	}
	srcKind := srcFieldValue.Kind()
	if srcKind == reflect.Pointer {
		srcFieldValue = srcFieldValue.Elem()
		srcKind = srcFieldValue.Kind()
	}
	//
	dstKind := fieldValue.Kind()
	if dstKind == reflect.Pointer {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		fieldValue = fieldValue.Elem()
		dstKind = fieldValue.Kind()
	}
	// 不相同
	if dstKind != srcKind {
		panic(fmt.Sprintf("field %s different type %v", fieldName, srcKind))
	}
	// 相同
	if (dstKind >= reflect.Bool && dstKind <= reflect.Uint64) ||
		(dstKind >= reflect.Float32 && dstKind <= reflect.Float64) ||
		dstKind == reflect.String {
		fieldValue.Set(srcFieldValue)
		return
	}
	panic(fmt.Sprintf("unsupported field type %v", dstKind))
}

// structFieldDefaultFunc 设置字段的为指定的字段的值
func structFieldDefaultFunc(structValue, fieldValue reflect.Value, funcName string) {
	kind := fieldValue.Kind()
	// nil 指针
	if kind == reflect.Pointer {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		fieldValue = fieldValue.Elem()
	}
	fieldValue.Set(structValue.MethodByName(funcName))
}

var (
	// StructDiffFieldTagName 是 StructDiffField 结构解析的 tag 名称
	StructDiffFieldTagName = "diff"
)

// StructDiffField 从 src 中找出 与 dst 相同类型（值可以与指针对比）
// 和名称但不同值的字段，然后返回这些字段的 map
//
//	type src struct {
//	   A string 默认与 dst.A 对比
//	   B string `diff:"BB"` -> 与 dst.BB 对比
//	   C string `diff:"omitempty"` -> 忽略零值
//	   D *string `diff:"omitempty"` -> 指针不为 nil 不算零值
//	   E string `diff:"-"` -> 忽略
//	}
func StructDiffField(dst, src any, ignore bool) map[string]any {
	// dst
	dv := reflect.ValueOf(dst)
	dk := dv.Kind()
	if dk == reflect.Pointer {
		dv = dv.Elem()
		dk = dv.Kind()
	}
	if dk != reflect.Struct {
		panic("dst must be struct type")
	}
	// src
	sv := reflect.ValueOf(src)
	sk := sv.Kind()
	if sk == reflect.Pointer {
		sv = sv.Elem()
		sk = sv.Kind()
	}
	if sk != reflect.Struct {
		panic("src must be struct")
	}
	//
	return structDiffField(dv, sv, make(map[string]any))
}

// structDiffField 是 StructDiffFields 的实现
func structDiffField(dst, src reflect.Value, m map[string]any) map[string]any {
	st := src.Type()
	for i := 0; i < st.NumField(); i++ {
		// src 值
		sfv := src.Field(i)
		// src 值无效
		if !sfv.IsValid() {
			continue
		}
		// src 字段
		sft := st.Field(i)
		// 不可导出的
		if !sft.IsExported() {
			continue
		}
		// src 嵌入字段
		if sft.Anonymous {
			sfk := sfv.Kind()
			if sfk == reflect.Pointer {
				sfv = sfv.Elem()
			}
			m = structDiffField(dst, sfv, m)
			continue
		}
		// src tag
		name, omitempty, ignore := structDiffFieldTag(&sft)
		// 忽略字段 / 零值
		if ignore || (omitempty && sfv.IsZero()) {
			continue
		}
		// dst 值
		dfv := dst.FieldByName(name)
		// sfk 类型
		sfk := sfv.Kind()
		if sfk == reflect.Pointer {
			sfv = sfv.Elem()
			if !sfv.IsValid() {
				m[sft.Name] = nil
				continue
			}
			sfk = sfv.Kind()
		}
		// dst 类型
		dfk := dfv.Kind()
		if dfk == reflect.Pointer {
			dfv = dfv.Elem()
			dfk = dfv.Kind()
		}
		// dst 值无效，一般是没有这个字段
		if !dfv.IsValid() || sfk != dfk {
			m[sft.Name] = sfv.Interface()
			continue
		}
		// 比较值
		sd, dd := sfv.Interface(), dfv.Interface()
		if !reflect.DeepEqual(sd, dd) {
			m[sft.Name] = sd
		}
	}
	return m
}

// structDiffFieldTag 是 structDiffField 解析 tag 的封装
func structDiffFieldTag(f *reflect.StructField) (name string, omitempty, ignore bool) {
	name = f.Name
	tag := f.Tag.Get(StructDiffFieldTagName)
	for tag != "" {
		var s string
		i := strings.IndexByte(tag, ',')
		if i < 0 {
			s = tag
			tag = ""
		} else {
			s = tag[:i]
			tag = tag[i+1:]
		}
		switch s {
		case "omitempty":
			omitempty = true
		case "-":
			ignore = true
		default:
			if s != "" {
				name = s
			}
		}
	}
	return
}

// // CopyStruct 拷贝 src 和 dst 中的相同名称和类型的字段，
// // 如果 dst 的字段不是零值则不拷贝。
// func CopyStruct(dst, src any) {
// 	copyStruct(copyStructCheck(dst, src))
// }

// // copyStruct 封装 CopyStruct 代码
// func copyStruct(dst, src reflect.Value) {
// 	// 结构类型
// 	srcType := src.Type()
// 	for i := 0; i < srcType.NumField(); i++ {
// 		srcField := src.Field(i)
// 		// src 零值
// 		if !srcField.IsValid() || srcField.IsZero() {
// 			continue
// 		}
// 		srcTypeField := srcType.Field(i)
// 		// dst 同名字段
// 		dstField := dst.FieldByName(srcTypeField.Name)
// 		// dst 不为零
// 		if !dstField.IsValid() || !dstField.CanSet() || !dstField.IsZero() {
// 			continue
// 		}
// 		dstFieldType := dstField.Type()
// 		// 不同类型
// 		if srcTypeField.Type != dstFieldType {
// 			srcFieldKind := srcField.Kind()
// 			// 看看是不是结构体
// 			if srcFieldKind == reflect.Pointer {
// 				srcField = srcField.Elem()
// 				srcFieldKind = srcField.Kind()
// 			}
// 			dstFieldKind := dstField.Kind()
// 			if dstFieldKind == reflect.Pointer {
// 				dstField = dstField.Elem()
// 				dstFieldKind = dstField.Kind()
// 			}
// 			// 都是结构体，进去赋值
// 			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
// 				copyStruct(dstField, srcField)
// 			}
// 			// 不是就算了
// 			continue
// 		}
// 		// 相同类型，赋值
// 		dstField.Set(srcField)
// 	}
// }

// // CopyStructAll 拷贝 src 和 dst 中的相同名称和类型的字段
// func CopyStructAll(dst, src any) {
// 	copyStructAll(copyStructCheck(dst, src))
// }

// // copyStructAll 封装 CopyStructAll 代码
// func copyStructAll(dst, src reflect.Value) {
// 	// type
// 	srcType := src.Type()
// 	for i := 0; i < srcType.NumField(); i++ {
// 		srcField := src.Field(i)
// 		if !srcField.IsValid() {
// 			continue
// 		}
// 		srcTypeField := srcType.Field(i)
// 		dstField := dst.FieldByName(srcTypeField.Name)
// 		if !dstField.IsValid() || !dstField.CanSet() {
// 			continue
// 		}
// 		dstFieldType := dstField.Type()
// 		// 不同类型
// 		if srcTypeField.Type != dstFieldType {
// 			srcFieldKind := srcField.Kind()
// 			// 看看是不是结构体
// 			if srcFieldKind == reflect.Pointer {
// 				srcField = srcField.Elem()
// 				srcFieldKind = srcField.Kind()
// 			}
// 			dstFieldKind := dstField.Kind()
// 			if dstFieldKind == reflect.Pointer {
// 				dstField = dstField.Elem()
// 				dstFieldKind = dstField.Kind()
// 			}
// 			// 都是结构体，进去赋值
// 			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
// 				copyStructAll(dstField, srcField)
// 			}
// 			continue
// 		}
// 		// 相同类型，赋值
// 		dstField.Set(srcField)
// 	}
// }

// // CopyStructNotEmpty 拷贝 src 和 dst 中的相同名称和类型的字段
// // src 为零值不拷贝，dst 不为零值，也拷贝哦
// func CopyStructNotEmpty(dst, src any) {
// 	copyStructNotEmpty(copyStructCheck(dst, src))
// }

// func copyStructNotEmpty(dst, src reflect.Value) {
// 	// type
// 	srcType := src.Type()
// 	for i := 0; i < srcType.NumField(); i++ {
// 		srcField := src.Field(i)
// 		// src 零值
// 		if !srcField.IsValid() || srcField.IsZero() {
// 			continue
// 		}
// 		srcTypeField := srcType.Field(i)
// 		dstField := dst.FieldByName(srcTypeField.Name)
// 		if !dstField.IsValid() || !dstField.CanSet() {
// 			continue
// 		}
// 		dstFieldType := dstField.Type()
// 		// 不同类型
// 		if srcTypeField.Type != dstFieldType {
// 			srcFieldKind := srcField.Kind()
// 			// 看看是不是结构体
// 			if srcFieldKind == reflect.Pointer {
// 				srcField = srcField.Elem()
// 				srcFieldKind = srcField.Kind()
// 			}
// 			dstFieldKind := dstField.Kind()
// 			if dstFieldKind == reflect.Pointer {
// 				dstField = dstField.Elem()
// 				dstFieldKind = dstField.Kind()
// 			}
// 			// 都是结构体，进去赋值
// 			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
// 				copyStructNotEmpty(dstField, srcField)
// 			}
// 			// 不是就算了
// 			continue
// 		}
// 		// 相同类型，赋值
// 		dstField.Set(srcField)
// 	}
// }
