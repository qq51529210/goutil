package reflect

import (
	"reflect"
	"strings"
)

// ParseTag 解析 "name,omitempty","-",",omitempty"
func ParseTag(f *reflect.StructField, tagName string) (name string, omitempty, ignore bool) {
	tag := f.Tag.Get(tagName)
	// ignore
	if tag == "-" {
		ignore = true
		return
	}
	// name
	i := strings.IndexByte(tag, ',')
	if i < 0 {
		name = tag
		return
	}
	name = tag[:i]
	// omitempty
	omitempty = tag[i+1:] == "omitempty"
	//
	return
}
