package sip

import (
	"io"
)

// reader 用于读取 crlf 的每一行数据
type reader struct {
	r io.Reader
	// 缓存
	buf []byte
	// 有效数据起始下标
	begin int
	// 有效数据终止下标
	end int
	// 已经解析的下标
	parsed int
}

// newReader 返回新的 reader
func newReader(r io.Reader, n int) *reader {
	if n < 1 {
		n = MaxMessageLen
	}
	return &reader{r: r, buf: make([]byte, n)}
}

func (r *reader) Raw() []byte {
	return r.buf[:r.end]
}

func (r *reader) Reset(reader io.Reader) {
	r.r = reader
	r.begin = 0
	r.end = 0
	r.parsed = 0
}

func (r *reader) ReadLine() (string, error) {
	for {
		// 尝试查找
		for r.parsed < r.end {
			// 找到回车
			if r.buf[r.parsed] == '\n' {
				i := r.parsed - 1
				// 找到换行
				if r.buf[i] == '\r' {
					line := string(r.buf[r.begin:i])
					r.parsed++
					r.begin = r.parsed
					r.checkEmpty()
					return line, nil
				}
			}
			r.parsed++
		}
		// 缓存中没有，那么要读取数据了
		if r.end == len(r.buf) {
			// 这一行数据太大了，指定的缓存装不下了
			if r.begin == 0 {
				return "", errLargeMessage
			}
			// 缓存向前移
			if r.begin > 0 {
				copy(r.buf, r.buf[r.begin:r.end])
				r.end -= r.begin
				r.parsed -= r.begin
				r.begin = 0
			}
		}
		// 继续读
		n, err := r.r.Read(r.buf[r.end:])
		if err != nil {
			return "", err
		}
		r.end += n
	}
}

func (r *reader) Read(b []byte) (int, error) {
	if r.begin == r.end {
		return r.r.Read(b)
	}
	n := copy(b, r.buf[r.begin:r.end])
	r.begin += n
	r.parsed += n
	r.checkEmpty()
	return n, nil
}

func (r *reader) checkEmpty() {
	if r.begin == r.end {
		r.begin = 0
		r.parsed = 0
		r.end = 0
	}
}
