package sip

import (
	"goutil/log"
	gosync "goutil/sync"
	"net"
)

// tcpServer 包装 tcp 相关的数据
type tcpServer struct {
	l  *net.TCPListener
	c  gosync.Map[connKey, *tcpConn]
	at gosync.Map[string, *activeTx]
	pt gosync.Map[string, *passiveTx]
}

// serveTCP 开始 tcp 服务
func (s *Server) serveTCP() error {
	log.InfofTrace(logTraceTCP, "listen %s", s.Addr)
	// 初始化
	a, err := net.ResolveTCPAddr("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.tcp.l, err = net.ListenTCP(a.Network(), a)
	if err != nil {
		return err
	}
	s.udp.at.Init()
	s.udp.pt.Init()
	//
	s.w.Add(3)
	// 监听
	go s.listenTCPRoutine()
	// 检查
	go s.checkActiveTxTimeoutRoutine(logTraceTCP, &s.tcp.at)
	go s.checkPassiveTxTimeoutRoutine(logTraceTCP, &s.tcp.pt)
	// 返回
	return nil
}

// listenTCPRoutine 监听 tcp 连接，然后启动协程处理
func (s *Server) listenTCPRoutine() {
	defer func() {
		// 异常
		log.Recover(recover())
		// 结束
		s.w.Done()
	}()
	for s.isOK() {
		// 监听
		conn, err := s.tcp.l.AcceptTCP()
		if err != nil {
			log.ErrorfTrace(logTraceTCP, "accept %v", err)
			continue
		}
		// 处理
		s.w.Add(1)
		go s.handleTCPConnRoutine(s.addTCPConn(conn))
	}
}

// addTCPConn 添加并返回
func (s *Server) addTCPConn(conn *net.TCPConn) *tcpConn {
	// 初始化
	c := new(tcpConn)
	c.init(conn)
	// 添加
	s.tcp.c.Set(c.key, c)
	//
	return c
}

// delTCPConn 移除并关闭
func (s *Server) delTCPConn(c *tcpConn) {
	s.tcp.c.Del(c.key)
	c.conn.Close()
}

// getTCPConn 获取
func (s *Server) getTCPConn(a *net.TCPAddr) *tcpConn {
	k := connKey{}
	k.Init(a.IP, a.Port)
	return s.tcp.c.Get(k)
}

// dialTCPConn 创建连接
func (s *Server) dialTCPConn(a *net.TCPAddr) (*tcpConn, error) {
	conn, err := net.DialTimeout(a.Network(), a.String(), s.TxTimeout)
	if err != nil {
		return nil, err
	}
	return s.addTCPConn(conn.(*net.TCPConn)), nil
}

// handleTCPConnRoutine 处理 tcp conn 消息
func (s *Server) handleTCPConnRoutine(c *tcpConn) {
	defer func() {
		// 异常
		log.Recover(recover())
		// 移除
		s.delTCPConn(c)
		// 结束
		s.w.Done()
	}()
	//
	r := newReader(c.conn, s.MaxMessageLen)
	for s.isOK() {
		// 解析，错误直接返回关闭连接
		m := new(message)
		err := m.Dec(r, s.MaxMessageLen)
		if err != nil {
			log.WarnfTrace(logTraceTCP, "read message %v", err)
			return
		}
		// 处理
		err = s.handleMsg(c, m, &s.tcp.at, &s.tcp.pt)
		if err != nil {
			log.ErrorfTrace(logTraceTCP, "handle %v", err)
			continue
		}
	}
}

// closeTCP 停止监听，关闭所有的连接
func (s *Server) closeTCP() {
	if s.tcp.l != nil {
		log.InfoTrace(logTraceTCP, "close")
		s.tcp.l.Close()
		s.tcp.l = nil
		//
		cs := s.tcp.c.TakeAll()
		for _, c := range cs {
			c.conn.Close()
		}
	}
}
