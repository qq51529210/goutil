package strings

import (
	"strings"

	"github.com/mozillazg/go-pinyin"
)

// ToPinyin 转成字母拼音
func ToPinyin(str string) string {
	var buf strings.Builder
	a := pinyin.NewArgs()
	for _, s := range str {
		buf.WriteString("/")
		ps := pinyin.Pinyin(string(s), a)
		if len(ps) < 1 {
			buf.WriteRune(s)
			continue
		}
		for _, p := range ps {
			for i := 0; i < len(p); i++ {
				buf.WriteString(p[i])
			}
		}
	}
	return buf.String()[1:]
}
