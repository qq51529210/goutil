package sip

import (
	"bytes"
	"fmt"
	gostrings "goutil/strings"
	"io"
	"strconv"
	"strings"
)

// URI 表示 <scheme:name@domain>
type URI struct {
	Scheme string
	Name   string
	Domain string
}

// Reset 重置
func (m *URI) Reset() {
	m.Scheme = ""
	m.Name = ""
	m.Domain = ""
}

// Dec 解析
func (m *URI) Dec(line string) bool {
	line = gostrings.TrimByte(line, gostrings.CharLessThan, gostrings.CharGreaterThan)
	// 有这种格式的
	if line != "*" {
		// scheme:
		m.Scheme, line = gostrings.Split(line, gostrings.CharColon)
		if m.Scheme == "" {
			return false
		}
		// name@domain
		m.Name, m.Domain = gostrings.Split(line, gostrings.CharAt)
		if m.Name == "" || m.Domain == "" {
			return false
		}
	}
	//
	return true
}

// Enc 格式化
func (m *URI) Enc(w *bytes.Buffer) {
	w.WriteByte(gostrings.CharLessThan)
	w.WriteString(m.Scheme)
	w.WriteByte(gostrings.CharColon)
	w.WriteString(m.Name)
	w.WriteByte(gostrings.CharAt)
	w.WriteString(m.Domain)
	w.WriteByte(gostrings.CharGreaterThan)
}

// Address 表示 name<uri>;tag
type Address struct {
	Name string
	URI  URI
	Tag  string
}

// Reset 重置
func (m *Address) Reset() {
	m.Name = ""
	m.Tag = ""
	m.URI.Reset()
}

// Enc 格式化
func (m *Address) Enc(w *bytes.Buffer, key string) {
	// From: / To:
	w.WriteString(key)
	// name
	if m.Name != "" {
		w.WriteString(m.Name)
	}
	// uri
	m.URI.Enc(w)
	// tag
	if m.Tag != "" {
		w.WriteByte(gostrings.CharSemicolon)
		w.WriteString("tag=")
		w.WriteString(m.Tag)
	}
	// crlf
	w.WriteString(gostrings.CRLF)
}

// Dec 解析
func (m *Address) Dec(line string) bool {
	var prefix, suffix string
	// name
	m.Name, suffix = gostrings.Split(line, gostrings.CharLessThan)
	// uri>;tag=
	prefix, suffix = gostrings.Split(suffix, gostrings.CharSemicolon)
	// uri>
	if !m.URI.Dec(prefix) {
		return false
	}
	// 找 tag=
	for suffix != "" {
		prefix, suffix = gostrings.Split(suffix, gostrings.CharSemicolon)
		if s := strings.TrimPrefix(prefix, "tag="); s != prefix {
			m.Tag = s
			break
		}
	}
	return true
}

// CSeq 表示 sn method
type CSeq struct {
	SN     string
	Method string
}

// Reset 重置
func (m *CSeq) Reset() {
	m.SN = ""
	m.Method = ""
}

// Dec 解析
func (m *CSeq) Dec(line string) bool {
	m.SN, m.Method = gostrings.Split(line, gostrings.CharSpace)
	return m.SN != "" && m.Method != ""
}

// Enc 格式化
func (m *CSeq) Enc(w *bytes.Buffer) {
	w.WriteString("CSeq: ")
	w.WriteString(m.SN)
	w.WriteByte(gostrings.CharSpace)
	w.WriteString(m.Method)
	w.WriteString(gostrings.CRLF)
}

// Via 表示 proto address;rport=x;branch=x
type Via struct {
	// SIP/2.0/UDP
	Proto    string
	Address  string
	Branch   string
	RPort    string
	Received string
	rport    bool
}

func (m *Via) HasRPort() bool {
	return m.rport
}

// Dec 解析
func (m *Via) Dec(line string) bool {
	var prefix, suffix string
	// proto
	m.Proto, suffix = gostrings.Split(line, gostrings.CharSpace)
	if m.Proto == "" {
		return false
	}
	// address
	m.Address, suffix = gostrings.Split(suffix, gostrings.CharSemicolon)
	if m.Address == "" {
		return false
	}
	for suffix != "" {
		prefix, suffix = gostrings.Split(suffix, gostrings.CharSemicolon)
		// branch
		s := strings.TrimPrefix(prefix, "branch=")
		if s != prefix {
			m.Branch = s
			continue
		}
		// rport
		s = strings.TrimPrefix(prefix, "rport=")
		if s != prefix {
			m.RPort = s
			m.rport = true
			continue
		}
		// received
		s = strings.TrimPrefix(prefix, "received=")
		if s != prefix {
			m.Received = s
			continue
		}
	}
	return true
}

// Enc 格式化到 w
func (m *Via) Enc(w *bytes.Buffer) {
	w.WriteString("Via: ")
	// proto address
	w.WriteString(m.Proto)
	w.WriteByte(gostrings.CharSpace)
	w.WriteString(m.Address)
	// rport
	w.WriteString(";rport")
	if m.RPort != "" {
		w.WriteByte(gostrings.CharEqual)
		w.WriteString(m.RPort)
	}
	// received
	if m.Received != "" {
		w.WriteString(";received=")
		w.WriteString(m.Received)
	}
	// branch
	if m.Branch != "" {
		w.WriteString(";branch=")
		w.WriteString(m.Branch)
	}
	// crlf
	w.WriteString(gostrings.CRLF)
}

// Header 表示消息的一些必需的头字段
type Header struct {
	Via           []*Via
	From          Address
	To            Address
	CallID        string
	CSeq          CSeq
	Contact       URI
	MaxForwards   string
	Expires       string
	ContentType   string
	UserAgent     string
	others        map[string]string
	contentLength int64
}

func (m *Header) errorLine(line string) error {
	return fmt.Errorf("parse header error line: %s", line)
}

// Dec 解析，max 是
func (m *Header) Dec(r Reader, max int) error {
	m.others = make(map[string]string)
	var from, to, cseq, contentLength bool
	for {
		// 读取一行数据
		line, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		// 最后一个空行
		if line == "" {
			break
		}
		// 是否超出最大字节
		max = max - len(line) - 2
		if max < 0 {
			return ErrLargeMessage
		}
		// key: value
		key, value := gostrings.Split(line, gostrings.CharColon)
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		// 大写
		uKey := strings.ToUpper(key)
		// 挑选出必要的头
		switch uKey {
		case "VIA":
			var via Via
			if !via.Dec(value) {
				return m.errorLine(line)
			}
			m.Via = append(m.Via, &via)
		case "FROM":
			if !m.From.Dec(value) {
				return m.errorLine(line)
			}
			from = true
		case "TO":
			if !m.To.Dec(value) {
				return m.errorLine(line)
			}
			to = true
		case "CALL-ID":
			m.CallID = strings.TrimSpace(value)
			if m.CallID == "" {
				return m.errorLine(line)
			}
		case "CSEQ":
			if !m.CSeq.Dec(value) {
				return m.errorLine(line)
			}
			cseq = true
		case "CONTACT":
			if !m.Contact.Dec(value) {
				return m.errorLine(line)
			}
		case "MAX-FORWARDS":
			m.MaxForwards = value
		case "EXPIRES":
			m.Expires = value
		case "CONTENT-TYPE":
			m.ContentType = value
		case "CONTENT-LENGTH":
			m.contentLength, err = strconv.ParseInt(value, 10, 64)
			if err != nil || m.contentLength < 0 {
				return m.errorLine(line)
			}
			if m.contentLength > int64(max) {
				return ErrLargeMessage
			}
			contentLength = true
		case "USER-AGENT":
			m.UserAgent = value
		default:
			m.others[uKey] = value
		}
	}
	// 检查必须的字段
	if len(m.Via) < 1 {
		return errMissHeaderVia
	}
	if !from {
		return errMissHeaderFrom
	}
	if !to {
		return errMissHeaderTo
	}
	if !cseq {
		return errMissHeaderCSeq
	}
	if m.CallID == "" {
		return errMissHeaderCallID
	}
	if !contentLength {
		return errMissHeaderContentLength
	}
	return nil
}

// Enc 格式化
func (m *Header) Enc(w *bytes.Buffer) {
	// Via
	for i := 0; i < len(m.Via); i++ {
		m.Via[i].Enc(w)
	}
	// From
	m.From.Enc(w, "From: ")
	// To
	m.To.Enc(w, "To: ")
	// Call-ID
	w.WriteString("Call-ID: ")
	w.WriteString(m.CallID)
	w.WriteString(gostrings.CRLF)
	// CSeq
	m.CSeq.Enc(w)
	// Contact
	if m.Contact.Scheme != "" && m.Contact.Name != "" && m.Contact.Domain != "" {
		w.WriteString("Contact: ")
		m.Contact.Enc(w)
		w.WriteString(gostrings.CRLF)
	}
	// Expires
	if m.Expires != "" {
		w.WriteString("Expires: ")
		w.WriteString(m.Expires)
		w.WriteString(gostrings.CRLF)
	}
	// Max-Forwards
	if m.MaxForwards != "" {
		w.WriteString("Max-Forwards: ")
		w.WriteString(m.MaxForwards)
		w.WriteString(gostrings.CRLF)
	}
	// Content-Type
	if m.ContentType != "" {
		w.WriteString("Content-Type: ")
		w.WriteString(m.ContentType)
		w.WriteString(gostrings.CRLF)
	}
	// Others
	for k, v := range m.others {
		w.WriteString(k)
		w.WriteString(": ")
		w.WriteString(v)
		w.WriteString(gostrings.CRLF)
	}
	// User-Agent
	if m.UserAgent != "" {
		w.WriteString("User-Agent: ")
		w.WriteString(m.UserAgent)
		w.WriteString(gostrings.CRLF)
	}
	// Content-Length
	fmt.Fprintf(w, "Content-Length: %d\r\n", m.contentLength)
}

// Reset 重置
func (m *Header) Reset() {
	m.Via = m.Via[:0]
	m.From.Reset()
	m.To.Reset()
	m.CallID = ""
	m.CSeq.Reset()
	m.Contact.Reset()
	m.MaxForwards = ""
	m.Expires = ""
	m.ContentType = ""
	m.UserAgent = ""
	m.others = map[string]string{}
}

// Has 是否存在
func (m *Header) Has(key string) (ok bool) {
	if m.others != nil {
		_, ok = m.others[key]
	}
	return
}

// Get 返回，key 需要时大写，因为底层解析的时候，全部转化成大写了
func (m *Header) Get(key string) (value string) {
	if m.others != nil {
		value = m.others[key]
	}
	return
}

// Set 设置
func (m *Header) Set(key, value string) {
	if m.others == nil {
		m.others = make(map[string]string)
	}
	m.others[key] = value
}

// Del 删除
func (m *Header) Del(key string) {
	if m.others != nil {
		delete(m.others, key)
	}
}

// ResetOther 重置 others
func (m *Header) ResetOther() {
	m.others = make(map[string]string)
}

// KeepBasic 重置 contact、contentType、useragent、other
func (m *Header) KeepBasic() {
	m.Contact.Reset()
	m.ContentType = ""
	m.UserAgent = ""
	m.others = make(map[string]string)
}

// Message 表示 start line + header + body 的结构
type Message struct {
	StartLine [3]string
	Header    Header
	Body      bytes.Buffer
	isReq     bool
}

// txKey 返回事务的 key
func (m *Message) txKey() string {
	return m.Header.CSeq.Method + m.Header.CallID + m.Header.Via[0].Branch
}

// decStartLine 解析 start line ，返回剩余的 max
func (m *Message) decStartLine(reader Reader, max int) (int, error) {
	// 读取一行
	line, err := reader.ReadLine()
	if err != nil {
		return max, err
	}
	max = max - len(line) - 2
	// 数据太大，返回错误
	if max < 0 {
		return max, ErrLargeMessage
	}
	// 解析
	if !m.decStartLine2(line) {
		return max, fmt.Errorf("error start line %s", line)
	}
	return max, nil
}

// Dec 从 r 中读取并解析一个完整的 Message
// 如果读取的字节数大于 max 返回错误
func (m *Message) Dec(r Reader, max int) (err error) {
	// start line
	max, err = m.decStartLine(r, max)
	if err != nil {
		return
	}
	// header
	if err = m.Header.Dec(r, max); err != nil {
		return
	}
	// body
	if m.Header.contentLength > 0 {
		_, err = io.CopyN(&m.Body, r, m.Header.contentLength)
	}
	return
}

// Enc 格式化 header 和 body（如果 body 不为空）到 w 中
// Content-Length 字段是根据 body 的大小自动添加的
func (m *Message) Enc(w *bytes.Buffer) {
	// start line
	w.WriteString(m.StartLine[0])
	w.WriteByte(gostrings.CharSpace)
	w.WriteString(m.StartLine[1])
	w.WriteByte(gostrings.CharSpace)
	w.WriteString(m.StartLine[2])
	w.WriteString(gostrings.CRLF)
	// header
	m.Header.contentLength = int64(m.Body.Len())
	m.Header.Enc(w)
	// header 和 body 的空行
	w.WriteString(gostrings.CRLF)
	// body
	if m.Header.contentLength > 0 {
		w.Write(m.Body.Bytes())
	}
}

// decStartLine2 解析 start line
func (m *Message) decStartLine2(line string) bool {
	m.StartLine[0], line = gostrings.Split(strings.TrimSpace(line), gostrings.CharSpace)
	if m.StartLine[0] == "" {
		return false
	}
	m.StartLine[1], m.StartLine[2] = gostrings.Split(line, gostrings.CharSpace)
	if m.StartLine[1] == "" || m.StartLine[2] == "" {
		return false
	}
	if m.StartLine[2] == SIPVersion {
		m.isReq = true
		return true
	}
	if m.StartLine[0] != SIPVersion {
		return false
	}
	//
	return true
}

// String 返回格式化后的字符串。
func (m *Message) String() string {
	var str bytes.Buffer
	m.Enc(&str)
	return str.String()
}

// 实现 io.WriterTo 接口
func (m *Message) WriteTo(w io.Writer) (int64, error) {
	var str bytes.Buffer
	m.Enc(&str)
	return io.Copy(w, &str)
}

// KeepBasic 重置 contact、contentType、useragent、other、body
func (m *Message) KeepBasic() {
	m.Header.KeepBasic()
	m.Body.Reset()
}

// SetResponseStartLine 设置 start line
func (m *Message) SetResponseStartLine(status, phrase string) {
	m.StartLine[0] = SIPVersion
	m.StartLine[1] = string(status)
	m.StartLine[2] = phrase
	if m.StartLine[2] == "" {
		m.StartLine[2] = StatusPhrase(status)
	}
}
