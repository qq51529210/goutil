package util

import (
	"reflect"
)

// IsNilOrEmpty 如果 v 是空指针，或者空值，返回 true
// 指针的值是空值，不算空值，也返回 true
func IsNilOrEmpty(v any) bool {
	return isNilOrEmpty(reflect.ValueOf(v))
}

// isNilOrEmpty 如果 v 是空指针，或者空值，返回 true
// 指针的值是空值，不算空值，也返回 true
func isNilOrEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			return true
		}
		v = v.Elem()
		switch v.Kind() {
		case reflect.Struct:
			return isStructNilOrEmpty(v)
		default:
			return false
		}
	case reflect.Func:
		return v.IsNil()
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Struct:
		return isStructNilOrEmpty(v)
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.IsZero()
	}
	return false
}

// isStructNilOrEmpty 封装 IsNilOrEmpty 中判断 struct 的代码
func isStructNilOrEmpty(v reflect.Value) bool {
	for i := 0; i < v.NumField(); i++ {
		if !isNilOrEmpty(v.Field(i)) {
			return false
		}
	}
	return true
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

// StructToMap 将 v 转换为 map，v 必须是结构体
func StructToMap(v any) map[string]any {
	return structToMap(reflect.ValueOf(v))
}

// structToMap 封装 StructToMap 的代码
func structToMap(vVal reflect.Value) map[string]any {
	if vVal.Kind() == reflect.Pointer {
		vVal = vVal.Elem()
	}
	vType := vVal.Type()
	result := make(map[string]any)
	for i := 0; i < vType.NumField(); i++ {
		field := vVal.Field(i)
		if !field.IsValid() {
			continue
		}
		fieldName := vType.Field(i).Name
		fieldKind := field.Kind()
		if fieldKind == reflect.Pointer {
			field = field.Elem()
			fieldKind = field.Kind()
			if fieldKind == reflect.Invalid {
				result[fieldName] = nil
				continue
			}
		}
		if fieldKind == reflect.Struct {
			result[fieldName] = structToMap(field)
		} else {
			result[fieldName] = field.Interface()
		}
	}
	return result
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
		if v1Field.Kind() == reflect.Pointer {
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
