package model

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	gbu "goutil/gb28181/util"
	"goutil/sip"
	gs "goutil/strings"
	"strings"
)

// 消息头字段名称
const (
	StrDate            = "DATE"
	StrAuthorization   = "AUTHORIZATION"
	StrWWWAuthenticate = "WWW-AUTHENTICATE"
)

// WWW-Authenticate 算法的名称
const (
	StrDigest     = "Digest"
	StrCapability = "Capability"
	StrAsymmetric = "Asymmetric"
)

// WWW-Authenticate 子字段名称
const (
	StrRealm     = "realm"
	StrNonce     = "nonce"
	StrQOP       = "qop"
	StrUsername  = "username"
	StrURI       = "uri"
	StrResponse  = "response"
	StrAlgorithm = "algorithm"
	StrCNonce    = "cnonce"
	StrNC        = "nc"
)

// kvQuotationMark 解析 k="v" ，返回 k, v
func kvQuotationMark(line string) (string, string) {
	// k = "v"
	k, v := gs.KeyValue(line)
	// v 去 ""
	return k, gs.TrimByte(v, gs.CharQuotationMark, gs.CharQuotationMark)
}

// rsaGenNonce 返回 pri.sign(hash(rand)), pub.enc(pub.enc(rand)), rand
func rsaGenNonce(pri *rsa.PrivateKey, pub *rsa.PublicKey, hash crypto.Hash) ([]byte, []byte, []byte, error) {
	// 随机 c
	c := make([]byte, 32)
	rand.Reader.Read(c)
	// 公钥加密 c ，得到 b
	b, err := rsa.EncryptPKCS1v15(rand.Reader, pub, c)
	if err != nil {
		return nil, nil, nil, err
	}
	// 对 c 哈希，得到 d
	h := hash.New()
	h.Reset()
	h.Write(c)
	d := h.Sum(nil)
	// 使用 ser 私钥签名 d ，得到 a
	a, err := rsa.SignPKCS1v15(rand.Reader, pri, hash, d)
	if err != nil {
		return nil, nil, nil, err
	}
	//
	return a, b, c, nil
}

// rsaVerifyNonce 验证 nonce ，返回随机数 c
// 算法 pub.verify(hash(pri.dec(base64.dec(nonce2))), base64.dec(nonce1))
func rsaVerifyNonce(pri *rsa.PrivateKey, pub *rsa.PublicKey, hash crypto.Hash, nonce string) ([]byte, error) {
	// 得到 a 和 b
	p1, p2 := gs.Split(nonce, gs.CharAmpersand)
	a, err := base64.StdEncoding.DecodeString(p1)
	if err != nil {
		return nil, err
	}
	b, err := base64.StdEncoding.DecodeString(p2)
	if err != nil {
		return nil, err
	}
	// 私钥解密 b ，得到 c
	c, err := rsa.DecryptPKCS1v15(rand.Reader, pri, b)
	if err != nil {
		return nil, err
	}
	// 对 c 哈希，得到 d
	h := hash.New()
	h.Reset()
	h.Write(c)
	d := h.Sum(nil)
	// 公钥验证 d 和 a
	return c, rsa.VerifyPKCS1v15(pub, hash, d, a)
}

// RegisterHeader 表示 REGISTER 消息的独有字段
type RegisterHeader struct {
	WWWAuthenticate *RegisterHeaderWWWAuthenticate
	Authorization   *RegisterHeaderAuthorization
	Expires         string
}

// Parse 解析
func (m *RegisterHeader) Parse(msg *sip.Message) bool {
	// Expires
	m.Expires = msg.Header.Expires
	// 先检查有没有 Authorization
	value := msg.Header.Get(StrAuthorization)
	if value != "" {
		m.Authorization = new(RegisterHeaderAuthorization)
		return m.Authorization.Parse(value)
	}
	// 没有 Authorization 在看看有没有 WWW-Authenticate
	value = msg.Header.Get(StrWWWAuthenticate)
	if value != "" {
		m.WWWAuthenticate = new(RegisterHeaderWWWAuthenticate)
		return m.WWWAuthenticate.Parse(value)
	}
	// 都没有，那就是普通认证的第一个请求了
	return true
}

// RegisterHeaderWWWAuthenticate 是服务端发送的消息，表示
// WWW-Authenticate: Digest/Asymmetric ...
type RegisterHeaderWWWAuthenticate struct {
	Digest     *RegisterHeaderWWWAuthenticateDigest
	Asymmetric *RegisterHeaderWWWAuthenticateAsymmetric
}

// Parse 解析
func (m *RegisterHeaderWWWAuthenticate) Parse(line string) bool {
	prefix, suffix := gs.Split(line, gs.CharSpace)
	// 类型，Digest/Asymmetric
	switch prefix {
	case StrDigest:
		m.Digest = new(RegisterHeaderWWWAuthenticateDigest)
		return m.Digest.Parse(suffix)
	case StrAsymmetric:
		m.Asymmetric = new(RegisterHeaderWWWAuthenticateAsymmetric)
		return m.Asymmetric.Parse(suffix)
	default:
		return false
	}
}

// RegisterHeaderWWWAuthenticateDigest 是服务端发送的消息，表示
// WWW-Authenticate: Digest realm="x",nonce="x",qop="x"
type RegisterHeaderWWWAuthenticateDigest struct {
	Realm string
	Nonce string
	QOP   string
}

// Parse 解析 realm="x",nonce="x",qop="x"
func (m *RegisterHeaderWWWAuthenticateDigest) Parse(line string) bool {
	prefix, suffix := "", line
	for {
		// 去空白
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		// kv,kv,
		prefix, suffix = gs.Split(suffix, gs.CharComma)
		//
		k, v := kvQuotationMark(prefix)
		switch k {
		case StrRealm:
			m.Realm = v
		case StrNonce:
			m.Nonce = v
		case StrQOP:
			m.QOP = v
		}
	}
	return true
}

// String 返回 Digest realm="x",nonce="x",qop="x"
func (m *RegisterHeaderWWWAuthenticateDigest) String() string {
	return fmt.Sprintf(`%s %s="%s",%s="%s",%s="%s"`, StrDigest, StrRealm, m.Realm, StrQOP, m.QOP, StrNonce, m.Nonce)
}

// RegisterHeaderAuthorizationDigest 是客户端发送的消息，表示
// Authorization: Digest user="x",realm="x",nonce="x",uri="x",response="x",algorithm=x,cnonce="x",nc="x",qop="x"
type RegisterHeaderAuthorizationDigest struct {
	Username  string
	Realm     string
	Nonce     string
	URI       string
	Response  string
	Algorithm string
	CNonce    string
	NC        string
	QOP       string
}

// Parse 解析 user="x",realm="x",nonce="x",uri="x",response="x",algorithm=x,cnonce="x",nc="x",qop="x"
func (m *RegisterHeaderAuthorizationDigest) Parse(line string) bool {
	prefix, suffix := "", line
	for {
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		// x , x
		prefix, suffix = gs.Split(suffix, gs.CharComma)
		// k="v"
		k, v := kvQuotationMark(prefix)
		switch k {
		case StrUsername:
			m.Username = v
		case StrRealm:
			m.Realm = v
		case StrNonce:
			m.Nonce = v
		case StrURI:
			m.URI = v
		case StrResponse:
			m.Response = v
		case StrAlgorithm:
			m.Algorithm = v
		case StrCNonce:
			m.CNonce = v
		case StrNC:
			m.NC = v
		case StrQOP:
			m.QOP = v
		}
	}
	return true
}

// sign 返回签名
// 算法 hex(hash(hash(username:realm:password):nonce:hash(method:uri)))
func (m *RegisterHeaderAuthorizationDigest) sign(password string) string {
	h := gbu.NewHash(m.Algorithm)
	w := bufio.NewWriter(h)
	// hash(username:realm:password)
	fmt.Fprintf(w, "%s:%s:%s", m.Username, m.Realm, password)
	w.Flush()
	h1 := hex.EncodeToString(h.Sum(nil))
	// hash(method:uri)
	h.Reset()
	fmt.Fprintf(w, "%s:%s", sip.MethodRegister, m.URI)
	w.Flush()
	h2 := hex.EncodeToString(h.Sum(nil))
	// hash(h1:nonce[:nc:cnonce:qop]:h2)
	h.Reset()
	w.WriteString(h1)
	w.WriteByte(':')
	w.WriteString(m.Nonce)
	if strings.Contains(m.QOP, "auth") {
		if m.NC != "" {
			w.WriteByte(':')
			w.WriteString(m.NC)
		}
		if m.CNonce != "" {
			w.WriteByte(':')
			w.WriteString(m.CNonce)
		}
		w.WriteByte(':')
		w.WriteString(m.QOP)
	}
	w.WriteByte(':')
	w.WriteString(h2)
	w.Flush()
	//
	return hex.EncodeToString(h.Sum(nil))
}

// GenResponse 生成 response 字段，客户端调用
// 算法 hex(hash(hash(username:realm:password):nonce:hash(method:uri)))
func (m *RegisterHeaderAuthorizationDigest) GenResponse(password string) {
	m.Response = m.sign(password)
}

// VerifyResponse 验证 response 字段，服务端调用
// 算法 hex(hash(hash(username:realm:password):nonce:hash(method:uri)))
func (m *RegisterHeaderAuthorizationDigest) VerifyResponse(username, password string) bool {
	return m.Username == username && m.Response == m.sign(password)
}

// String 返回 Digest username="x",realm="x",nonce="x",uri="x",response="x",algorithm=x,cnonce="x",nc="x",qop="x"
func (m *RegisterHeaderAuthorizationDigest) String() string {
	var str strings.Builder
	fmt.Fprintf(&str, `%s %s="%s",%s="%s",%s="%s",%s="%s",%s=%s`, StrDigest,
		StrUsername, m.Username,
		StrRealm, m.Realm,
		StrNonce, m.Nonce,
		StrURI, m.URI,
		StrAlgorithm, m.Algorithm)
	if m.Response != "" {
		fmt.Fprintf(&str, `,%s="%s"`, StrResponse, m.Response)
	}
	if m.CNonce != "" {
		fmt.Fprintf(&str, `,%s="%s"`, StrCNonce, m.CNonce)
	}
	if m.QOP != "" {
		fmt.Fprintf(&str, `,%s="%s"`, StrQOP, m.QOP)
	}
	if m.NC != "" {
		fmt.Fprintf(&str, `,%s="%s"`, StrNC, m.NC)
	}
	return str.String()
}

// RegisterHeaderWWWAuthenticateAsymmetric 是服务端发送的消息，表示
// WWW-Authenticate: Asymmetric nonce="x&x" algorithm="A:x;H:x;S:x;"
type RegisterHeaderWWWAuthenticateAsymmetric struct {
	Nonce   string
	A, H, S string
}

// Parse 解析 nonce="x&x" algorithm="A:x;H:x;S:x;"
func (m *RegisterHeaderWWWAuthenticateAsymmetric) Parse(line string) bool {
	prefix, suffix := "", line
	for {
		// 去空白
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		// x空白x
		prefix, suffix = gs.Split(suffix, gs.CharSpace)
		//
		k, v := kvQuotationMark(prefix)
		switch k {
		case StrNonce:
			m.Nonce = v
		case StrAlgorithm:
			return m.parseAlgorithm(v)
		}
	}
	return true
}

// parseAlgorithm 解析 A:x;H:x;S:x
func (m *RegisterHeaderWWWAuthenticateAsymmetric) parseAlgorithm(line string) bool {
	prefix, suffix := "", line
	for {
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		prefix, suffix = gs.Split(suffix, gs.CharSemicolon)
		prefix = gs.TrimByte(prefix, gs.CharSpace, gs.CharSpace)
		str := strings.TrimPrefix(prefix, "A:")
		if str != prefix {
			m.A = str
			continue
		}
		str = strings.TrimPrefix(prefix, "H:")
		if str != prefix {
			m.H = str
			continue
		}
		str = strings.TrimPrefix(prefix, "S:")
		if str != prefix {
			m.S = str
			continue
		}
	}
	return true
}

// String 返回 Asymmetric nonce="x&x" algorithm="A:x;H:x;S:x;"
func (m *RegisterHeaderWWWAuthenticateAsymmetric) String() string {
	var str strings.Builder
	if m.A != "" {
		fmt.Fprintf(&str, "A:%s;", m.A)
	}
	if m.H != "" {
		fmt.Fprintf(&str, "H:%s;", m.H)
	}
	if m.S != "" {
		fmt.Fprintf(&str, "S:%s;", m.S)
	}
	return fmt.Sprintf(`%s %s="%s" %s="%s"`, StrAsymmetric, StrNonce, m.Nonce, StrAlgorithm, str.String())
}

// GenNonce 返回签名，服务端调用
// 算法 base64(ser.sign(hash(rand)))&base64(cli.enc(rand))
func (m *RegisterHeaderWWWAuthenticateAsymmetric) GenNonce(ser *rsa.PrivateKey, cli *rsa.PublicKey) error {
	a, b, _, err := rsaGenNonce(ser, cli, gbu.CryptoHash(m.H))
	if err != nil {
		return err
	}
	//
	m.Nonce = base64.StdEncoding.EncodeToString(a) + "&" + base64.StdEncoding.EncodeToString(b)
	//
	return nil
}

// Verify 验证并返回随机数，客户端使用
// 算法 ser.verify(hash(cli.dec(base64.dec(nonce2))), base64.dec(nonce1))
func (m *RegisterHeaderWWWAuthenticateAsymmetric) VerifyNonce(ser *rsa.PublicKey, cli *rsa.PrivateKey) ([]byte, error) {
	return rsaVerifyNonce(cli, ser, gbu.CryptoHash(m.H), m.Nonce)
}

// RegisterHeaderAuthorization 表示
// Authorization: Digest/Capability/Asymmetric
type RegisterHeaderAuthorization struct {
	Digest     *RegisterHeaderAuthorizationDigest
	Capability *RegisterHeaderAuthorizationCapability
	Asymmetric *RegisterHeaderAuthorizationAsymmetric
}

// Parse 解析
func (m *RegisterHeaderAuthorization) Parse(line string) bool {
	prefix, suffix := gs.Split(line, gs.CharSpace)
	// 类型，Digest/Capability/Asymmetric
	switch prefix {
	case StrDigest:
		m.Digest = new(RegisterHeaderAuthorizationDigest)
		return m.Digest.Parse(suffix)
	case StrCapability:
		m.Capability = new(RegisterHeaderAuthorizationCapability)
		return m.Capability.Parse(suffix)
	case StrAsymmetric:
		m.Asymmetric = new(RegisterHeaderAuthorizationAsymmetric)
		return m.Asymmetric.Parse(suffix)
	default:
		return false
	}
}

// RegisterHeaderAuthorizationCapability 是服务端发送的消息，表示
// Authorization: Capability algorithm = "A:x;H:x;S:x"
type RegisterHeaderAuthorizationCapability struct {
	A, H, S string
}

// Parse 解析
func (m *RegisterHeaderAuthorizationCapability) Parse(line string) bool {
	prefix, suffix := "", line
	for {
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		// x , x
		prefix, suffix = gs.Split(suffix, gs.CharComma)
		// k="v"
		k, v := kvQuotationMark(prefix)
		switch k {
		case StrAlgorithm:
			m.parseAlgorithm(v)
		}
	}
	return true
}

// parseAlgorithm 解析 A:x;H:x;S:x
func (m *RegisterHeaderAuthorizationCapability) parseAlgorithm(line string) {
	prefix, suffix := "", line
	for {
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		prefix, suffix = gs.Split(suffix, gs.CharSemicolon)
		prefix = gs.TrimByte(prefix, gs.CharSpace, gs.CharSpace)
		str := strings.TrimPrefix(prefix, "A:")
		if str != prefix {
			m.A = str
			continue
		}
		str = strings.TrimPrefix(prefix, "H:")
		if str != prefix {
			m.H = str
			continue
		}
		str = strings.TrimPrefix(prefix, "S:")
		if str != prefix {
			m.S = str
			continue
		}
	}
}

// String 返回 Capability algorithm="A:x;H:x;S:x;"
func (m *RegisterHeaderAuthorizationCapability) String() string {
	var str strings.Builder
	if m.A != "" {
		fmt.Fprintf(&str, "A:%s;", m.A)
	}
	if m.H != "" {
		fmt.Fprintf(&str, "H:%s;", m.H)
	}
	if m.S != "" {
		fmt.Fprintf(&str, "S:%s;", m.S)
	}
	return fmt.Sprintf(`%s %s="%s"`, StrCapability, StrAlgorithm, str.String())
}

// RegisterHeaderAuthorizationAsymmetric 是客户端发送的消息，表示
// Authorization: Asymmetric nonce="x&x" algorithm=SHA1
type RegisterHeaderAuthorizationAsymmetric struct {
	Nonce     string
	Response  string
	Algorithm string
	// 服务端的随机数，由 RegisterHeaderWWWAuthenticateAsymmetric.Verify 返回
	C []byte
}

// Parse 解析 nonce="x&x" algorithm=SHA1
func (m *RegisterHeaderAuthorizationAsymmetric) Parse(line string) bool {
	prefix, suffix := "", line
	for {
		suffix = gs.TrimByte(suffix, gs.CharSpace, gs.CharSpace)
		if suffix == "" {
			break
		}
		prefix, suffix = gs.Split(suffix, gs.CharComma)
		// k="v"
		k, v := kvQuotationMark(prefix)
		switch k {
		case StrNonce:
			m.Nonce = v
		case StrResponse:
			m.Response = v
		case StrAlgorithm:
			m.Algorithm = v
		}
	}
	return true
}

// GenResponse 生成 response ，客户端调用
// 算法 hash(c+nonce)
func (m *RegisterHeaderAuthorizationAsymmetric) GenResponse() (string, error) {
	hash := gbu.CryptoHash(m.Algorithm)
	h := hash.New()
	h.Reset()
	h.Write(m.C)
	h.Write([]byte(m.Nonce))
	d := h.Sum(nil)
	return string(d), nil
}

// Verify 验证 response ，服务端调用
// 算法 cli.verify(hash(ser.dec(base64.dec(nonce2))), base64.dec(nonce1))
func (m *RegisterHeaderAuthorizationAsymmetric) VerifyResponse(ser *rsa.PrivateKey, cli *rsa.PublicKey) (bool, error) {
	hash := gbu.CryptoHash(m.Algorithm)
	// 验证 nonce
	c, err := rsaVerifyNonce(ser, cli, hash, m.Nonce)
	if err != nil {
		return false, err
	}
	// 验证 hash(c+nonce)==response
	h := hash.New()
	h.Reset()
	h.Write(c)
	h.Write([]byte(m.Nonce))
	d := h.Sum(nil)
	//
	return string(d) == m.Response, nil
}
