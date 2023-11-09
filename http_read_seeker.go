package util

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	httpReadSeekCloserHeaderLen = 1024
)

var (
	errHTTPReadSeekCloserSeekOffset = errors.New("error seek offset")
)

// HTTPReadSeekCloser 实现 io.ReadSeekCloseer，服务必须支持 Range 哦
type HTTPReadSeekCloser struct {
	req *http.Request
	res *http.Response
	// 请求头的 ETag
	eTag string
	// 当前位置
	offset int64
	// 数据总量
	total int64
}

// Close 实现 io.Closer 接口，主要用于关闭 response.body
func (r *HTTPReadSeekCloser) Close() error {
	res := r.res
	r.res = nil
	if res != nil {
		res.Body.Close()
	}
	return nil
}

// Read 实现 io.Reader 接口
func (r *HTTPReadSeekCloser) Read(buf []byte) (int, error) {
	//
	if r.offset >= r.total {
		return 0, io.EOF
	}
	//
	n, err := r.res.Body.Read(buf)
	if err != nil {
		return n, err
	}
	r.offset += int64(n)
	//
	return n, nil
}

// Seek 实现 io.Seeker 接口
func (r *HTTPReadSeekCloser) Seek(offset int64, whence int) (int64, error) {
	// 如果 offset 计算后的距离比 httpReadSeekerHeaderLen 小，
	// 那么没必要再去发起一个请求，直接读就可以了。
	switch whence {
	case io.SeekCurrent:
		// 不变
		if offset == 0 {
			return r.offset, nil
		}
		// 偏移
		offset += r.offset
	case io.SeekEnd:
		// 整数
		if offset == 0 {
			return r.total, nil
		}
		// 偏移
		offset += r.total
	case io.SeekStart:
		// 不变
		if offset == r.offset {
			return r.offset, nil
		}
	default:
		// 错误的
		return r.offset, fmt.Errorf("seek error whence %d", whence)
	}
	// 负数
	if offset < 0 {
		return r.offset, errHTTPReadSeekCloserSeekOffset
	}
	// 结尾
	if offset >= r.total {
		r.offset = offset
		return offset, nil
	}
	// 偏移小，读取
	if offset > r.offset {
		n := offset - r.offset
		if n < httpReadSeekCloserHeaderLen {
			_, err := io.CopyN(io.Discard, r.res.Body, n)
			if err != nil {
				return r.offset, err
			}
			r.offset = offset
			//
			return r.offset, nil
		}
	}
	// 请求
	err := r.seek(offset)
	if err != nil {
		return r.offset, err
	}
	return r.offset, nil
}

func (r *HTTPReadSeekCloser) seek(offset int64) error {
	// 关闭上一个响应的 body
	res := r.res
	r.res = nil
	if res != nil {
		res.Body.Close()
	}
	// 重新设置
	if r.total < 1 {
		r.req.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	} else {
		r.req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", offset, r.total))
	}
	if r.eTag != "" {
		r.req.Header.Set("If-Range", r.eTag)
	}
	// 发起请求
	var err error
	res, err = http.DefaultClient.Do(r.req)
	if err != nil {
		return err
	}
	r.res = res
	// 状态码
	if r.res.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("error status code %d", r.res.StatusCode)
	}
	r.eTag = r.res.Header.Get("ETag")
	// 解析出 pos 和 total
	contentRange := r.res.Header.Get("Content-Range")
	if contentRange == "" {
		return errors.New("missing content-range header")
	}
	i := strings.LastIndexByte(contentRange, '/')
	if i < 1 {
		return nil
	}
	total, err := strconv.ParseInt(contentRange[i+1:], 10, 64)
	if err != nil {
		return fmt.Errorf("error content range %s", contentRange)
	}
	r.offset = offset
	r.total = total
	//
	return nil
}

// NewHTTPReadSeeker 返回新的 HTTPReadSeekCloser ，开始会请求一次
func NewHTTPReadSeeker(url string, offset int64) (*HTTPReadSeekCloser, error) {
	// 请求
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept-Ranges", "bytes")
	// 创建
	r := new(HTTPReadSeekCloser)
	r.req = req
	err = r.seek(offset)
	if err != nil {
		return nil, err
	}
	//
	return r, nil
}
