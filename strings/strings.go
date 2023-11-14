package strings

import "strings"

// 字符
const (
	CharComma     = ','
	CharSpace     = ' '
	CharEqual     = '='
	CharSemicolon = ';'
	CharAmpersand = '&'
	CharColon     = ':'
	CharAt        = '@'
	QuotationMark = '"'
)

// TrimByte 去掉两端的字符
func TrimByte(str string, left, right byte) string {
	if str == "" {
		return str
	}
	i1 := 0
	for ; i1 < len(str); i1++ {
		if str[i1] != left {
			break
		}
	}
	i2 := len(str) - 1
	for ; i2 > i1; i2-- {
		if str[i2] != right {
			break
		}
	}
	return str[i1 : i2+1]
}

// Split 根据第一个 c 切分 s
func Split(s string, c byte) (prefix, suffix string) {
	i := strings.IndexByte(s, c)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}
