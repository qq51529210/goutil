package util

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// structParseTag 解析 tag 的封装
func structParseTag(f *reflect.StructField, tagName string) (name string, omitempty, ignore bool) {
	name = f.Name
	tag := f.Tag.Get(tagName)
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

// IsNilOrEmpty 如果 v 是空指针或者零值，返回 true
// 指针有零值不算零值
func IsNilOrEmpty(v any) bool {
	return isNilOrEmpty(reflect.ValueOf(v))
}

// isNilOrEmpty 是 IsNilOrEmpty 的实现
func isNilOrEmpty(v reflect.Value) bool {
	// 无效
	if !v.IsValid() {
		return true
	}
	k := v.Kind()
	pk := k
	// 指针
	if k == reflect.Pointer {
		v = v.Elem()
		k = v.Kind()
		// nil
		if !v.IsValid() {
			return true
		}
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
		// 其他，指针不是 nil 就不是零
		if pk == reflect.Pointer {
			return false
		}
		return v.IsZero()
	}
}

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
	m := make(map[string]any)
	return structToMap(reflect.ValueOf(v), m)
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
		name, omitempty, ignore := structParseTag(&ft, StructToMapTagName)
		// 忽略
		if ignore {
			continue
		}
		// 值
		fv := v.Field(i)
		// 无效/零值
		if !fv.IsValid() || (omitempty && fv.IsZero()) {
			continue
		}
		// 数据类型
		fk := ft.Type.Kind()
		if fk == reflect.Pointer {
			fv = fv.Elem()
			fk = ft.Type.Elem().Kind()
		}
		// 嵌入字段
		if ft.Anonymous {
			// 嵌入的必须是结构
			if fk == reflect.Struct {
				// 而且不能是 nil 指针
				if !fv.IsValid() {
					continue
				}
				m = structToMap(fv, m)
			}
			continue
		}
		// 不可导出
		if !ft.IsExported() {
			continue
		}
		// 其他
		if fv.IsValid() && fv.CanInterface() {
			// 结构
			if fk == reflect.Struct {
				m[name] = structToMap(fv, make(map[string]any))
				continue
			}
			m[name] = fv.Interface()
		} else {
			m[name] = nil
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
	//
	structCopy(dv, sv)
}

// structCopy 是 StructCopy 的实现
func structCopy(dstStructValue, srcStructValue reflect.Value) {
	// 类型
	srcStructType := srcStructValue.Type()
	// 所有字段
	for i := 0; i < srcStructType.NumField(); i++ {
		// src 字段类型
		srcFieldType := srcStructType.Field(i)
		// tag
		name, omitempty, ignore := structParseTag(&srcFieldType, StructCopyTagName)
		// 忽略字段
		if ignore {
			continue
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
				structCopy(dstStructValue, srcFieldValue)
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
			structCopy(dstFieldValue, srcFieldValue)
			continue
		}
		// 其他类型赋值
		dstFieldValue.Set(srcFieldValue)
	}
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
		name, omitempty, ignore := structParseTag(&sft, StructDiffFieldTagName)
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
