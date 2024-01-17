package soap

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/xml"
	"math/rand"
	"time"
)

const (
	// SecurityNamespace 命名空间
	SecurityNamespace = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd"
)

var (
	// NamespaceAttr 命名空间属性
	NamespaceAttr = NewSecurityNamespaceAttr()
)

var (
	_rand = rand.New(rand.NewSource(time.Now().UnixMilli()))
)

// Security 表示
type Security struct {
	XMLName       xml.Name `xml:"wsse:Security"`
	UsernameToken UsernameToken
}

// NewSecurity 返回初始化好的
func NewSecurity(username, password string) *Security {
	s := new(Security)
	s.Init(username, password)
	return s
}

// Init 初始化
func (s *Security) Init(username, password string) {
	//
	s.UsernameToken.Username = username
	s.UsernameToken.Password.Type = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-username-token-profile-1.0#PasswordDigest"
	s.UsernameToken.Created = time.Now().UTC().Format(time.RFC3339Nano)
	s.UsernameToken.Nonce.Type = "http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary"
	//
	nonce := make([]byte, 32)
	_rand.Read(nonce)
	//
	var buf bytes.Buffer
	buf.Write(nonce)
	buf.WriteString(s.UsernameToken.Created)
	buf.WriteString(password)
	//
	hash := sha1.New()
	hash.Write(buf.Bytes())
	//
	s.UsernameToken.Password.Password = base64.StdEncoding.EncodeToString(hash.Sum(nil))
	s.UsernameToken.Nonce.Nonce = base64.StdEncoding.EncodeToString(nonce)
}

// UsernameToken 表示 Security 的 UsernameToken 字段
type UsernameToken struct {
	XMLName  xml.Name `xml:"wsse:UsernameToken"`
	Username string   `xml:"wsse:Username"`
	Password Password `xml:"wsse:Password"`
	Nonce    Nonce    `xml:"wsse:Nonce"`
	Created  string   `xml:"http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd Created"`
}

// Password 表示 UsernameToken 的 Password 字段
type Password struct {
	Type     string `xml:"Type,attr"`
	Password string `xml:",chardata"`
}

// Nonce 表示 UsernameToken 的 Nonce 字段
type Nonce struct {
	Type  string `xml:"EncodingType,attr"`
	Nonce string `xml:",chardata"`
}

// NewSecurityNamespaceAttr 返回命名空间属性
func NewSecurityNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:wsse",
		},
		Value: SecurityNamespace,
	}
}
