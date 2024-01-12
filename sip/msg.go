package sip

import (
	"bytes"
	"fmt"
	gostrings "goutil/strings"
	"goutil/uid"
	"io"
	"strconv"
	"strings"
)

// URI 表示 <sip:name@domain>
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
	line = gostrings.TrimByte(line, '<', '>')
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
func (m *URI) Enc(w Writer) (err error) {
	_, err = fmt.Fprintf(w, "<%s:%s@%s>", m.Scheme, m.Name, m.Domain)
	return
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
func (m *Address) Enc(w Writer, prefix string) (err error) {
	if _, err = w.WriteString(prefix); err != nil {
		return
	}
	// name
	if m.Name != "" {
		if _, err = w.WriteString(m.Name); err != nil {
			return
		}
	}
	// uri
	err = m.URI.Enc(w)
	if err != nil {
		return
	}
	// tag
	if m.Tag != "" {
		if _, err = fmt.Fprintf(w, ";tag=%s", m.Tag); err != nil {
			return
		}
	}
	// crlf
	_, err = w.WriteString(CRLF)
	return
}

// Dec 解析
func (m *Address) Dec(line string) bool {
	prefix, suffix := gostrings.Split(line, gostrings.CharSemicolon)
	if prefix == "" {
		return false
	}
	// name
	m.Name, prefix = gostrings.Split(prefix, '<')
	// uri
	if !m.URI.Dec(prefix) {
		return false
	}
	for suffix != "" {
		prefix, suffix = gostrings.Split(suffix, gostrings.CharSemicolon)
		// tag
		s := strings.TrimPrefix(prefix, "tag=")
		if s != prefix {
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
	if m.SN == "" || m.Method == "" {
		return false
	}
	return true
}

// Enc 格式化
func (m *CSeq) Enc(w Writer) (err error) {
	_, err = fmt.Fprintf(w, "CSeq: %s %s\r\n", m.SN, m.Method)
	return
}

// Via 表示 proto address;rport=x;branch=x
type Via struct {
	// SIP/2.0/UDP
	Proto    string
	Address  string
	Branch   string
	RProt    string
	Received string
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
			m.RProt = s
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
func (m *Via) Enc(w Writer) (err error) {
	// proto address
	if _, err = fmt.Fprintf(w, "Via: %s %s", m.Proto, m.Address); err != nil {
		return
	}
	// rport
	if _, err = w.WriteString(";rport"); err != nil {
		return
	}
	if m.RProt != "" {
		if _, err = fmt.Fprintf(w, "=%s", m.RProt); err != nil {
			return
		}
	}
	// received
	if m.Received != "" {
		if _, err = fmt.Fprintf(w, ";received=%s", m.Received); err != nil {
			return
		}
	}
	// branch
	if m.Branch != "" {
		if _, err = fmt.Fprintf(w, ";branch=%s", m.Branch); err != nil {
			return
		}
	}
	// crlf
	_, err = w.WriteString(CRLF)
	return
}

// header 表示消息的一些必需的头字段
type header struct {
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

func (m *header) errorLine(line string) error {
	return fmt.Errorf("error header %s", line)
}

// Dec 解析
func (m *header) Dec(r Reader, max int) (int, error) {
	m.others = make(map[string]string)
	var from, to, cseq, contentLength bool
	for {
		// 读取一行数据
		line, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return max, err
		}
		// 最后一个空行
		if line == "" {
			break
		}
		// 是否超出最大字节
		max = max - len(line) - 2
		if max < 0 {
			return max, errLargeMessage
		}
		key, value := gostrings.Split(line, gostrings.CharColon)
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		uKey := strings.ToUpper(key)
		// 挑选出必要的头
		switch uKey {
		case "VIA":
			var via Via
			if !via.Dec(value) {
				return max, m.errorLine(line)
			}
			m.Via = append(m.Via, &via)
		case "FROM":
			if !m.From.Dec(value) {
				return max, m.errorLine(line)
			}
			from = true
		case "TO":
			if !m.To.Dec(value) {
				return max, m.errorLine(line)
			}
			to = true
		case "CALL-ID":
			m.CallID = strings.TrimSpace(value)
			if m.CallID == "" {
				return max, m.errorLine(line)
			}
		case "CSEQ":
			if !m.CSeq.Dec(value) {
				return max, m.errorLine(line)
			}
			cseq = true
		case "CONTACT":
			if !m.Contact.Dec(value) {
				return max, m.errorLine(line)
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
				return max, m.errorLine(line)
			}
			if m.contentLength > int64(max) {
				return max, errLargeMessage
			}
			contentLength = true
		case "USER-AGENT":
			m.UserAgent = value
		default:
			m.others[key] = value
		}
	}
	// 检查必须的字段
	if len(m.Via) < 1 {
		return max, errMissHeaderVia
	}
	if !from {
		return max, errMissHeaderFrom
	}
	if !to {
		return max, errMissHeaderTo
	}
	if !cseq {
		return max, errMissHeaderCSeq
	}
	if m.CallID == "" {
		return max, errMissHeaderCallID
	}
	if !contentLength {
		return max, errMissHeaderContentLength
	}
	return max, nil
}

// Enc 格式化
func (m *header) Enc(w Writer) (err error) {
	// Via
	for i := 0; i < len(m.Via); i++ {
		if err = m.Via[i].Enc(w); err != nil {
			return
		}
	}
	// From
	if err = m.From.Enc(w, "From: "); err != nil {
		return
	}
	// To
	if err = m.To.Enc(w, "To: "); err != nil {
		return
	}
	// Call-ID
	if _, err = fmt.Fprintf(w, "Call-ID: %s\r\n", m.CallID); err != nil {
		return
	}
	// CSeq
	if _, err = fmt.Fprintf(w, "CSeq: %s %s\r\n", m.CSeq.SN, m.CSeq.Method); err != nil {
		return
	}
	// Contact
	if m.Contact.Scheme != "" && m.Contact.Name != "" && m.Contact.Domain != "" {
		if _, err = fmt.Fprintf(w, "Contact: <%s:%s@%s>\r\n", m.Contact.Scheme, m.Contact.Name, m.Contact.Domain); err != nil {
			return
		}
	}
	// Expires
	if m.Expires != "" {
		if _, err = fmt.Fprintf(w, "Expires: %s\r\n", m.Expires); err != nil {
			return
		}
	}
	// Max-Forwards
	if m.MaxForwards != "" {
		if _, err = fmt.Fprintf(w, "Max-Forwards: %s\r\n", m.MaxForwards); err != nil {
			return
		}
	}
	// Content-Type
	if m.ContentType != "" {
		if _, err = fmt.Fprintf(w, "Content-Type: %s\r\n", m.ContentType); err != nil {
			return
		}
	}
	// Others
	for k, v := range m.others {
		if _, err = fmt.Fprintf(w, "%s: %s\r\n", k, v); err != nil {
			return
		}
	}
	// User-Agent
	if m.UserAgent != "" {
		if _, err = fmt.Fprintf(w, "User-Agent: %s\r\n", m.UserAgent); err != nil {
			return
		}
	}
	// Content-Length
	_, err = fmt.Fprintf(w, "Content-Length: %d\r\n", m.contentLength)
	return
}

// Reset 重置
func (m *header) Reset() {
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
func (m *header) Has(key string) (ok bool) {
	if m.others != nil {
		_, ok = m.others[key]
	}
	return
}

// Get 返回
func (m *header) Get(key string) (value string) {
	if m.others != nil {
		value = m.others[key]
	}
	return
}

// Set 设置
func (m *header) Set(key, value string) {
	if m.others == nil {
		m.others = make(map[string]string)
	}
	m.others[key] = value
}

// Del 删除
func (m *header) Del(key string) {
	if m.others != nil {
		delete(m.others, key)
	}
}

// ResetOther 重置 others
func (m *header) ResetOther() {
	m.others = make(map[string]string)
}

// KeepBasic 重置 contact、contentType、useragent、other
func (m *header) KeepBasic() {
	m.Contact.Reset()
	m.ContentType = ""
	m.UserAgent = ""
	m.others = make(map[string]string)
}

// message 表示 start line + header + body 的结构
type message struct {
	StartLine [3]string
	Header    header
	Body      bytes.Buffer
	isReq     bool
}

// txKey 返回事务的 key
func (m *message) txKey() string {
	return m.Header.CSeq.Method + m.Header.CallID + m.Header.Via[0].Branch
}

// Dec 从 r 中读取并解析一个完整的 message
// 如果读取的字节数大于 max 返回错误
func (m *message) Dec(r Reader, max int) (err error) {
	// start line
	max, err = m.DecStartLine(r, max)
	if err != nil {
		return
	}
	// header
	_, err = m.Header.Dec(r, max)
	if err != nil {
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
func (m *message) Enc(w Writer) (err error) {
	// start line
	fmt.Fprintf(w, "%s %s %s\r\n", m.StartLine[0], m.StartLine[1], m.StartLine[2])
	// header
	m.Header.contentLength = int64(m.Body.Len())
	if err = m.Header.Enc(w); err != nil {
		return
	}
	// header 和 body 的空行
	if _, err = w.WriteString(CRLF); err != nil {
		return
	}
	// body
	if m.Header.contentLength > 0 {
		_, err = w.Write(m.Body.Bytes())
	}
	return
}

// DecStartLine 解析 start line ，返回剩余的 max
func (m *message) DecStartLine(reader Reader, max int) (int, error) {
	// 读取一行
	line, err := reader.ReadLine()
	if err != nil {
		return max, err
	}
	max = max - len(line) - len(CRLF)
	// 数据太大，返回错误
	if max < 0 {
		return max, fmt.Errorf("large bytes %d of start line", len(line))
	}
	// 解析
	if !m.decStartLine(line) {
		return max, fmt.Errorf("error start line %s", line)
	}
	return max, nil
}

// decStartLine 解析 start line
func (m *message) decStartLine(line string) bool {
	// 一部分
	line = strings.TrimSpace(line)
	i := strings.Index(line, " ")
	if i < 0 {
		return false
	}
	m.StartLine[0] = strings.ToUpper(line[:i])
	// 二部分
	line = strings.TrimSpace(line[i+1:])
	i = strings.Index(line, " ")
	if i < 0 {
		return false
	}
	m.StartLine[1] = strings.ToUpper(line[:i])
	// 三部分
	m.StartLine[2] = strings.TrimSpace(line[i+1:])
	// 检查
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
func (m *message) String() string {
	var str strings.Builder
	m.Enc(&str)
	return str.String()
}

// Request 表示请求消息
type Request struct {
	*Server
	*message
	tx
}

// NewRequest 创建请求
func NewRequest() *Request {
	return &Request{message: new(message)}
}

// Response 用于发送响应消息，在处理请求消息的时候使用
// 直接修改的是 request 的数据哦
func (m *Request) Response(status string) {
	m.ResponseWith(status, StatusPhrase(status))
}

// ResponseWith 用于发送响应消息，在处理请求消息的时候使用
// 直接修改的是 request 的数据哦
func (m *Request) ResponseWith(status, phrase string) {
	// start line
	m.StartLine[0] = SIPVersion
	m.StartLine[1] = string(status)
	m.StartLine[2] = phrase
	if m.StartLine[2] == "" {
		m.StartLine[2] = StatusPhrase(status)
	}
	// to tag
	if m.Header.To.Tag == "" {
		m.Header.To.Tag = fmt.Sprintf("%d", uid.SnowflakeID())
	}
	if m.Header.Via[0].Received == "" {
		m.Header.Via[0].Received = m.RemoteIP()
	}
	if m.Header.Via[0].RProt == "" {
		m.Header.Via[0].RProt = m.RemotePort()
	}
	//
	m.Header.UserAgent = m.Server.UserAgent
	// 格式化
	buf := m.dataBuffer()
	buf.Reset()
	m.Enc(buf)
	// 日志
	m.Server.Logger.DebugfTrace(m.tx.TxKey(), "write response to %s %s\n%s", m.Network(), m.RemoteAddrString(), buf.String())
	// 立刻发送
	err := m.write(buf.Bytes())
	if err != nil {
		m.Server.Logger.ErrorDepthTrace(2, m.txKey(), err)
	}
}

// ResponseError 调用 Response
func (m *Request) ResponseError(err *ResponseError) {
	m.ResponseWith(err.Status, err.Phrase)
}

// NewResponse 根据 Request 创建 Response
func (m *Request) NewResponse(status, phrase string) *Response {
	r := new(Response)
	r.Server = m.Server
	r.tx = m.tx
	r.message = new(message)
	// start line
	r.StartLine[0] = SIPVersion
	r.StartLine[1] = string(status)
	r.StartLine[2] = phrase
	if r.StartLine[2] == "" {
		r.StartLine[2] = StatusPhrase(status)
	}
	// header
	for _, v := range m.Header.Via {
		vv := new(Via)
		*vv = *v
		r.Header.Via = append(r.Header.Via, vv)
	}
	if r.Header.Via[0].Received == "" {
		r.Header.Via[0].Received = m.RemoteIP()
	}
	if r.Header.Via[0].RProt == "" {
		r.Header.Via[0].RProt = m.RemotePort()
	}
	r.Header.From = m.Header.From
	r.Header.To = m.Header.To
	r.Header.To.Tag = fmt.Sprintf("%d", uid.SnowflakeID())
	r.Header.CallID = m.Header.CallID
	r.Header.CSeq = m.Header.CSeq
	r.Header.MaxForwards = m.Header.MaxForwards
	r.Header.Expires = m.Header.Expires
	r.Header.UserAgent = m.Server.UserAgent
	r.Header.others = make(map[string]string)
	//
	return r
}

// Response 表示响应消息
type Response struct {
	*Server
	*message
	tx
}

// IsStatus 返回是否与 code 相等，因为 StartLine[1] 不好记忆
func (m *Response) IsStatus(status string) bool {
	return m.StartLine[1] == status
}

// Error 返回当前的 status 和 pharse
func (m *Response) Error() *ResponseError {
	return NewResponseError(m.StartLine[1], m.StartLine[2], "")
}
