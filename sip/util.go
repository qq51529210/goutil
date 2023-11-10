package sip

import (
	"fmt"
	"strings"
	"sync/atomic"
)

var (
	// CSeq 的递增 SN
	sn32   = int64(0)
	cseq32 = int64(0)
)

// GetSN 返回全局递增的 sn
func GetSN() int64 {
	return atomic.AddInt64(&sn32, 1)
}

// GetCSeq 返回全局递增的 sn
func GetCSeq() int64 {
	return atomic.AddInt64(&cseq32, 1)
}

// GetSNString 返回字符串形式的全局递增的 sn
func GetSNString() string {
	return fmt.Sprintf("%d", atomic.AddInt64(&sn32, 1))
}

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

// TrimQuotationMark TrimByte(s, '"', '"')
func TrimQuotationMark(s string) string {
	return TrimByte(s, '"', '"')
}

// 字符
const (
	CharComma     = ','
	CharSpace     = ' '
	CharEqual     = '='
	CharSemicolon = ';'
	CharAmpersand = '&'
	CharColon     = ':'
	CharAt        = '@'
)

// Split 根据第一个 c 切分 s
func Split(s string, c byte) (prefix, suffix string) {
	i := strings.IndexByte(s, c)
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i+1:]
}
