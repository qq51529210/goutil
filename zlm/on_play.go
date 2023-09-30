package zlm

import (
	"net/url"

	"github.com/gin-gonic/gin"
)

// OnPlayReq 表示 on_play 提交的数据
type OnPlayReq struct {
	// 服务器id,通过配置文件设置
	MediaServerID string `json:"mediaServerId"`
	// 流虚拟主机
	VHost string `json:"vhost"`
	// 推流的协议，可能是rtsp、rtmp
	Schema string `json:"schema"`
	// 流应用名
	App string `json:"app"`
	// 流ID
	Stream string `json:"stream"`
	// 推流url参数
	Params string `json:"params"`
	// // TCP链接唯一ID
	// ID string `json:"id"`
	// // 推流器ip
	// IP string `json:"ip"`
	// // 推流器端口号
	// Port int `json:"port"`
}

// OnPlayRes 表示 on_play 的返回值
type OnPlayRes struct {
	apiError
}

// OnPlay 处理 zlm 的 on_play 回调
func OnPlay(ctx *gin.Context, req *OnPlayReq, res *OnPlayRes) {
	res.Code = -1
	// 解析查询参数
	query, err := url.ParseQuery(req.Params)
	if err != nil {
		return
	}
	// 检查
	if _playToken.has(query.Get(QueryNameToken)) {
		res.Code = 0
	}
}
