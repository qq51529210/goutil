package zlm

// OnStreamNoneReaderReq 表示 on_stream_none_reader 提交的数据
type OnStreamNoneReaderReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 播放的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 日志追踪
	TraceID string `json:"-"`
}

// OnStreamNoneReaderRes 表示 on_stream_none_reader 返回值
type OnStreamNoneReaderRes struct {
	Close *bool `json:"close"`
	Code  int   `json:"code"`
}

// Keep 不关闭流
func (r *OnStreamNoneReaderRes) Keep() {
	ok := false
	r.Close = &ok
}

// // OnStreamNoneReader 处理 zlm 的 on_stream_none_reader 回调
// func OnStreamNoneReader(ctx *gin.Context, req *OnStreamNoneReaderReq, res *OnStreamNoneReaderRes) {
// 	// 获取实例
// 	ser := GetServer(req.MediaServerID)
// 	if !ser.IsOK() {
// 		return
// 	}
// 	req.TraceID, _ = ctx.Value(CtxKeyTraceID).(string)
// 	// 回调
// 	close := true
// 	res.Close = &close
// 	HandleStreamNoneReader(ctx, ser, req, res)
// 	// 日志
// 	log.Debugf("%s res close %t", req.TraceID, *res.Close)
// }
