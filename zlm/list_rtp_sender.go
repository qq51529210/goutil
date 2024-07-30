package zlm

import "context"

// ListRTPSenderReq 是 ListRTPSender 的参数
type ListRTPSenderReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
}

// listRTPSenderRes 是 ListRTPSender 返回值
type listRTPSenderRes struct {
	CodeMsg
	// 文档没有说明，我看代码的
	// 应该是返回了 ssrc 数组
	// val["data"].append(ssrc);
	Data []string
}

const (
	ListRTPSenderPath = apiPathPrefix + "/listRtpSender"
)

// ListRTPSender 调用 /index/api/listRtpSender ，返回所有 rtp server 的 ssrc
func ListRTPSender(ctx context.Context, ser Server, req *ListRTPSenderReq, res *ListRTPSenderReq) error {
	return Request(ctx, ser, ListRTPSenderPath, req, res)
}
