package reflect

import "reflect"

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
	// 指针不是 nil 就不是零
	if k == reflect.Pointer {
		return v.IsNil()
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
		return v.IsZero()
	}
}
