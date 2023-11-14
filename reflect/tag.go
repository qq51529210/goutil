package reflect

import (
	"reflect"
	"strings"
)

// parseTag 解析 tag 的封装
func parseTag(f *reflect.StructField, tagName string) (name string, omitempty, ignore bool) {
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
