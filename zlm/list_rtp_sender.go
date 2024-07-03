package zlm

import "context"

// ListRTPSenderReq 是 ListRTPSender 的参数
type ListRTPSenderReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 添加的流的虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 添加的应用名，例如 live
	App string `query:"app"`
	// 添加的流id，例如 test
	Stream string `query:"stream"`
}

// listRTPSenderRes 是 ListRTPSender 返回值
type listRTPSenderRes struct {
	apiError
	// 文档没有说明，我看代码的
	// 应该是返回了 ssrc 数组
	// val["data"].append(ssrc);
	Data []string
}

const (
	apiListRTPSender = "listRtpSender"
)

// ListRTPSender 调用 /index/api/listRtpSender
func ListRTPSender(ctx context.Context, req *ListRTPSenderReq) ([]string, error) {
	// 请求
	var res listRTPSenderRes
	if err := request(ctx, req.BaseURL, apiListRTPSender, req, &res); err != nil {
		return nil, err
	}
	if res.apiError.Code != codeTrue {
		res.apiError.Path = apiListRTPSender
		return nil, &res.apiError
	}
	return res.Data, nil
}
