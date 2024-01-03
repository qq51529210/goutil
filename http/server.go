package http

import "net/http"

// Server 封装代码
type Server struct {
	S http.Server
	// 证书路径
	CertFile string
	KeyFile  string
}

// Serve 如果证书路径不为空，监听 tls
func (s *Server) Serve() error {
	if s.CertFile != "" && s.KeyFile != "" {
		return s.S.ListenAndServeTLS(s.CertFile, s.KeyFile)
	}
	return s.S.ListenAndServe()
}
