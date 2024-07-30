package zlm

import "context"

// DeleteRecordDirectoryReq 是 DeleteRecordDirectory 的参数
type DeleteRecordDirectoryReq struct {
	// 流应用
	App string `query:"app"`
	// 流标识
	Stream string `query:"stream"`
	// 流的录像日期，格式为2020-01-01
	// 如果不是完整的日期，那么会删除失败
	Period string `query:"period"`
}

// DeleteRecordDirectoryRes 是 DeleteRecordDirectory 返回值
type DeleteRecordDirectoryRes struct {
	CodeMsg
}

const (
	DeleteRecordDirectoryPath = apiPathPrefix + "/deleteRecordDirectory"
)

// DeleteRecordDirectory 调用 /index/api/deleteRecordDirectory ，删除录像文件目录
// 经过测试，删除不存在的会返回 code=-1
func DeleteRecordDirectory(ctx context.Context, ser Server, req *DeleteRecordDirectoryReq, res *DeleteRecordDirectoryRes) error {
	return Request(ctx, ser, DeleteRecordDirectoryPath, req, res)
}
