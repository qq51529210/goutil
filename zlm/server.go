package zlm

// var (
// 	_servers gs.MapSlice[string, *Server]
// )

// func init() {
// 	_servers.Init()
// }

// Server 表示服务
type Server struct {
	ID         string
	Secret     string
	APIBaseURL string
}

// // RTMPURL 返回 rtmp 地址列表
// func (s *Server) RTMPURL(app, stream, token string) []string {
// 	var urls []string
// 	var str strings.Builder
// 	// rtmp
// 	if s.RTMPPort != "" {
// 		str.Reset()
// 		fmt.Fprintf(&str, "rtmp://%s:%s/%s/%s", s.PublicIP, s.RTMPPort, app, stream)
// 		if token != "" {
// 			fmt.Fprintf(&str, "?token=%s", token)
// 		}
// 		urls = append(urls, str.String())
// 	}
// 	if s.RTMPSSLPort != "" {
// 		str.Reset()
// 		fmt.Fprintf(&str, "rtmps://%s:%s/%s/%s", s.PublicIP, s.RTMPSSLPort, app, stream)
// 		if token != "" {
// 			fmt.Fprintf(&str, "?token=%s", token)
// 		}
// 		urls = append(urls, str.String())
// 	}
// 	//
// 	return urls
// }

// // RTSPURL 返回 rtsp 地址列表
// func (s *Server) RTSPURL(app, stream, token string) []string {
// 	var urls []string
// 	var str strings.Builder
// 	// rtsp
// 	if s.RTSPPort != "" {
// 		str.Reset()
// 		fmt.Fprintf(&str, "rtsp://%s:%s/%s/%s", s.PublicIP, s.RTSPPort, app, stream)
// 		if token != "" {
// 			fmt.Fprintf(&str, "?token=%s", token)
// 		}
// 		urls = append(urls, str.String())
// 	}
// 	if s.RTSPSSLPort != "" {
// 		str.Reset()
// 		fmt.Fprintf(&str, "rtsps://%s:%s/%s/%s", s.PublicIP, s.RTSPSSLPort, app, stream)
// 		if token != "" {
// 			fmt.Fprintf(&str, "?token=%s", token)
// 		}
// 		urls = append(urls, str.String())
// 	}
// 	//
// 	return urls
// }

// // FFMPEGPull 封装 AddFFMPEGSource
// func (s *Server) FFMPEGPull(ctx context.Context, url, app, stream, timeoutMS string) error {
// 	_, err := s.AddFFMPEGSource(ctx, &AddFFMPEGSourceReq{
// 		SrcURL:    url,
// 		DstURL:    fmt.Sprintf("rtmp://%s:%s/%s/%s", s.PrivateIP, s.RTMPPort, app, stream),
// 		TimeoutMS: timeoutMS,
// 		EnableMP4: False, // 这个在参数 dst_url=localhost 不起作用，还是会录制
// 	})
// 	return err
// }

// // GetAllServer 获取所有
// func GetAllServer() []*Server {
// 	return _servers.All()
// }

// // GetServer 获取
// func GetServer(id string) *Server {
// 	return _servers.Get(id)
// }

// // RestartServer 重启服务
// func RestartServer(ctx context.Context, id string) error {
// 	s := GetServer(id)
// 	if s != nil {
// 		return s.RestartServer(ctx)
// 	}
// 	return ErrServerNotAvailable
// }

// // BatchSetServer 更新全部
// func BatchSetServer(ms []*Server) {
// 	// 组装
// 	s := make([]*Server, 0, len(ms))
// 	d := make(map[string]*Server)
// 	for _, m := range ms {
// 		p := (*Server)(m)
// 		d[m.ID] = p
// 		s = append(s, p)
// 	}
// 	// 替换
// 	_servers.Lock()
// 	_servers.D = d
// 	_servers.S = s
// 	_servers.Unlock()
// }
