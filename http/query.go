package http

import (
	"fmt"
	"net/url"
	"reflect"
)

var (
	// QueryTag 是 Query 解析 tag 的名称
	QueryTag = "query"
)

// Query 将结构体 v 格式化到 url.Values
// 只扫描一层，并略过空值
func Query(v any, q url.Values) url.Values {
	return QueryWithTag(v, q, QueryTag)
}

// QueryWithTag 指定 tag
func QueryWithTag(v any, q url.Values, tag string) url.Values {
	rv := reflect.ValueOf(v)
	vk := rv.Kind()
	if vk == reflect.Invalid {
		return q
	}
	if vk == reflect.Pointer {
		rv = rv.Elem()
		vk = rv.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	if q == nil {
		q = make(url.Values)
	}
	return httpQuery(rv, q, tag)
}

func httpQuery(v reflect.Value, q url.Values, tag string) url.Values {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构，只一层
		if fvk == reflect.Struct {
			continue
		}
		if fvk == reflect.String {
			// 空字符串
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(tag)
		if tn == "" || tn == "-" {
			continue
		}
		q.Add(tn, fmt.Sprintf("%v", fv.Interface()))
	}
	return q
}
