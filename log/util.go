package log

import (
	"strconv"
	"strings"
)

const (
	// 单位
	kb = 1024
	mb = 1024 * kb
	gb = 1024 * mb
	tb = 1024 * gb
)

// parseSize 解析 1K ，10M 字符串形式到整数字节。
func ParseSize(str string) (int64, error) {
	str = strings.ToUpper(str)
	// TB
	n, err := parseSize(str, "TB", tb)
	if err != nil {
		return n, err
	}
	if n > 0 {
		return n, nil
	}
	// GB
	n, err = parseSize(str, "GB", gb)
	if err != nil {
		return n, err
	}
	if n > 0 {
		return n, nil
	}
	// MB
	n, err = parseSize(str, "MB", mb)
	if err != nil {
		return n, err
	}
	if n > 0 {
		return n, nil
	}
	// KB
	n, err = parseSize(str, "KB", kb)
	if err != nil {
		return n, err
	}
	if n > 0 {
		return n, nil
	}
	// 没有单位，字节
	return strconv.ParseInt(str, 10, 64)
}

// parseSize 是 ParseSize 的辅助函数。
func parseSize(str, unit string, bytes int64) (int64, error) {
	// XB
	p := strings.TrimSuffix(str, unit)
	if p != str {
		n, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, err
		}
		return int64(n * float64(bytes)), nil
	}
	// X
	p = strings.TrimSuffix(str, unit[:1])
	if p != str {
		n, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return 0, err
		}
		return int64(n * float64(bytes)), nil
	}
	return -1, nil
}
