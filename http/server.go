package http

import (
	"crypto/tls"
	"net"
	"net/http"
)

// Server 封装代码
type Server struct {
	http.Server
	// 证书路径
	CertFile string
	KeyFile  string
	// 证书数据，优先比路径高
	CertPem string
	KeyPem  string
}

// Serve 如果证书路径不为空，监听 tls
func (s *Server) Serve() error {
	// 证书
	if s.CertPem != "" && s.KeyPem != "" {
		cert, err := tls.X509KeyPair([]byte(s.CertPem), []byte(s.KeyPem))
		if err != nil {
			return err
		}
		var cfg tls.Config
		cfg.Certificates = append(cfg.Certificates, cert)
		cfg.NextProtos = append(cfg.NextProtos, "http/1.1", "h2")
		// 监听
		l, err := net.Listen("tcp", s.Addr)
		if err != nil {
			return err
		}
		l = tls.NewListener(l, &cfg)
		//
		return s.Server.Serve(l)
	}
	// 证书路径
	if s.CertFile != "" && s.KeyFile != "" {
		return s.ListenAndServeTLS(s.CertFile, s.KeyFile)
	}
	// 普通
	return s.ListenAndServe()
}
