package rand

import (
	"math/rand"
	"time"
)

var (
	_random random
)

func init() {
	// 种子
	_random.r = rand.New(rand.NewSource(time.Now().Unix()))
	// 数字
	for i := '0'; i <= '9'; i++ {
		_random.b[0] = append(_random.b[0], byte(i))
		_random.b[3] = append(_random.b[3], byte(i))
	}
	// 小写字母
	for i := 'a'; i <= 'z'; i++ {
		_random.b[1] = append(_random.b[1], byte(i))
		_random.b[3] = append(_random.b[3], byte(i))
	}
	// 大写字母
	for i := 'A'; i <= 'Z'; i++ {
		_random.b[2] = append(_random.b[2], byte(i))
		_random.b[3] = append(_random.b[3], byte(i))
	}
}

type random struct {
	// 0:数字
	// 1:小写字母
	// 2:大写字母
	// 3:混合
	b [4][]byte
	// 随机数
	r *rand.Rand
}

func (r *random) rand(n, bi int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = r.b[bi][r.r.Intn(len(r.b[bi]))]
	}
	return string(b)
}

// Number 返回 n 长度的随机数字
func Number(n int) string {
	return _random.rand(n, 0)
}

// Lower 返回 n 长度的随机小写字母
func Lower(n int) string {
	return _random.rand(n, 1)
}

// Upper 返回 n 长度的随机大写字母
func Upper(n int) string {
	return _random.rand(n, 2)
}

// String 返回 n 长度的随机字母混合数字
func String(n int) string {
	return _random.rand(n, 3)
}

// Int 返回随机的整数
func Int() int {
	return _random.r.Int()
}

// Intn 返回随机的整数
func Intn(n int) int {
	return _random.r.Intn(n)
}
