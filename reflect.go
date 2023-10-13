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
		zero := fv.IsZero()
		// tag
		name, omitempty, ignore := structToMapTag(&ft)
		// 忽略字段
		if ignore {
			continue
		}
		// 指针
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 忽略空值
		if omitempty && zero {
			continue
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

func copyStructCheck(dst, src any) (dstVal, srcVal reflect.Value) {
	// dst
	dstVal = reflect.ValueOf(dst)
	if dstVal.Kind() != reflect.Pointer {
		panic("dst must be pointer")
	}
	dstVal = dstVal.Elem()
	if dstVal.Kind() != reflect.Struct {
		panic("dst must be struct pointer")
	}
	// src
	srcVal = reflect.ValueOf(src)
	if srcVal.Kind() != reflect.Pointer {
		panic("src must be pointer")
	}
	srcVal = srcVal.Elem()
	if srcVal.Kind() != reflect.Struct {
		panic("src must be struct pointer")
	}
	//
	return
}

// CopyStruct 拷贝 src 和 dst 中的相同名称和类型的字段，
// 如果 dst 的字段不是零值则不拷贝。
func CopyStruct(dst, src any) {
	copyStruct(copyStructCheck(dst, src))
}

// copyStruct 封装 CopyStruct 代码
func copyStruct(dst, src reflect.Value) {
	// 结构类型
	srcType := src.Type()
	for i := 0; i < srcType.NumField(); i++ {
		srcField := src.Field(i)
		// src 零值
		if !srcField.IsValid() || srcField.IsZero() {
			continue
		}
		srcTypeField := srcType.Field(i)
		// dst 同名字段
		dstField := dst.FieldByName(srcTypeField.Name)
		// dst 不为零
		if !dstField.IsValid() || !dstField.CanSet() || !dstField.IsZero() {
			continue
		}
		dstFieldType := dstField.Type()
		// 不同类型
		if srcTypeField.Type != dstFieldType {
			srcFieldKind := srcField.Kind()
			// 看看是不是结构体
			if srcFieldKind == reflect.Pointer {
				srcField = srcField.Elem()
				srcFieldKind = srcField.Kind()
			}
			dstFieldKind := dstField.Kind()
			if dstFieldKind == reflect.Pointer {
				dstField = dstField.Elem()
				dstFieldKind = dstField.Kind()
			}
			// 都是结构体，进去赋值
			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
				copyStruct(dstField, srcField)
			}
			// 不是就算了
			continue
		}
		// 相同类型，赋值
		dstField.Set(srcField)
	}
}

// CopyStructAll 拷贝 src 和 dst 中的相同名称和类型的字段
func CopyStructAll(dst, src any) {
	copyStructAll(copyStructCheck(dst, src))
}

// copyStructAll 封装 CopyStructAll 代码
func copyStructAll(dst, src reflect.Value) {
	// type
	srcType := src.Type()
	for i := 0; i < srcType.NumField(); i++ {
		srcField := src.Field(i)
		if !srcField.IsValid() {
			continue
		}
		srcTypeField := srcType.Field(i)
		dstField := dst.FieldByName(srcTypeField.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}
		dstFieldType := dstField.Type()
		// 不同类型
		if srcTypeField.Type != dstFieldType {
			srcFieldKind := srcField.Kind()
			// 看看是不是结构体
			if srcFieldKind == reflect.Pointer {
				srcField = srcField.Elem()
				srcFieldKind = srcField.Kind()
			}
			dstFieldKind := dstField.Kind()
			if dstFieldKind == reflect.Pointer {
				dstField = dstField.Elem()
				dstFieldKind = dstField.Kind()
			}
			// 都是结构体，进去赋值
			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
				copyStructAll(dstField, srcField)
			}
			continue
		}
		// 相同类型，赋值
		dstField.Set(srcField)
	}
}

// CopyStructNotEmpty 拷贝 src 和 dst 中的相同名称和类型的字段
// src 为零值不拷贝，dst 不为零值，也拷贝哦
func CopyStructNotEmpty(dst, src any) {
	copyStructNotEmpty(copyStructCheck(dst, src))
}

func copyStructNotEmpty(dst, src reflect.Value) {
	// type
	srcType := src.Type()
	for i := 0; i < srcType.NumField(); i++ {
		srcField := src.Field(i)
		// src 零值
		if !srcField.IsValid() || srcField.IsZero() {
			continue
		}
		srcTypeField := srcType.Field(i)
		dstField := dst.FieldByName(srcTypeField.Name)
		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}
		dstFieldType := dstField.Type()
		// 不同类型
		if srcTypeField.Type != dstFieldType {
			srcFieldKind := srcField.Kind()
			// 看看是不是结构体
			if srcFieldKind == reflect.Pointer {
				srcField = srcField.Elem()
				srcFieldKind = srcField.Kind()
			}
			dstFieldKind := dstField.Kind()
			if dstFieldKind == reflect.Pointer {
				dstField = dstField.Elem()
				dstFieldKind = dstField.Kind()
			}
			// 都是结构体，进去赋值
			if srcFieldKind == reflect.Struct && dstFieldKind == reflect.Struct {
				copyStructNotEmpty(dstField, srcField)
			}
			// 不是就算了
			continue
		}
		// 相同类型，赋值
		dstField.Set(srcField)
	}
}

// StructToMapIgnore 将 v 转换为 map，v 必须是结构体
// 忽略零值字段
func StructToMapIgnore(v any) map[string]any {
	return structToMapIgnore(reflect.ValueOf(v))
}

// structToMapIgnore 封装 StructToMapIgnore 的代码
func structToMapIgnore(vVal reflect.Value) map[string]any {
	if vVal.Kind() == reflect.Pointer {
		vVal = vVal.Elem()
	}
	vType := vVal.Type()
	result := make(map[string]any)
	for i := 0; i < vType.NumField(); i++ {
		field := vVal.Field(i)
		if !field.IsValid() || field.IsZero() {
			continue
		}
		fieldName := vType.Field(i).Name
		fieldKind := field.Kind()
		if fieldKind == reflect.Pointer {
			field = field.Elem()
			fieldKind = field.Kind()
			if fieldKind == reflect.Invalid {
				continue
			}
		}
		if fieldKind == reflect.Struct {
			result[fieldName] = structToMapIgnore(field)
		} else {
			result[fieldName] = field.Interface()
		}
	}
	return result
}

// StructDiffFieldsIgnore 找出 v1 和 v2 中不同值的字段，然后返回 v2 的值
// 忽略零值字段
// v1 和 v2 必须是同一种类型结构，判断一层
func StructDiffFieldsIgnore(v1, v2 any) map[string]any {
	// v1
	v1Val := reflect.ValueOf(v1)
	v1Kind := v1Val.Kind()
	if v1Kind == reflect.Pointer {
		v1Val = v1Val.Elem()
		v1Kind = v1Val.Kind()
	}
	if v1Kind != reflect.Struct {
		panic("v1 must be struct")
	}
	// v2
	v2Val := reflect.ValueOf(v2)
	v2Kind := v2Val.Kind()
	if v2Kind == reflect.Pointer {
		v2Val = v2Val.Elem()
		v2Kind = v2Val.Kind()
	}
	if v2Kind != reflect.Struct {
		panic("v2 must be struct")
	}
	// 同一类型
	if v1Val.Type() != v2Val.Type() {
		panic("v1 v2 must be same type")
	}
	//
	return structDiffFieldsIgnore(v1Val, v2Val, make(map[string]any))
}

// structDiffFields 封装 StructDiffFields 的代码
func structDiffFieldsIgnore(v1, v2 reflect.Value, m map[string]any) map[string]any {
	_type := v1.Type()
	for i := 0; i < _type.NumField(); i++ {
		v1Field, v2Field := v1.Field(i), v2.Field(i)
		if !v1Field.IsValid() || !v2Field.IsValid() || v2Field.IsZero() {
			continue
		}
		kind := v1Field.Kind()
		if kind == reflect.Pointer {
			v1Field = v1Field.Elem()
			v2Field = v2Field.Elem()
		}
		kind = v1Field.Kind()
		if kind == reflect.Struct {
			structDiffFieldsIgnore(v1Field, v2Field, m)
		} else {
			v1Value := v1Field.Interface()
			v2Value := v2Field.Interface()
			if (kind >= reflect.Bool && kind <= reflect.Uint64) ||
				kind == reflect.Float32 || kind == reflect.Float64 ||
				kind == reflect.String {
				if v1Value != v2Value {
					m[_type.Field(i).Name] = v2Value
				}
			}
		}
	}
	return m
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
	structFieldValue(vv)
}

// structFieldValue 封装 StructFieldValue 的代码
func structFieldValue(structValue reflect.Value) {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		value := structField.Tag.Get("value")
		if value != "" {
			structFieldtDefaultValue(structValue.Field(i), value)
			continue
		}
		value = structField.Tag.Get("field")
		if value != "" {
			structFieldDefaultField(structValue, structValue.Field(i), value)
			continue
		}
		value = structField.Tag.Get("func")
		if value != "" {
			structFieldDefaultFunc(structValue, structValue.Field(i), value)
			continue
		}
	}
}

// structFieldtDefaultValue 设置字段的为指定的值
func structFieldtDefaultValue(fieldValue reflect.Value, defaultValue string) {
	kind := fieldValue.Kind()
	// nil 指针
	if kind == reflect.Pointer {
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
	dstKind := fieldValue.Kind()
	if dstKind == reflect.Pointer {
		if fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}
		fieldValue = fieldValue.Elem()
		dstKind = fieldValue.Kind()
	}
	srcKind := srcFieldValue.Kind()
	if srcKind == reflect.Pointer {
		srcFieldValue = srcFieldValue.Elem()
		srcKind = srcFieldValue.Kind()
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
