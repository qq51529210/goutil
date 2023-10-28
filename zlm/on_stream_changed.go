package zlm

// OnStreamChangedReq 表示 on_stream_changed 提交的数据，保存注册和注销的所有字段
type OnStreamChangedReq struct {
	// 服务器id，通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流注册或注销
	Regist bool `json:"regist"`
	// 流的媒体信息
	MediaListData
	// 日志追踪
	TraceID string
}

// // OnStreamChanged 处理 zlm 的 on_stream_changed 回调
// func OnStreamChanged(ctx *gin.Context, req *OnStreamChangedReq) {
// 	// 获取实例
// 	ser := GetServer(req.MediaServerID)
// 	if !ser.IsOK() {
// 		return
// 	}
// 	// 第一次才回调
// 	info := ser.onStreamChanged(req)
// 	if info != nil {
// 		req.TraceID, _ = ctx.Value(CtxKeyTraceID).(string)
// 		HandleStreamChanged(ctx, ser, req, info)
// 	}
// }

// // onStreamChanged 处理 zlm 的 on_stream_changed 回调，返回是否第一次
// func (s *Server) onStreamChanged(req *OnStreamChangedReq) *MediaInfo {
// 	// 处理 rtmp 即可，因为 gb / 推流 / 拉流 / onvif 都会启用它
// 	if req.Schema != RTMP {
// 		return nil
// 	}
// 	key := mediaInfoKey{App: req.App, Stream: req.Stream}
// 	// 上锁
// 	s.lock.Lock()
// 	defer s.lock.Unlock()
// 	//
// 	info := s.mediaInfos[key]
// 	// 流注册
// 	if req.Regist {
// 		// 表里没有
// 		if info == nil {
// 			info = new(MediaInfo)
// 			info.init(s, &req.MediaListData)
// 			// 添加
// 			s.mediaInfos[key] = info
// 			//
// 			return info
// 		}
// 		// 表里有
// 		return nil
// 	}
// 	// 流注销，移除
// 	if info != nil {
// 		delete(s.mediaInfos, key)
// 	}
// 	return info
// }
