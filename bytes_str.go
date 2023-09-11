package util

import (
	"errors"
	"fmt"
	"strconv"
)

const (
	bytesK = 1024
	bytesM = 1024 * bytesK
	bytesG = 1024 * bytesM
	bytesT = 1024 * bytesG
)

var (
	errEmptyArgument = errors.New("empty argument")
)

// StringToByte 解析 2KB，3.13M，4.5G，5.67Tb 这样的字符串到数值
func StringToByte(s string) (uint64, error) {
	if s == "" {
		return 0, errEmptyArgument
	}
	i := len(s) - 1
	if s[i] == 'B' || s[i] == 'b' {
		i--
		if i < 0 {
			return 0, fmt.Errorf("invalid argument '%v'", s)
		}
	}
	switch s[i] {
	case 'K', 'k':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * bytesK), nil
	case 'M', 'm':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * bytesM), nil
	case 'G', 'g':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * bytesG), nil
	case 'T', 't':
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v * bytesT), nil
	default:
		v, e := strconv.ParseFloat(s[:i], 64)
		if nil != e {
			return 0, e
		}
		return uint64(v), nil
	}
}

// ByteToString 格式化数值1234567，这样的到字符串1.177...MB，保留 p 位小数
func ByteToString(n uint64, p int) string {
	if n > bytesT {
		return strconv.FormatFloat(float64(n)/bytesT, 'f', p, 64) + "T"
	}
	if n > bytesG {
		return strconv.FormatFloat(float64(n)/bytesG, 'f', p, 64) + "G"
	}
	if n > bytesM {
		return strconv.FormatFloat(float64(n)/bytesM, 'f', p, 64) + "M"
	}
	if n > bytesK {
		return strconv.FormatFloat(float64(n)/bytesK, 'f', p, 64) + "K"
	}
	return strconv.FormatUint(n, 10)
}
