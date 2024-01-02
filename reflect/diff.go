package reflect

import (
	"reflect"
)

// var (
// 	// StructDiffFieldTagName 是 StructDiffField 结构解析的 tag 名称
// 	StructDiffFieldTagName = "diff"
// )

// StructDiffField 从 src 中找出 与 dst 不同值的字段，然后返回这些字段的 map
// src 和 dst 是相同数据类型的结构
func StructDiffField(dst, src any, tag string) map[string]any {
	// dst
	dv := reflect.ValueOf(dst)
	dk := dv.Kind()
	if dk == reflect.Pointer {
		dv = dv.Elem()
		dk = dv.Kind()
	}
	if dk != reflect.Struct {
		panic("dst must be struct")
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
	// type
	if dv.Type() != sv.Type() {
		panic("src dst must be same struct type")
	}
	//
	return structDiffField(dv, sv, tag, make(map[string]any))
}

// structDiffField 是 StructDiffFields 的实现
func structDiffField(dstStructValue, srcStructValue reflect.Value, tag string, m map[string]any) map[string]any {
	srcStructType := srcStructValue.Type()
	for i := 0; i < srcStructType.NumField(); i++ {
		// 字段类型
		fieldType := srcStructType.Field(i)
		// 忽略
		if fieldType.Tag.Get(tag) == "-" {
			continue
		}
		fieldKind := fieldType.Type.Kind()
		srcFieldValue := srcStructValue.Field(i)
		dstFieldValue := dstStructValue.Field(i)
		// 指针
		if fieldKind == reflect.Pointer {
			fieldKind = fieldType.Type.Elem().Kind()
			srcFieldValue = srcFieldValue.Elem()
			dstFieldValue = dstFieldValue.Elem()
			if srcFieldValue.IsValid() {
				if !dstFieldValue.IsValid() {
					// src 有效，dst 无效
					// 是否结构
					if fieldKind == reflect.Struct {
						// 嵌入
						if fieldType.Anonymous {
							structDiffField(dstStructValue, srcFieldValue, tag, m)
						} else {
							if !srcFieldValue.CanInterface() {
								continue
							}
							m[fieldType.Name] = srcFieldValue.Interface()
						}
					} else {
						if !srcFieldValue.CanInterface() {
							continue
						}
						m[fieldType.Name] = srcFieldValue.Interface()
					}
				} else {
					// src 有效，dst 有效
					// 是否结构
					if fieldKind == reflect.Struct {
						// 嵌入
						if fieldType.Anonymous {
							structDiffField(dstStructValue, srcFieldValue, tag, m)
						} else {
							if !srcFieldValue.CanInterface() {
								continue
							}
							m[fieldType.Name] = srcFieldValue.Interface()
						}
					} else {
						if !srcFieldValue.CanInterface() || !dstFieldValue.CanInterface() {
							continue
						}
						// 比较
						srcData, dstData := srcFieldValue.Interface(), dstFieldValue.Interface()
						if !reflect.DeepEqual(srcData, dstData) {
							m[fieldType.Name] = srcData
						}
					}
				}
				continue
			} else {
				// src 无效，dst 有效
				if dstFieldValue.IsValid() {
					m[fieldType.Name] = nil
				}
			}
			continue
		}
		// 不是指针
		if fieldKind == reflect.Struct {
			// 嵌入
			if fieldType.Anonymous {
				structDiffField(dstFieldValue, srcFieldValue, tag, m)
				continue
			}
		}
		// 比较
		srcData, dstData := srcFieldValue.Interface(), dstFieldValue.Interface()
		if !reflect.DeepEqual(srcData, dstData) {
			m[fieldType.Name] = srcData
		}
	}
	//
	return m
}
