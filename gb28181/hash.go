package gb28181

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

type HashName string

// hash 算法的名称
const (
	HashMD5    HashName = "MD5"
	HashSHA1   HashName = "SHA1"
	HashSHA256 HashName = "SHA256"
	HashSHA384 HashName = "SHA384"
	HashSHA512 HashName = "SHA512"
)

// CryptoHash 获取  name 的 crypto.Hash ，默认是 MD5
func CryptoHash(name HashName) crypto.Hash {
	switch name {
	case HashSHA1:
		return crypto.SHA1
	case HashSHA256:
		return crypto.SHA256
	case HashSHA384:
		return crypto.SHA384
	case HashSHA512:
		return crypto.SHA512
	default:
		return crypto.MD5
	}
}

// NewHash 返回 name 的 hash ，默认是 MD5
func NewHash(name HashName) hash.Hash {
	switch name {
	case HashSHA1:
		return sha1.New()
	case HashSHA256:
		return sha256.New()
	case HashSHA384:
		return sha512.New384()
	case HashSHA512:
		return sha512.New()
	default:
		return md5.New()
	}
}
