package reflect

// import (
// 	"fmt"
// 	"reflect"
// 	"strconv"
// )

// // StructFieldValue 给字段赋值，v 必须是结构指针
// // 如果字段为零值，或者 nil 指针，则赋值
// //
// //	type S1 struct {
// //	   A *string `value:"123"` -> A 为 nil 赋值 123
// //	   B string `field:"A"` -> B 为 零值，将 A 的值给它
// //	   C string `func:"Default"` -> C 为 零值，调用结构体方法赋值，接收者不能是指针
// //	}
// //
// //	func (s S1)Default()string {...}
// func StructFieldValue(v any) {
// 	sv := reflect.ValueOf(v)
// 	if sv.Kind() != reflect.Pointer {
// 		panic("v must be pointer")
// 	}
// 	sv = sv.Elem()
// 	if sv.Kind() != reflect.Struct {
// 		panic("v must be struct pointer")
// 	}
// 	//
// 	st := sv.Type()
// 	for i := 0; i < st.NumField(); i++ {
// 		fv := sv.Field(i)
// 		// 无效/不为空
// 		if !fv.IsValid() || !fv.IsZero() {
// 			continue
// 		}
// 		ft := st.Field(i)
// 		tag := ft.Tag.Get("value")
// 		if tag != "" {
// 			structFieldtDefaultValue(fv, tag)
// 			continue
// 		}
// 		tag = ft.Tag.Get("field")
// 		if tag != "" {
// 			structFieldDefaultField(sv, fv, tag)
// 			continue
// 		}
// 		tag = ft.Tag.Get("func")
// 		if tag != "" {
// 			structFieldDefaultFunc(sv.MethodByName(tag), fv, tag)
// 			continue
// 		}
// 	}
// }

// // structFieldtDefaultValue 设置字段的为指定的值
// func structFieldtDefaultValue(fieldValue reflect.Value, defaultValue string) {
// 	kind := fieldValue.Kind()
// 	// 指针
// 	if kind == reflect.Pointer {
// 		// nil new
// 		if fieldValue.IsNil() {
// 			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
// 		}
// 		fieldValue = fieldValue.Elem()
// 		kind = fieldValue.Kind()
// 	}
// 	// 字符串->数字
// 	if kind >= reflect.Uint && kind <= reflect.Uint64 {
// 		n, err := strconv.ParseInt(defaultValue, 10, 64)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fieldValue.SetUint(uint64(n))
// 		return
// 	}
// 	if kind >= reflect.Int && kind <= reflect.Int64 {
// 		n, err := strconv.ParseInt(defaultValue, 10, 64)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fieldValue.SetInt(n)
// 		return
// 	}
// 	if kind >= reflect.Float32 && kind <= reflect.Float64 {
// 		n, err := strconv.ParseFloat(defaultValue, 64)
// 		if err != nil {
// 			panic(err)
// 		}
// 		fieldValue.SetFloat(n)
// 		return
// 	}
// 	// 字符串
// 	if kind == reflect.String {
// 		fieldValue.SetString(defaultValue)
// 		return
// 	}
// 	//
// 	if kind == reflect.Bool {
// 		fieldValue.SetBool(defaultValue == "true")
// 		return
// 	}
// 	panic(fmt.Sprintf("unsupported field type %v", kind))
// }

// // structFieldDefaultField 设置字段的为指定的字段的值
// func structFieldDefaultField(structValue, fieldValue reflect.Value, fieldName string) {
// 	// 找到标记的字段
// 	srcFieldValue := structValue.FieldByName(fieldName)
// 	// 无效就算了
// 	if !srcFieldValue.IsValid() {
// 		return
// 	}
// 	srcKind := srcFieldValue.Kind()
// 	if srcKind == reflect.Pointer {
// 		srcFieldValue = srcFieldValue.Elem()
// 		srcKind = srcFieldValue.Kind()
// 	}
// 	//
// 	dstKind := fieldValue.Kind()
// 	if dstKind == reflect.Pointer {
// 		if fieldValue.IsNil() {
// 			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
// 		}
// 		fieldValue = fieldValue.Elem()
// 		dstKind = fieldValue.Kind()
// 	}
// 	// 不相同
// 	if dstKind != srcKind {
// 		panic(fmt.Sprintf("field %s different type %v", fieldName, srcKind))
// 	}
// 	// 相同
// 	if (dstKind >= reflect.Bool && dstKind <= reflect.Uint64) ||
// 		(dstKind >= reflect.Float32 && dstKind <= reflect.Float64) ||
// 		dstKind == reflect.String {
// 		fieldValue.Set(srcFieldValue)
// 		return
// 	}
// 	panic(fmt.Sprintf("unsupported field type %v", dstKind))
// }

// // structFieldDefaultFunc 设置字段的为指定的字段的值
// func structFieldDefaultFunc(structFunc, fieldValue reflect.Value, funcName string) {
// 	// nil 指针
// 	kind := fieldValue.Kind()
// 	if kind == reflect.Pointer {
// 		if fieldValue.IsNil() {
// 			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
// 		}
// 		fieldValue = fieldValue.Elem()
// 	}
// 	fieldValue.Set(structFunc.Call(nil)[0])
// }
