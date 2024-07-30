package zlm

import "context"

// SeekRecordStampReq 是 SeekRecordStamp 的参数
type SeekRecordStampReq struct {
	// 协议，这个 postman 上没有，但是不添不行
	Schema string `query:"schema"`
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 要设置的录像播放位置，单位毫秒
	Stamp string `query:"stamp"`
}

// SeekRecordStampRes 是 SeekRecordStamp 返回值
type SeekRecordStampRes struct {
	CodeMsg
}

const (
	SeekRecordStampPath = apiPathPrefix + "/seekRecordStamp"
)

// SeekRecordStamp 调用 /index/api/seekRecordStamp ，设置录像的播放位置
func SeekRecordStamp(ctx context.Context, ser Server, req *SeekRecordStampReq, res *SeekRecordStampRes) error {
	return Request(ctx, ser, SeekRecordStampPath, req, res)
}
