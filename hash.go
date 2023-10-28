package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"sync"
)

var (
	md5Pool  sync.Pool
	sha1Pool sync.Pool
)

func init() {
	md5Pool.New = func() interface{} {
		return &hashBuffer{
			hash: md5.New(),
			buf:  make([]byte, 0, md5.Size*2),
			sum:  make([]byte, md5.Size),
		}
	}
	sha1Pool.New = func() interface{} {
		return &hashBuffer{
			hash: sha1.New(),
			buf:  make([]byte, 0, sha1.Size*2),
			sum:  make([]byte, sha1.Size),
		}
	}
}

// MD5 返回 16 进制哈希字符串
func MD5(b []byte) string {
	h := md5Pool.Get().(*hashBuffer)
	s := h.Hash(b)
	md5Pool.Put(h)
	return s
}

// SHA1 返回 16 进制哈希字符串
func SHA1(b []byte) string {
	h := sha1Pool.Get().(*hashBuffer)
	s := h.Hash(b)
	sha1Pool.Put(h)
	return s
}

// MD5String 返回 16 进制哈希字符串
func MD5String(s string) string {
	h := md5Pool.Get().(*hashBuffer)
	s = h.HashString(s)
	md5Pool.Put(h)
	return s
}

// SHA1String 返回 16 进制哈希字符串
func SHA1String(s string) string {
	h := sha1Pool.Get().(*hashBuffer)
	s = h.HashString(s)
	sha1Pool.Put(h)
	return s
}

// 用于做hash运算的缓存
type hashBuffer struct {
	hash hash.Hash
	buf  []byte
	sum  []byte
}

func (h *hashBuffer) Hash(b []byte) string {
	h.hash.Reset()
	h.hash.Write(b)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}

func (h *hashBuffer) HashString(s string) string {
	h.buf = h.buf[:0]
	h.buf = append(h.buf, s...)
	h.hash.Reset()
	h.hash.Write(h.buf)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}
