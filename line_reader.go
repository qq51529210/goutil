package util

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

var (
	crlf = []byte("\r\n")
)

// LineReader 用于读取 lf 的每一行数据
type LineReader struct {
	// 数据源
	scaner *bufio.Scanner
}

// ReadLine读取一行数据
func (r *LineReader) ReadLine() (string, error) {
	var line strings.Builder
	for r.scaner.Scan() {
		b := r.scaner.Bytes()
		if len(b) < 1 {
			continue
		}
		line.Write(b)
		break
	}
	if err := r.scaner.Err(); err != nil && err != io.EOF {
		return "", err
	}
	return line.String(), nil
}

// NewLineReader 返回新的 LineReader
func NewLineReader(reader io.Reader) *LineReader {
	r := new(LineReader)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	r.scaner = scanner
	return r
}

// CRLFReader 用于读取 crlf 的每一行数据
type CRLFReader struct {
	// 数据源
	scaner *bufio.Scanner
}

// ReadLine读取一行数据
func (r *CRLFReader) ReadLine() (string, error) {
	var line strings.Builder
	for r.scaner.Scan() {
		b := r.scaner.Bytes()
		if len(b) < 1 {
			continue
		}
		line.Write(b)
		break
	}
	if err := r.scaner.Err(); err != nil && err != io.EOF {
		return "", err
	}
	return line.String(), nil
}

// NewCRLFReader 返回新的 CRLFReader
func NewCRLFReader(reader io.Reader) *CRLFReader {
	r := new(CRLFReader)
	scanner := bufio.NewScanner(reader)
	scanner.Split(r.scanLines)
	r.scaner = scanner
	return r
}

// scanLines 是 bufio.Scanner 的 Split 函数
func (r *CRLFReader) scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// 查找数据
	if len(data) > 2 {
		if i := bytes.Index(data, crlf); i >= 0 {
			return i + 2, data[:i], nil
		}
	}
	// 结尾
	if atEOF && len(data) > 0 {
		if i := bytes.Index(data, crlf); i >= 0 {
			return i + 2, data[:i], nil
		}
		return len(data), data, nil
	}
	// 返回
	return 0, nil, nil
}
