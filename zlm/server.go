package zlm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"util"

	"util/log"
)

var (
	_servers servers
)

// Server 表示一个 zlm 服务
type Server struct {
	lock sync.RWMutex
	// 是否有效
	ok int32
	// 是否在线
	online bool
	// 数据版本
	Version int64
	// 标识
	ID string
	// 密钥
	Secret string
	// api 基本路径，http(https)://hostname:port/
	APIBaseURL string
	// 外部 ip
	PublicIP string
	// 内部 ip
	PrivateIP string
	// 是否ip6 ，sdp 会用到
	IsIPV6 bool
	// 接口调用超时，后台使用
	APICallTimeout time.Duration
	// 媒体流同步/自动截图间隔
	SyncMediaDur time.Duration
	// 截图超时，单位秒
	SnapshotTimeout int64
	// 截图存放目录，空则不截图
	SnapshotDir string
	// 配置
	*Config
	// 心跳时间，on_server_keepalive 回调
	keepaliveTime *time.Time
	// 心跳超时，从 zlm 的配置读取的，加了 3 秒作为缓冲
	keepaliveTimeout time.Duration
	// 最近一次 api 调用时间
	// 用于判断是否在线（主动调用成功，也算它在线）
	apiCallTime *time.Time
	// 流列表
	mediaInfos util.MapSlice[mediaInfoKey, *MediaInfo]
	// 是否正在请求配置
	syncConfig int32
	// 是否正在请求媒体流
	syncMedia int32
	// 同步媒体流时间
	syncMediaTime *time.Time
	// 正在请求拉流的表，防止重复调用
	pullStreams util.Set[mediaInfoKey]
	// 正在请求截图的表，防止重复调用
	snapshots util.Set[mediaInfoKey]
	// 国标实时流 ssrc 表，接收 rtp 推流使用
	realTimeSSRC ssrcPool
	// 国标历史流 ssrc 表，接收 rtp 推流使用
	historySSRC ssrcPool
	// 是否有更新
	isUpdated bool
	// 总播放人数
	player int
}

// IsOK 返回 s 是否可用
func (s *Server) IsOK() bool {
	return s != nil && s.online
}

// checkAndSync 检查，同步数据
func (s *Server) checkAndSync(now *time.Time) {
	// 未准备好
	if atomic.LoadInt32(&s.ok) == 0 {
		// 获取配置
		if atomic.CompareAndSwapInt32(&s.syncConfig, 0, 1) {
			go s.getConfigRoutine()
		}
		return
	}
	// 同步媒体流
	t := s.syncMediaTime
	if now.Sub(*t) >= s.SyncMediaDur && atomic.CompareAndSwapInt32(&s.syncMedia, 0, 1) {
		go s.getMediaRoutine()
	}
	// 是否离线
	online := s.online
	keealiveTime := s.keepaliveTime
	apiCallTime := s.apiCallTime
	if now.Sub(*keealiveTime) >= s.keepaliveTimeout && now.Sub(*apiCallTime) >= s.keepaliveTimeout {
		online = false
	} else {
		online = true
	}
	s.isUpdated = online != s.online
}

// getConfigRoutine 在协程中请求服务配置
func (s *Server) getConfigRoutine() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.APICallTimeout))
	defer func() {
		cancel()
		// 异常
		log.Recover(recover())
		// 开关
		atomic.StoreInt32(&s.syncConfig, 0)
	}()
	// 请求
	err := s.GetServerConfig(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	// 成功
	atomic.StoreInt32(&s.ok, 1)
	//
	s.online = true
	s.isUpdated = true
}

// getMediaRoutine 在协程中请求媒体流
func (s *Server) getMediaRoutine() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(s.APICallTimeout))
	defer func() {
		cancel()
		// 异常
		log.Recover(recover())
		// 开关
		atomic.StoreInt32(&s.syncMedia, 0)
	}()
	// 请求
	_, err := s.GetMediaList(ctx, &GetMediaListReq{
		Schema: RTMP,
	})
	if err != nil {
		log.Error(err)
		return
	}
	// 时间
	now := time.Now()
	s.syncMediaTime = &now
	// 请求截图
	if s.SnapshotDir == "" {
		return
	}
	timeout := fmt.Sprintf("%d", s.SnapshotTimeout)
	for _, v := range s.mediaInfos.All() {
		// 没有视频轨道，不截图
		if len(v.Video) < 1 {
			continue
		}
		// 	// 回放，下载，不截图
		// 	if v.App == RTP || v.App == ONVIF {
		// 		if strings.HasPrefix(v.Stream, PLAYBACK) ||
		// 			strings.HasPrefix(v.Stream, DOWNLOAD) {
		// 			continue
		// 		}
		// 	}
		key := mediaInfoKey{App: v.App, Stream: v.Stream}
		if s.snapshots.TrySet(key) {
			go s.getSnapshotRoutine(&key, timeout)
		}
	}
}

// getSnapshotRoutine 在协程中请求截图
func (s *Server) getSnapshotRoutine(key *mediaInfoKey, timeout string) {
	defer func() {
		// 异常
		log.Recover(recover())
		// 移除列表
		s.snapshots.Del(*key)
	}()
	// 截图
	err := s.SaveSnap(key.App, key.Stream, timeout)
	if err != nil {
		log.Error(err)
		return
	}
}

// servers 用于管理所有 Server
type servers struct {
	util.MapSlice[string, *Server]
	// 最小心跳
	d time.Duration
	// 计时器
	t *time.Timer
}

// // init 初始化
// func (ss *servers) init() error {
// 	// 数据库
// 	var q db.ServerQuery
// 	q.Enable = &db.True
// 	ms, err := db.AllServer(context.Background(), &q)
// 	if err != nil {
// 		return err
// 	}
// 	// 内存
// 	ss.D = make(map[string]*Server)
// 	for _, m := range ms {
// 		ss.D[m.ID] = ss.new(m)
// 	}
// 	ss.resetSlice()
// 	// 计时器
// 	ss.t = time.NewTimer(ss.d)
// 	// 检查协程
// 	go ss.checkServerRoutine()
// 	// 更新数据库协程
// 	go ss.updateDBRoutine()
// 	//
// 	return nil
// }

// // add 添加
// func (ss *servers) add(m *db.Server) {
// 	// 上锁
// 	ss.Lock()
// 	defer ss.Unlock()
// 	// 检查
// 	s := ss.D[m.ID]
// 	if s != nil {
// 		// 检查数据版本，防止并发修改操作出错
// 		if s.modelVersion > m.UpdatedAt {
// 			// 内存是新数据
// 			return
// 		}
// 	}
// 	// 添加
// 	ss.D[m.ID] = ss.new(m)
// 	ss.resetSlice()
// }

// // new 创建
// func (ss *servers) new(m *db.Server) *Server {
// 	s := new(Server)
// 	s.Server = m
// 	s.modelVersion = m.UpdatedAt
// 	s.mediaInfos = make(map[mediaInfoKey]*MediaInfo)
// 	s.pullStreams = make(map[mediaInfoKey]struct{})
// 	s.snapshots = make(map[mediaInfoKey]struct{})
// 	keepaliveTime := time.Unix(m.KeepaliveTime, 0)
// 	s.keepaliveTime = &keepaliveTime
// 	s.syncMediaTime = &time.Time{}
// 	s.syncMediaDur = time.Duration(m.SyncMediaInterval) * time.Second
// 	s.apiCallTime = &time.Time{}
// 	// 不判断了，api 那里判断过了
// 	ip, _ := netip.ParseAddr(m.PublicIP)
// 	s.IsIPV6 = ip.Is6()
// 	return s
// }

// // resetSlice 重新组装一下切片，顺便找出最小心跳
// func (ss *servers) resetSlice() {
// 	// 心跳比这个还大就没意义了
// 	dur := time.Hour
// 	// 重置数据
// 	ss.S = make([]*Server, 0, len(ss.D))
// 	for _, v := range ss.D {
// 		if v.keepaliveTimeout < dur {
// 			dur = v.keepaliveTimeout
// 		}
// 		ss.S = append(ss.S, v)
// 	}
// 	ss.d = dur
// 	// 重置计时器
// 	ss.t = time.NewTimer(ss.d)
// }

// // initSSRC 查询 invite 信息，初始化所有服务的 ssrc 池
// func (ss *servers) initSSRC() error {
// 	// // 数据库
// 	// ms, err := db.DeviceInviteDA.All()
// 	// if err != nil {
// 	// 	return fmt.Errorf("zlm init ssrc Err:%s", err.Error())
// 	// }
// 	// // 更新
// 	// for _, m := range ms {
// 	// 	s := ss.get(m.ID)
// 	// 	if s != nil {
// 	// 		s.PutSSRC(m.SSRC)
// 	// 	}
// 	// }
// 	//
// 	return nil
// }

// // checkServerRoutine 在协程中定时检查每一个 Server
// func (ss *servers) checkServerRoutine() {
// 	log.InfoTrace(logTraceID, "check server routine start")
// 	defer func() {
// 		log.InfoTrace(logTraceID, "check server routine stop")
// 		// 异常
// 		log.Recover(recover())
// 	}()
// 	// 循环检查
// 	for {
// 		now := <-ss.t.C
// 		// 检查
// 		ms := ss.All()
// 		for _, m := range ms {
// 			m.checkAndSync(&now)
// 		}
// 		// 重置计时器
// 		ss.t.Reset(ss.d)
// 	}
// }

// // updateDBRoutine 在协程中定时更新需要更新的数据
// // 集中更新是为了避免并发更新出现的 too many connections
// func (ss *servers) updateDBRoutine() {
// 	log.InfoTrace(logTraceID, "udpate db routine start")
// 	defer func() {
// 		log.InfoTrace(logTraceID, "udpate db routine stop")
// 		// 异常
// 		log.Recover(recover())
// 	}()
// 	// 循环检查
// 	for {
// 		<-ss.t.C
// 		ctx := context.Background()
// 		ms := ss.All()
// 		for _, m := range ms {
// 			// 如果更新失败，等下一次
// 			if atomic.CompareAndSwapInt32(&m.updateDB, 1, 0) {
// 				db.UpdateServerOnline(ctx, m.Server)
// 			}
// 		}
// 		// 重置计时器
// 		ss.t.Reset(ss.d)
// 	}
// }

// // loadServerRoutine 在协程中加载指定 id 数据
// func (ss *servers) loadServerRoutine(id string) {
// 	// 计时器
// 	timer := time.NewTimer(0)
// 	defer func() {
// 		// 异常
// 		log.Recover(recover())
// 		// 计时器
// 		timer.Stop()
// 	}()
// 	// 数据库
// 	m := ss.mustLoadServer(timer, id)
// 	// 不存在/禁用
// 	if m == nil || *m.Enable == db.False {
// 		// 清理内存
// 		RemoveServer(id)
// 		return
// 	}
// 	// 添加
// 	ss.add(m)
// }

// // mustLoadServer 确保成功加载指定 id 数据
// func (ss *servers) mustLoadServer(timer *time.Timer, id string) *db.Server {
// 	m := new(db.Server)
// 	m.ID = id
// 	// 加载
// 	for {
// 		<-timer.C
// 		// 查询
// 		err := db.GetServer(context.Background(), m)
// 		if err == nil {
// 			return m
// 		}
// 		// 没有数据
// 		if db.IsDataNotFound(err) {
// 			return nil
// 		}
// 		log.ErrorfTrace(logTraceID, "load db %s error %s", id, err.Error())
// 		// 失败，重试
// 		timer.Reset(time.Second)
// 	}
// }

// LoadServer 加载
func LoadServer(id string) {
	// go _servers.loadServerRoutine(id)
}

// GetServer 获取
func GetServer(id string) *Server {
	return _servers.Get(id)
}

// RemoveServer 删除
func RemoveServer(id string) {
	_servers.Del(id)
}

// BatchRemoveServer 批量删除
func BatchRemoveServer(ids []string) {
	_servers.BatchDel(ids)
}

// GetMinLoadServer 返回所有服务中，最小收流负载的服务
func GetMinLoadServer() *Server {
	// 上锁
	_servers.RLock()
	defer _servers.RUnlock()
	// 负载
	var ser *Server
	load := -1
	for _, s := range _servers.D {
		// 在线
		if !s.online {
			continue
		}
		// 负载
		n := s.mediaInfos.Len()
		if load < 1 || load > n {
			load = n
			ser = s
		}
	}
	// 返回
	return ser
}

// GetMinPlayerServer 返回所有服务中，最小播放负载的服务
func GetMinPlayerServer() *Server {
	// 上锁
	_servers.RLock()
	defer _servers.RUnlock()
	// 负载
	var ser *Server
	load := -1
	for _, s := range _servers.D {
		// 在线
		if !s.online {
			continue
		}
		// 负载
		s.mediaInfos.RLock()
		n := s.player
		s.mediaInfos.RUnlock()
		if load < 1 || load > n {
			load = n
			ser = s
		}
	}
	// 返回
	return ser
}

// // PullStream 尝试拉流，并发控制的，多次调用也只会有一次调用
// func (s *Server) PullStream(m *db.PullStream) {
// 	// 禁用
// 	if *m.Enable != db.True {
// 		return
// 	}
// 	// 出错后台继续
// 	key := mediaInfoKey{App: PULL, Stream: m.ID}
// 	// 上锁
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	// 看看流是否存在
// 	if _, ok := s.pullStreams[key]; !ok {
// 		s.pullStreams[key] = struct{}{}
// 		go s.pullStreamRoutine(m, key)
// 	}
// }

// // pullStream 拉流
// func (s *Server) pullStream(m *db.PullStream) error {
// 	// 请求
// 	var req AddFFMPEGSourceReq
// 	req.SrcURL = m.SrcURL
// 	// req.DstURL = fmt.Sprintf("rtmp://localhost:%s/%s/%s", s.RTMPPort, PULL, m.ID)
// 	// 如果用localhost那么它就不会触发on_publish
// 	req.DstURL = fmt.Sprintf("rtmp://%s:%s/%s/%s", s.PrivateIP, s.RTMPPort, PULL, m.ID)
// 	req.TimeoutMS = fmt.Sprintf("%d", m.Timeout)
// 	if m.FFMPEGCmd != "" {
// 		req.CmdKey = m.FFMPEGCmd
// 	}
// 	// 然后这个参数又不起作用
// 	req.EnableMP4 = False
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.Timeout)*time.Millisecond)
// 	_, err := s.AddFFMPEGSource(ctx, &req)
// 	cancel()
// 	return err
// }

// // pullStreamRoutine 在协程中尝试拉流
// func (s *Server) pullStreamRoutine(m *db.PullStream, k mediaInfoKey) {
// 	defer func() {
// 		// 异常
// 		log.Recover(recover())
// 		// 上锁
// 		s.lock.Lock()
// 		delete(s.pullStreams, k)
// 		s.lock.Unlock()
// 	}()
// 	// 请求
// 	err := s.pullStream(m)
// 	if err != nil {
// 		log.Error(err)
// 		return
// 	}
// }
