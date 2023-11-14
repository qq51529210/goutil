package hash

import (
	"crypto/sha1"
	"sync"
)

var (
	sha1Pool sync.Pool
)

func init() {
	sha1Pool.New = func() interface{} {
		return &buffer{
			hash: sha1.New(),
			buf:  make([]byte, 0, sha1.Size*2),
			sum:  make([]byte, sha1.Size),
		}
	}
}

// SHA1 返回 16 进制哈希字符串
func SHA1(b []byte) string {
	h := sha1Pool.Get().(*buffer)
	s := h.Hash(b)
	sha1Pool.Put(h)
	return s
}

// SHA1String 返回 16 进制哈希字符串
func SHA1String(s string) string {
	h := sha1Pool.Get().(*buffer)
	s = h.HashString(s)
	sha1Pool.Put(h)
	return s
}
