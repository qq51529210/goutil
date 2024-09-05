package reflect

import "reflect"

var (
	// StructCopyTagName 是 StructCopy 结构解析的 tag 名称
	StructCopyTagName = "copy"
)

// StructCopy 将 src 的字段拷贝到 dst
// 嵌入的字段不是结构，或者是 nil 的结构指针，不处理
// 不可导出的字段，不处理
//
//	type src struct {
//	   A string 默认拷贝到相同名称 dst.A
//	   B string `copy:"BB"` -> 指定字段 dst.BB
//	   C string `copy:"omitempty"` -> 忽略零值
//	   D *string `copy:"d,omitempty"` -> 指针不为 nil 不算零值
//	   E string `copy:"-"` -> 忽略
//	   f *common 不导出，忽略
//	   *common 忽略 nil
//	}
//
//	type dst struct {
//	   A *string 有拷贝则自动 new
//	   BB string
//	   C string
//	   d *string 不导出，忽略
//	   D *string
//	   E string
//	   A1 string
//	   B1 string
//	}
//
//	type common struct {
//	   A1 string
//	   B1 string
//	}
func StructCopy(dst, src any) {
	StructCopyWithTag(dst, src, StructCopyTagName)
}

// StructCopyWithTag 使用自定义 tag
func StructCopyWithTag(dst, src any, tag string) {
	dv := reflect.ValueOf(dst)
	if !dv.IsValid() {
		panic("dst is invalid")
	}
	dk := dv.Kind()
	if dk != reflect.Pointer {
		panic("dst must be pointer")
	}
	dv = dv.Elem()
	dk = dv.Kind()
	if dk != reflect.Struct {
		panic("dst must be struct pointer")
	}
	//
	sv := reflect.ValueOf(src)
	if !sv.IsValid() {
		panic("src is invalid")
	}
	sk := sv.Kind()
	if sk == reflect.Pointer {
		sv = sv.Elem()
	}
	//
	//
	structCopy(dv, sv, tag)
}

// structCopy 是 StructCopy 的实现
func structCopy(dstStructValue, srcStructValue reflect.Value, tag string) {
	// 类型
	srcStructType := srcStructValue.Type()
	// 所有字段
	for i := 0; i < srcStructType.NumField(); i++ {
		// src 字段类型
		srcFieldType := srcStructType.Field(i)
		// tag
		name, omitempty, ignore := ParseTag(&srcFieldType, tag)
		// 忽略字段
		if ignore {
			continue
		}
		if name == "" {
			name = srcFieldType.Name
		}
		// src 字段值
		srcFieldValue := srcStructValue.Field(i)
		srcFieldKind := srcFieldValue.Kind()
		if srcFieldKind == reflect.Pointer {
			srcFieldKind = srcFieldType.Type.Elem().Kind()
			// nil 指针
			if srcFieldValue.IsNil() {
				// 结构
				if srcFieldKind == reflect.Struct {
					continue
				}
				// 嵌入 / 忽略零值 / 不可导出
				if srcFieldType.Anonymous || omitempty || !srcFieldType.IsExported() {
					continue
				}
			}
			srcFieldValue = srcFieldValue.Elem()
			// 往下
		} else {
			// 无论导出，处理嵌入的结构
			if srcFieldType.Anonymous && srcFieldKind == reflect.Struct {
				structCopy(dstStructValue, srcFieldValue, tag)
				continue
			}
			// 不可导出 / 忽略零值
			if !srcFieldType.IsExported() || (omitempty && srcFieldValue.IsZero()) {
				continue
			}
			// 往下
		}
		// dst 值
		dstFieldValue := dstStructValue.FieldByName(name)
		// 找不到
		if !dstFieldValue.IsValid() || !dstFieldValue.CanSet() {
			continue
		}
		dstFieldKind := dstFieldValue.Kind()
		if dstFieldKind == reflect.Pointer {
			// 需要 new 一个接值
			if dstFieldValue.IsNil() {
				// src 是无效，不处理
				if !srcFieldValue.IsValid() {
					continue
				}
				dataType := dstFieldValue.Type().Elem()
				dstFieldKind = dataType.Kind()
				// 数据类型不同
				if dstFieldKind != srcFieldKind {
					continue
				}
				dstFieldValue.Set(reflect.New(dataType))
				dstFieldValue = dstFieldValue.Elem()
			} else {
				// src 是无效，但是没有忽略零值，这里可能需要设置 dst 为 nil
				if !srcFieldValue.IsValid() {
					dataType := dstFieldValue.Type().Elem()
					dstFieldKind = dataType.Kind()
					// 数据类型不同
					if dstFieldKind != srcFieldKind {
						continue
					}
					dstFieldValue.Set(reflect.Zero(srcFieldType.Type))
					continue
				}
				// src 有效
				dstFieldValue = dstFieldValue.Elem()
				dstFieldKind = dstFieldValue.Kind()
				// 数据类型不同
				if dstFieldKind != srcFieldKind {
					continue
				}
			}
			// 往下
		} else {
			// src 是无效，不处理
			if !srcFieldValue.IsValid() {
				continue
			}
			// 数据类型不同
			if dstFieldKind != srcFieldKind {
				continue
			}
			// 往下
		}
		// 结构拷贝
		if srcFieldKind == reflect.Struct {
			structCopy(dstFieldValue, srcFieldValue, tag)
			continue
		}
		// 其他类型赋值
		dstFieldValue.Set(srcFieldValue)
	}
}
