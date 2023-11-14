package hash

import (
	"crypto/md5"
	"sync"
)

var (
	md5Pool sync.Pool
)

func init() {
	md5Pool.New = func() interface{} {
		return &buffer{
			hash: md5.New(),
			buf:  make([]byte, 0, md5.Size*2),
			sum:  make([]byte, md5.Size),
		}
	}
}

// MD5 返回 16 进制哈希字符串
func MD5(b []byte) string {
	h := md5Pool.Get().(*buffer)
	s := h.Hash(b)
	md5Pool.Put(h)
	return s
}

// MD5String 返回 16 进制哈希字符串
func MD5String(s string) string {
	h := md5Pool.Get().(*buffer)
	s = h.HashString(s)
	md5Pool.Put(h)
	return s
}
