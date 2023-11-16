package sdp

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	zero = "0"
)

var (
	// ErrFormat 表示 sdp 格式错误
	ErrFormat = errors.New("error format")
)

// Session 表示 sdp 的会话描述
type Session struct {
	// 版本
	// v=
	V string
	// 会话创建者
	// o=
	O *Origin
	// 名称
	// s=
	S string
	// 连接信息
	// c=
	C *Connection
	// 时间描述
	// t=
	T *Time
	// 媒体信息
	// m=
	M []*Media
}

// Init 初始化
func (s *Session) Init() {
	s.O = new(Origin)
	s.O.SessionID = zero
	s.O.SessionVersion = zero
	s.O.NetType = NetTypeIN
	s.O.AddrType = AddrTypeIP4
	s.C = new(Connection)
	s.C.NetType = NetTypeIN
	s.C.AddrType = AddrTypeIP4
	s.T = new(Time)
}

// ParseFrom 解析
func (s *Session) ParseFrom(reader io.Reader) error {
	scaner := bufio.NewScanner(reader)
	for scaner.Scan() {
		line := scaner.Text()
		if line == "" {
			continue
		}
		err := s.parse(scaner, line)
		if err != nil {
			return err
		}
	}
	if err := scaner.Err(); err != nil {
		return err
	}
	// 检查是否缺少必要的字段
	if s.O == nil || s.C == nil || s.T == nil || len(s.M) < 1 {
		return ErrFormat
	}
	//
	return nil
}

// parse 解析
func (s *Session) parse(scaner *bufio.Scanner, line string) error {
	// v=
	value := strings.TrimPrefix(line, "v=")
	if value != line {
		s.V = value
		return nil
	}
	// o=
	value = strings.TrimPrefix(line, "o=")
	if value != line {
		s.O = new(Origin)
		return s.O.Parse(value)
	}
	// s=
	value = strings.TrimPrefix(line, "s=")
	if value != line {
		s.S = value
		return nil
	}
	// t=
	value = strings.TrimPrefix(line, "t=")
	if value != line {
		s.T = new(Time)
		return s.T.Parse(value)
	}
	// c=
	value = strings.TrimPrefix(line, "c=")
	if value != line {
		s.C = new(Connection)
		return s.C.Parse(value)
	}
	// m=
	value = strings.TrimPrefix(line, "m=")
	if value != line {
		return s.parseM(scaner, value)
	}
	//
	return nil
}

// parseM 解析 m= 和它的子项
func (s *Session) parseM(scaner *bufio.Scanner, line string) error {
	m := new(Media)
	if err := m.Parse(line); err != nil {
		return err
	}
	m.A = make(map[string][]string)
	s.M = append(s.M, m)
	// 其他
	var value string
	for scaner.Scan() {
		line := scaner.Text()
		if line == "" {
			break
		}
		// a=
		value = strings.TrimPrefix(line, "a=")
		if value != line {
			m.parseA(value)
			continue
		}
		// c=
		value = strings.TrimPrefix(line, "c=")
		if value != line {
			m.C = new(Connection)
			if err := m.C.Parse(value); err != nil {
				return err
			}
			continue
		}
		// y=
		value = strings.TrimPrefix(line, "y=")
		if value != line {
			m.Y = value
			return nil
		}
		// f=
		value = strings.TrimPrefix(line, "f=")
		if value != line {
			m.F = value
			return nil
		}
		return s.parse(scaner, line)
	}
	//
	return scaner.Err()
}

// FormatTo 格式化
func (s *Session) FormatTo(buf *bytes.Buffer) {
	// v=
	fmt.Fprint(buf, "v=0\r\n")
	// o=
	if s.O != nil {
		s.O.FormatTo(buf)
	}
	// s=
	if s.S != "" {
		fmt.Fprintf(buf, "s=%s\r\n", s.S)
	}
	// c=
	if s.C != nil {
		s.C.FormatTo(buf)
	}
	// t=
	if s.T != nil {
		s.T.FormatTo(buf)
	}
	// m=
	for _, m := range s.M {
		m.FormatTo(buf)
	}
}