package hash

import (
	"encoding/hex"
	"hash"
)

// buffer 用于做hash运算的缓存
type buffer struct {
	hash hash.Hash
	buf  []byte
	sum  []byte
}

func (h *buffer) Hash(b []byte) string {
	h.hash.Reset()
	h.hash.Write(b)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}

func (h *buffer) HashString(s string) string {
	h.buf = h.buf[:0]
	h.buf = append(h.buf, s...)
	h.hash.Reset()
	h.hash.Write(h.buf)
	h.hash.Sum(h.sum[:0])
	h.buf = h.buf[:h.hash.Size()*2]
	hex.Encode(h.buf, h.sum)
	return string(h.buf)
}
