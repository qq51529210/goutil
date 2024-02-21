package util

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
	"strings"
)

// hash 算法的名称
const (
	nameMD5    = "MD5"
	nameSHA1   = "SHA1"
	nameSHA256 = "SHA256"
	nameSHA384 = "SHA384"
	nameSHA512 = "SHA512"
)

// CryptoHash 获取  name 的 crypto.Hash ，默认是 MD5
func CryptoHash(name string) crypto.Hash {
	name = strings.TrimSpace(name)
	switch name {
	case nameSHA1:
		return crypto.SHA1
	case nameSHA256:
		return crypto.SHA256
	case nameSHA384:
		return crypto.SHA384
	case nameSHA512:
		return crypto.SHA512
	default:
		return crypto.MD5
	}
}

// NewHash 返回 name 的 hash ，默认是 MD5
func NewHash(name string) hash.Hash {
	name = strings.ToUpper(name)
	switch name {
	case nameMD5:
		return md5.New()
	case nameSHA1:
		return sha1.New()
	case nameSHA256:
		return sha256.New()
	case nameSHA384:
		return sha512.New384()
	case nameSHA512:
		return sha512.New()
	default:
		return md5.New()
	}
}
