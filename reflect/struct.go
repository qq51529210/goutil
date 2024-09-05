package reflect

// import "reflect"

// func StructWalker(s any, tag string, cb func(name string, value any)) {
// 	v := reflect.ValueOf(s)
// 	k := v.Kind()
// 	if k == reflect.Pointer {
// 		// 空指针
// 		if v.IsNil() {
// 			return
// 		}
// 		v = v.Elem()
// 	}
// 	t := v.Type()
// 	for i := 0; i < t.NumField(); i++ {
// 		ft := t.Field(i)
// 		// tag
// 		name, omitempty, ignore := ParseTag(&ft, tag)
// 		// 忽略字段
// 		if ignore {
// 			continue
// 		}
// 		fv := v.Field(i)

// 	}
// }
